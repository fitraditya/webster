package node

import (
	"github.com/fitraditya/webster/internal/hub"
	"github.com/hashicorp/memberlist"
	"github.com/obrel/go-lib/pkg/log"
)

type Node struct {
	hub        *hub.Hub
	memberlist *memberlist.Memberlist
	Broadcast  chan []byte
}

func NewNode(h *hub.Hub, m *memberlist.Memberlist) *Node {
	n := &Node{
		hub:        h,
		memberlist: m,
		Broadcast:  make(chan []byte),
	}

	return n
}

func (n *Node) Run(d *Delegate) {
	run := true

	for run {
		data := <-d.messages
		n.hub.Gossip(data)

		log.For("node", "run").Infof("Received msg: %s", string(data))
	}
}
