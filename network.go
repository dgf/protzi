package protzi

import (
	"fmt"
	"reflect"
)

// Network thats flows.
type Network interface {
	Add(string, Component)
	Connect(out, in string)
	In(string, interface{})
	Out(string, interface{})
	Run()
}

type connection struct {
	// out > in
	name string
	// input channel (out port)
	in reflect.Value
	// output channel (in port)
	out reflect.Value
}

func (c *connection) run() {
	fmt.Println("run connection", c.name)
	ot := c.out.Type().Elem()
	for { // loop forever
		if v, ok := c.in.Recv(); !ok {
			panic(fmt.Sprintf("connection %q closed", c.name))
		} else {
			fmt.Println("send", c.name, ot, v.Type())

			// convert?
			if c.out.Type().Elem() != v.Type() {
				panic(fmt.Sprintf("Type differs %s != %s\n", c.out.Type().Elem(), v.Type()))
			}

			c.out.Send(v)
		}
	}
}

type network struct {
	connections map[string]connection
	components  map[string]Component
	// input channels by component port name
	ins map[string]reflect.Value
	// output channels by component port name
	outs map[string]reflect.Value
}

// New creates a fresh empty network.
func New(i interface{}) Network {
	return &network{
		connections: map[string]connection{},
		components:  map[string]Component{},
		ins:         map[string]reflect.Value{},
		outs:        map[string]reflect.Value{},
	}
}

func (n *network) Add(name string, c Component) {
	fmt.Println("add", name)
	n.components[name] = c

	// get underlying value and type
	v := reflect.ValueOf(c).Elem()
	t := v.Type()

	// loop type fields
	for i := 0; i < t.NumField(); i++ {
		f := t.Field(i)

		// channel?
		if f.Type.Kind() == reflect.Chan {
			if f.Type.ChanDir() == reflect.BothDir {
				continue // ignore this
			}

			// create undirected channel
			ct := reflect.ChanOf(reflect.BothDir, f.Type.Elem())
			dc := reflect.MakeChan(ct, 1)

			// assign channel
			v.Field(i).Set(dc)

			// remember channel by port name
			cn := fmt.Sprintf("%s.%s", name, f.Name)
			if f.Type.ChanDir() == reflect.RecvDir {
				n.ins[cn] = dc
			} else {
				n.outs[cn] = dc
			}
		}
	}
}

func (n *network) Connect(out, in string) {
	fmt.Println("connect", out, in)
	if op, ok := n.outs[out]; !ok {
		panic(fmt.Sprintf("input port %q not found", out))
	} else if ip, ok := n.ins[in]; !ok {
		panic(fmt.Sprintf("output port %q not found", in))
	} else {
		cn := out + " > " + in
		if _, ok := n.connections[cn]; ok {
			panic(fmt.Sprintf("connected %q before", cn))
		}
		n.connections[cn] = connection{name: cn, in: op, out: ip}
	}
}

func (n *network) Run() {
	fmt.Println("run net")
	for _, c := range n.connections {
		go func(c connection) { c.run() }(c)
	}
	for _, c := range n.components {
		go func(c Component) { c.Run() }(c)
	}
}

func (n *network) In(in string, c interface{}) {
	fmt.Println("in", in)
	if ip, ok := n.ins[in]; !ok {
		panic(fmt.Sprintf("input port %q not found", in))
	} else {
		// TODO validate channel
		cn := "net > " + in
		if _, ok := n.connections[cn]; ok {
			panic(fmt.Sprintf("connected %q before", cn))
		}
		n.connections[cn] = connection{name: cn, in: reflect.ValueOf(c), out: ip}
	}
}

func (n *network) Out(out string, c interface{}) {
	fmt.Println("out", out)
	if op, ok := n.outs[out]; !ok {
		panic(fmt.Sprintf("output port %q not found", out))
	} else {
		// TODO validate channel
		cn := out + " > net"
		if _, ok := n.connections[cn]; ok {
			panic(fmt.Sprintf("connected %q before", cn))
		}
		n.connections[cn] = connection{name: cn, in: op, out: reflect.ValueOf(c)}
	}
}
