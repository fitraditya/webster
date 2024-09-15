package cmd

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/fitraditya/webster/internal/handler"
	"github.com/fitraditya/webster/internal/hub"
	"github.com/fitraditya/webster/internal/node"
	"github.com/gorilla/mux"
	"github.com/hashicorp/memberlist"
	"github.com/obrel/go-lib/pkg/log"
	"github.com/spf13/cobra"
)

var (
	nodePort  int
	nodeJoin  string
	address   string
	list      *memberlist.Memberlist
	delegate  *node.Delegate
	serverCmd = &cobra.Command{
		Use:   "server",
		Short: "Run websocket cluster",
	}
)

func serverPreRun(cmd *cobra.Command, args []string) error {
	var err error

	delegate = node.NewDelegate()

	list, err = node.CreateMemberList(cmd.Context(), delegate, nodePort, nodeJoin)
	if err != nil {
		return err
	}

	local := list.LocalNode()
	log.For("server", "prerun").Printf("Node at %s:%d", local.Addr.To4().String(), local.Port)

	return nil
}

func serverRun(cmd *cobra.Command, args []string) error {
	hub := hub.New(list)
	handler := handler.NewHandler(hub)
	server := createServer(handler)
	node := node.NewNode(hub, list)

	go handler.Run()
	go node.Run(delegate)

	go func() {
		log.For("server", "run").Info("Starting websocket server")

		if err := server.ListenAndServe(); err != nil {
			log.For("server", "run").Error(err)
		}
	}()

	if err := handleExitSignal(server, list); err != nil {
		return err
	}

	return nil
}

func createServer(h *handler.Handler) *http.Server {
	router := mux.NewRouter()
	router.HandleFunc("/ws", h.Websocket)
	router.HandleFunc("/", h.Home)

	server := &http.Server{
		Handler: router,
		Addr:    address,
	}

	return server
}

func handleExitSignal(s *http.Server, m *memberlist.Memberlist) error {
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)

	<-stop
	log.For("server", "exit").Info("Got shutdown signal")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	log.For("server", "exit").Info("Leaving cluster")
	if err := m.Leave(time.Second * 5); err != nil {
		return err
	}

	log.For("server", "exit").Info("Shutting down server")
	if err := s.Shutdown(ctx); err != nil {
		return err
	}

	return nil
}

func init() {
	serverCmd.RunE = serverRun
	serverCmd.PreRunE = serverPreRun
	serverCmd.PersistentFlags().StringVarP(&address, "address", "a", "0.0.0.0:4000", "Websocket address")
	serverCmd.PersistentFlags().IntVarP(&nodePort, "port", "p", 2500, "Node port")
	serverCmd.PersistentFlags().StringVarP(&nodeJoin, "join", "j", "", "Node to join")
}
