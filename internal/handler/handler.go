package handler

import (
	"github.com/fitraditya/webster/internal/hub"
	"github.com/gorilla/websocket"
	"github.com/obrel/go-lib/pkg/log"
)

var upgrader = websocket.Upgrader{}

type Handler struct {
	hub *hub.Hub
}

func NewHandler(n *hub.Hub) *Handler {
	return &Handler{
		hub: n,
	}
}

func (h *Handler) Run() {
	for {
		select {
		case client := <-h.hub.Connect:
			log.For("handler", "run").Infof("%s connected\n", client.ID())
		case client := <-h.hub.Register:
			h.hub.AddClient(client)
			log.For("handler", "run").Infof("Client %s created\n", client.ID())
		case client := <-h.hub.Unregister:
			if c := h.hub.GetClient(client.ID()); c != nil {
				c.Close()
				h.hub.RemoveClient(c.ID())
				log.For("handler", "run").Infof("Client %s removed\n", c.ID())
			}
		}
	}
}
