package component

import "fmt"

// Output prints the messages on StdOut.
type Output struct {
	Message <-chan interface{}
}

// Run reads every message and prints it on StdOut.
func (o *Output) Run() {
	for o := range o.Message {
		fmt.Println(o)
	}
}
