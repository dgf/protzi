package component

import "fmt"

// Output prints the message on stdout
type Output struct {
	Message <-chan interface{}
}

func (o *Output) Run() {
	for o := range o.Message {
		fmt.Println(o)
	}
}
