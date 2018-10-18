package api

import "reflect"

// Output receives payloads.
type Output interface {
	Receive() (payload interface{}, ok bool)
}

type output struct {
	channel reflect.Value
}

// Receive awaits payload of the output channel.
func (o *output) Receive() (interface{}, bool) {
	return o.channel.Recv()
}
