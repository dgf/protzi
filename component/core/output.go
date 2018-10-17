package core

import "fmt"

// Output prints the messages on StdOut.
type Output struct {
	Message <-chan interface{}
	Printed chan<- bool
}

// Run reads every message and prints it on StdOut.
func (o *Output) Run() {
	for m := range o.Message {
		fmt.Println(m)
		o.Printed <- true
	}
}
