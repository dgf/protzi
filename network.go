package protzi

import (
	"fmt"
	"reflect"
)

// Component that can run
type Component interface {
	Run()
}

// Network that flows.
type Network interface {
	Add(string, Component)
	Connect(out, in string) error
	In(string, interface{}) error
	Init(string, interface{}) error
	Out(string, interface{}) error
}

func convertibleTo(out, in reflect.Type) error {
	if !in.ConvertibleTo(out) {
		return fmt.Errorf("types not convertible (%s > %s)", in, out)
	}
	return nil
}

type network struct {
	name string
	// TODO sync.Map or Mutex

	// component by name
	components map[string]Component

	// input channels by component port name
	ins map[string]reflect.Value

	// output channels by component port name
	outs map[string]reflect.Value

	// component output to input connections
	connections map[string][]string
}

// listen on output channel to forward payloads to every connected input
func forward(n *network, sender string, out reflect.Value) error {
	for { // read forever TODO use context to cancel
		if v, ok := out.Recv(); !ok {
			return fmt.Errorf("connection %q closed", sender)
		} else {
			for _, receiver := range n.connections[sender] {
				n.ins[receiver].Send(v)
			}
		}
	}
	return nil
}

// Add a component with a unique name (initializes all unidirectional channels)
func (n *network) Add(name string, c Component) {
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
				go forward(n, cn, dc)
			}
		}
	}

	// run component loop
	go c.Run()
}

// Connect two channels from the output of one to the input of another.
func (n *network) Connect(out, in string) error {
	if op, ok := n.outs[out]; !ok {
		return fmt.Errorf("output port %q not found", out)
	} else if ip, ok := n.ins[in]; !ok {
		return fmt.Errorf("input port %q not found", in)
	} else if err := convertibleTo(ip.Type().Elem(), op.Type().Elem()); err != nil {
		return err
	} else {
		n.connections[out] = append(n.connections[out], in)
	}
	return nil
}

// In maps an input channel
func (n *network) In(in string, c interface{}) error {
	dc := reflect.ValueOf(c)
	valueType := dc.Type().Elem()
	cn := "net > " + in

	if ip, ok := n.ins[in]; !ok {
		return fmt.Errorf("input port %q not found", in)
	} else if err := convertibleTo(ip.Type().Elem(), valueType); err != nil {
		return err
	} else {
		n.outs[cn] = dc
		n.connections[cn] = append(n.connections[cn], in)
		go forward(n, cn, dc)
	}
	return nil
}

// Init infinitely forwards the payload to the input channel.
func (n *network) Init(in string, p interface{}) error {
	value := reflect.ValueOf(p)
	if ip, ok := n.ins[in]; !ok {
		return fmt.Errorf("init port %q not found", in)
	} else if err := convertibleTo(ip.Type().Elem(), value.Type()); err != nil {
		return err
	} else {
		go func(port reflect.Value) {
			for { // forever TODO cancel with context
				port.Send(value)
			}
		}(ip)
	}
	return nil
}

// Out maps an output channel
func (n *network) Out(out string, c interface{}) error {
	valueType := reflect.ValueOf(c).Type().Elem()
	cn := out + " > net"

	if op, ok := n.outs[out]; !ok {
		return fmt.Errorf("output port %q not found", out)
	} else if err := convertibleTo(valueType, op.Type().Elem()); err != nil {
		return err
	} else {
		n.ins[cn] = reflect.ValueOf(c)
		n.connections[out] = append(n.connections[out], cn)
	}
	return nil
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
