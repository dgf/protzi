package protzi

import (
	"fmt"
	"log"
	"reflect"
)

type connection struct {
	// out > in
	name string

	// input channel (out port)
	in reflect.Value

	// output channel (in port)
	out reflect.Value
}

func (c *connection) run() {
	log.Println("run connection", c.name)
	for { // read forever
		if v, ok := c.in.Recv(); !ok {
			panic(fmt.Sprintf("connection %q closed", c.name))
		} else {
			log.Println("send", c.name, v)
			c.out.Send(v)
		}
	}
}
