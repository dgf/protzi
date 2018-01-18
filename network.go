package protzi

import (
	"fmt"
	"reflect"
)

// Network thats flows
type Network interface {
	Connect(out, in string)
	Add(string, Component)
	In(string)
	Out(string)
	Run()
	Process(interface{}) interface{}
}

type network struct {
	name       string
	start      reflect.Value
	end        reflect.Value
	components map[string]Component
	ins        map[string]reflect.Value
	outs       map[string]reflect.Value
}

// New creates a fresh empty network.
func New(name string) Network {
	return &network{
		name:       name,
		components: map[string]Component{},
		ins:        map[string]reflect.Value{},
		outs:       map[string]reflect.Value{},
	}
}

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
			}
		}
	}
}

func (n *network) Connect(out, in string) {
	if op, ok := n.outs[out]; !ok {
		panic(fmt.Sprintf("port %q not found", out))
	} else if ip, ok := n.ins[in]; !ok {
		panic(fmt.Sprintf("port %q not found", in))
	} else {
		go func() {
			if v, ok := op.Recv(); !ok {
				panic(fmt.Sprintf("port %q closed", out))
			} else {
				ip.Send(v)
			}
		}()
	}
}

func (n *network) Run() {
	for _, c := range n.components {
		go c.Run()
	}
}

func (n *network) Process(i interface{}) interface{} {
	n.start.Send(reflect.ValueOf(i))
	if v, ok := n.end.Recv(); !ok {
		panic("out port is closed")
	} else {
		return v.Interface()
	}
}

func (n *network) In(in string) {
	if ip, ok := n.ins[in]; !ok {
		panic(fmt.Sprintf("port %q not found", in))
	} else {
		n.start = ip
	}
}

func (n *network) Out(out string) {
	if op, ok := n.outs[out]; !ok {
		panic(fmt.Sprintf("port %q not found", out))
	} else {
		n.end = op
	}
}
