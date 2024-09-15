package handler

import (
	"net/http"

	"github.com/fitraditya/webster/internal/hub"
	"github.com/google/uuid"
	"github.com/obrel/go-lib/pkg/log"
)

func (h *Handler) Websocket(w http.ResponseWriter, r *http.Request) {
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.For("handler", "websocket").Error(err)
		return
	}

	client := hub.NewClient(uuid.NewString(), h.hub, c)

	go client.Read()
	go client.Write()
}
