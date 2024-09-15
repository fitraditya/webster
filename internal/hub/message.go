package hub

type IncomingMessage struct {
	Request string      `json:"request"`
	Data    interface{} `json:"data,omitempty"`
}

type OutgoingMessage struct {
	Response string      `json:"response,omitempty"`
	Data     interface{} `json:"data,omitempty"`
}

type EventMessage struct {
	Event  string      `json:"event,omitempty"`
	Sender string      `json:"sender,omitempty"`
	Data   interface{} `json:"data,omitempty"`
}

type ResponseSuccess struct {
	Success bool        `json:"success"`
	Message interface{} `json:"message,omitempty"`
}
