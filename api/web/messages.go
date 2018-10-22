package web

import (
	"encoding/json"
	"fmt"

	"github.com/dgf/protzi/api"
)

type Message struct {
	Type    string          `json:"type"`
	Payload json.RawMessage `json:"payload"`
}

type ErrorPayload struct {
	Message json.RawMessage `json:"message"`
}

type FlowPayload struct {
	Name string `json:"name"`
}

type OutputPayload struct {
	OutputData
	Value json.RawMessage `json:"value"`
}

func ErrorMessage(err error) *Message {
	return &Message{
		Type:    "error",
		Payload: json.RawMessage([]byte(fmt.Sprintf(`{"message":%q}`, err.Error()))),
	}
}

func CommandMessage(c Command) *Message {
	payload, err := json.Marshal(c)
	if err != nil {
		panic(err)
	}
	return &Message{
		Type:    "command",
		Payload: payload,
	}
}

func FlowMessage(i *api.Flow) *Message {
	return &Message{
		Type:    "flow",
		Payload: json.RawMessage([]byte(fmt.Sprintf(`{"name":%q}`, i.Name()))),
	}
}
