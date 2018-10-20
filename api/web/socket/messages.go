package socket

import (
	"encoding/json"
	"fmt"
)

type Message struct {
	Type    string          `json:"type"`
	Payload json.RawMessage `json:"payload"`
}

type Command struct {
	Call string          `json:"call"`
	Data json.RawMessage `json:"data"`
}

func Error(err error) *Message {
	return &Message{Type: "error", Payload: json.RawMessage([]byte(fmt.Sprintf(`{"message":%q}`, err.Error())))}
}

type FlowsMessage map[string]string

type ComponentsMessage map[string]string
