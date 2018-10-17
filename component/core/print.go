package core

import "fmt"

// Print prints the messages on Stdout.
type Print struct {
	Message <-chan interface{}
	Printed chan<- bool
}

// Run reads every message and prints it on Stdout.
func (o *Print) Run() {
	for m := range o.Message {
		fmt.Println(m)
		o.Printed <- true
	}
}
