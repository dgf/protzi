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
	ot := c.out.Type().Elem()
	for { // read forever
		if v, ok := c.in.Recv(); !ok {
			panic(fmt.Sprintf("connection %q closed", c.name))
		} else {
			log.Println("send", c.name, ot, v.Type())

			// convert?
			if c.out.Type().Elem() != v.Type() {
				panic(fmt.Sprintf("Type differs %s != %s\n", c.out.Type().Elem(), v.Type()))
			}

			c.out.Send(v)
		}
	}
}
