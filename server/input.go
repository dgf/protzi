package server

import "reflect"

// Input sends payloads.
type Input interface {
	Send(interface{})
}

type input struct {
	channel reflect.Value
}

// Send forwards the payload into the input channel.
func (i *input) Send(payload interface{}) {
	i.channel.Send(reflect.ValueOf(payload))
}
