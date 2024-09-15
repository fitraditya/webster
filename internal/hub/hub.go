package hub

import (
	"fmt"
	"sync"

	"github.com/hashicorp/memberlist"
)

type Hub struct {
	mux        sync.Mutex
	clients    map[string]*Client
	memberlist *memberlist.Memberlist
	Connect    chan *Client
	Register   chan *Client
	Unregister chan *Client
}

func New(m *memberlist.Memberlist) *Hub {
	h := &Hub{
		memberlist: m,
		clients:    make(map[string]*Client),
		Connect:    make(chan *Client),
		Register:   make(chan *Client),
		Unregister: make(chan *Client),
	}

	return h
}

func (h *Hub) AddClient(c *Client) error {
	h.mux.Lock()
	defer h.mux.Unlock()

	if _, ok := h.clients[c.id]; ok {
		return fmt.Errorf("Client with id %s is already exists", c.id)
	}

	h.clients[c.id] = c

	return nil
}

func (h *Hub) GetClient(id string) *Client {
	h.mux.Lock()
	defer h.mux.Unlock()

	if _, ok := h.clients[id]; ok {
		return h.clients[id]
	}

	return nil
}

func (h *Hub) GetClients() map[string]*Client {
	h.mux.Lock()
	defer h.mux.Unlock()

	return h.clients
}

func (h *Hub) NumClients() int {
	h.mux.Lock()
	defer h.mux.Unlock()

	return len(h.clients)
}

func (h *Hub) RemoveClient(id string) {
	h.mux.Lock()
	defer h.mux.Unlock()

	delete(h.clients, id)
}

func (h *Hub) Broadcast(msg []byte) {
	h.mux.Lock()
	defer h.mux.Unlock()

	for _, node := range h.memberlist.Members() {
		if node.Name == h.memberlist.LocalNode().Name {
			continue
		}

		h.memberlist.SendReliable(node, msg)
	}

	for client := range h.clients {
		if c, ok := h.clients[client]; ok {
			c.send <- msg
		}
	}
}

func (h *Hub) Gossip(msg []byte) {
	h.mux.Lock()
	defer h.mux.Unlock()

	for client := range h.clients {
		if c, ok := h.clients[client]; ok {
			c.send <- msg
		}
	}
}
