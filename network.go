package protzi

import (
	"fmt"
	"log"
	"reflect"
)

// Network that flows.
type Network interface {
	Add(string, Component)
	Connect(out, in string)
	In(string, interface{})
	Out(string, interface{})
	Run()
}

type network struct {
	name string

	// component by name
	components map[string]Component

	// input channels by component port name
	ins map[string]reflect.Value

	// output channels by component port name
	outs map[string]reflect.Value

	// component connections by combined name
	connections map[string]connection
}

// New creates a fresh empty network.
func New(name string) Network {
	return &network{
		name:        name,
		components:  map[string]Component{},
		ins:         map[string]reflect.Value{},
		outs:        map[string]reflect.Value{},
		connections: map[string]connection{},
	}
}

// Add a component with a unique name (initializes all unidirectional channels)
func (n *network) Add(name string, c Component) {
	log.Println("add", name)
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

// Connect two channels from the output of one to the input of another.
func (n *network) Connect(out, in string) {
	log.Println("connect", out, in)

	if op, ok := n.outs[out]; !ok {
		panic(fmt.Sprintf("input port %q not found", out))
	} else if ip, ok := n.ins[in]; !ok {
		panic(fmt.Sprintf("output port %q not found", in))
	} else {
		cn := out + " > " + in
		if _, ok := n.connections[cn]; ok {
			panic(fmt.Sprintf("connected %q before", cn))
		}

		outTypeElem := ip.Type().Elem()
		inTypeElem := op.Type().Elem()
		if !inTypeElem.ConvertibleTo(outTypeElem) {
			panic(fmt.Sprintf("Type differs %s != %s\n", inTypeElem, outTypeElem))
		}

		n.connections[cn] = connection{name: cn, in: op, out: ip}
	}
}

// Run starts the network by co-routine-ing all connection channels and components
func (n *network) Run() {
	log.Println("run net")
	for _, c := range n.connections {
		go func(c connection) { c.run() }(c)
	}
	for _, c := range n.components {
		go func(c Component) { c.Run() }(c)
	}
}

// In maps an input channel
func (n *network) In(in string, c interface{}) {
	log.Println("in", in)
	if ip, ok := n.ins[in]; !ok {
		panic(fmt.Sprintf("input port %q not found", in))
	} else {
		cn := "net > " + in
		if _, ok := n.connections[cn]; ok {
			panic(fmt.Sprintf("Connected %q before", cn))
		}

		outTypeElem := ip.Type().Elem()
		valueOfChannel := reflect.ValueOf(c)
		inTypeElem := valueOfChannel.Type().Elem()
		if !inTypeElem.ConvertibleTo(outTypeElem) {
			panic(fmt.Sprintf("Type differs %s > %s\n", inTypeElem, outTypeElem))
		}

		n.connections[cn] = connection{name: cn, in: valueOfChannel, out: ip}
	}
}

// Out maps an output channel
func (n *network) Out(out string, c interface{}) {
	log.Println("out", out)
	if op, ok := n.outs[out]; !ok {
		panic(fmt.Sprintf("output port %q not found", out))
	} else {
		cn := out + " > net"
		if _, ok := n.connections[cn]; ok {
			panic(fmt.Sprintf("connected %q before", cn))
		}

		inTypeElem := op.Type().Elem()
		valueOfChannel := reflect.ValueOf(c)
		outTypeElem := valueOfChannel.Type().Elem()
		if !inTypeElem.ConvertibleTo(outTypeElem) {
			panic(fmt.Sprintf("Type differs %s > %s\n", inTypeElem, outTypeElem))
		}

		n.connections[cn] = connection{name: cn, in: op, out: reflect.ValueOf(c)}
	}
}
