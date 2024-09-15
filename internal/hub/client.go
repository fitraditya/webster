package hub

import (
	"encoding/json"

	"github.com/gorilla/websocket"
	"github.com/obrel/go-lib/pkg/log"
)

type Client struct {
	id     string
	hub    *Hub
	socket *websocket.Conn
	send   chan []byte
}

func NewClient(id string, hub *Hub, conn *websocket.Conn) *Client {
	return &Client{
		id:     id,
		hub:    hub,
		socket: conn,
		send:   make(chan []byte),
	}
}

func (c *Client) ID() string {
	return c.id
}

func (c *Client) Send() chan []byte {
	return c.send
}

func (c *Client) Close() {
	close(c.send)
}

func (c *Client) Read() {
	defer func() {
		c.socket.Close()

		if c.id != "" {
			c.hub.Unregister <- c
		}
	}()

	run := true

	for run {
		_, message, err := c.socket.ReadMessage()
		if err != nil {
			log.For("client", "read").Error(err)
			c.socket.Close()

			if c.id != "" {
				c.hub.Unregister <- c
			}

			break
		}

		msg := &IncomingMessage{}

		err = json.Unmarshal([]byte(message), msg)
		if err != nil {
			log.For("client", "read").Error(err)
			return
		}

		switch msg.Request {
		case "register":
			log.For("client", "read").Debugf("Registering %s\n", c.id)

			c.hub.Register <- c

			jsonMessage, _ := json.Marshal(&OutgoingMessage{
				Response: "register",
				Data: &ResponseSuccess{
					Success: true,
				},
			})

			c.send <- jsonMessage
		case "broadcast":
			log.For("client", "read").Infof("Broadcast from %s\n", c.id)

			message, _ := json.Marshal(&EventMessage{
				Event:  "broadcast",
				Sender: c.id,
				Data:   msg.Data,
			})
			c.hub.Broadcast(message)

			jsonMessage, _ := json.Marshal(&OutgoingMessage{
				Response: "broadcast",
				Data: &ResponseSuccess{
					Success: true,
				},
			})

			c.send <- jsonMessage
		}
	}
}

func (c *Client) Write() {
	defer func() {
		c.socket.Close()
	}()

	run := true

	for run {
		message, ok := <-c.send

		if !ok {
			c.socket.WriteMessage(websocket.CloseMessage, []byte{})
			return
		}

		c.socket.WriteMessage(websocket.TextMessage, message)
	}
}
