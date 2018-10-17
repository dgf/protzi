package protzi

import (
	"fmt"
	"log"
	"reflect"
)

// Component that can run
type Component interface {
	Run()
}

// Network that flows.
type Network interface {
	Add(string, Component)
	Connect(out, in string)
	In(string, interface{})
	Out(string, interface{})
	Run()
}

func convertibleTo(out, in reflect.Type) error {
	if !in.ConvertibleTo(out) {
		return fmt.Errorf("Types not convertible %s > %s\n", in, out)
	}
	return nil
}

type network struct {
	name string

	// component by name
	components map[string]Component

	// input channels by component port name
	ins map[string]reflect.Value

	// output channels by component port name
	outs map[string]reflect.Value

	// component output to input connections
	connections map[string][]string
}

// Add a component with a unique name (initializes all unidirectional channels)
func (n *network) Add(name string, c Component) {
	log.Println("Add", name)
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
	cn := out + " > " + in
	log.Println("Connect", cn)

	if op, ok := n.outs[out]; !ok {
		panic(fmt.Sprintf("Input port %q not found", out))
	} else if ip, ok := n.ins[in]; !ok {
		panic(fmt.Sprintf("Output port %q not found", in))
	} else if err := convertibleTo(ip.Type().Elem(), op.Type().Elem()); err != nil {
		panic(err)
	} else {
		n.connections[out] = append(n.connections[out], in)
	}
}

// In maps an input channel
func (n *network) In(in string, c interface{}) {
	valueType := reflect.ValueOf(c).Type().Elem()
	cn := "net > " + in
	log.Println("In", cn, valueType)

	if ip, ok := n.ins[in]; !ok {
		panic(fmt.Sprintf("Input port %q not found", in))
	} else if _, ok := n.connections[cn]; ok {
		panic(fmt.Sprintf("Connected %q before", cn))
	} else if err := convertibleTo(ip.Type().Elem(), valueType); err != nil {
		panic(err)
	} else {
		n.outs[cn] = reflect.ValueOf(c)
		n.connections[cn] = append(n.connections[cn], in)
	}
}

// Out maps an output channel
func (n *network) Out(out string, c interface{}) {
	valueType := reflect.ValueOf(c).Type().Elem()
	cn := out + " > net"
	log.Println("Out", cn, valueType)

	if op, ok := n.outs[out]; !ok {
		panic(fmt.Sprintf("output port %q not found", out))
	} else if _, ok := n.connections[cn]; ok {
		panic(fmt.Sprintf("connected %q before", cn))
	} else if err := convertibleTo(valueType, op.Type().Elem()); err != nil {
		panic(err)
	} else {
		n.ins[cn] = reflect.ValueOf(c)
		n.connections[out] = append(n.connections[out], cn)
	}
}

// Run starts the network by co-routine-ing all connection channels and components
func (n *network) Run() {
	log.Println("Run network", n.name)

	// connections
	for s, o := range n.outs {
		go func(sender string, out reflect.Value) {
			for { // read forever
				if v, ok := out.Recv(); !ok {
					panic(fmt.Sprintf("connection %q closed", sender))
				} else {
					for _, receiver := range n.connections[sender] {
						log.Println("send", sender, ">", receiver)
						n.ins[receiver].Send(v)
					}
				}
			}
		}(s, o)
	}

	// instances
	for _, c := range n.components {
		go func(c Component) {
			c.Run()
		}(c)
	}
}

// New creates a fresh empty network.
func New(name string) Network {
	return &network{
		name:        name,
		components:  map[string]Component{},
		ins:         map[string]reflect.Value{},
		outs:        map[string]reflect.Value{},
		connections: map[string][]string{},
	}
}
