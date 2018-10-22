package api

import (
	"fmt"
	"reflect"

	"github.com/dgf/protzi"
)

// Flow declares a network on registered components.
type Flow struct {
	name string
	net  protzi.Network
	add  reflect.Method
	in   reflect.Method
	out  reflect.Method
	// TODO mutex maps
	components map[string]reflect.Value
	inputs     map[string]Input
	outputs    map[string]Output
}

func (f *Flow) Name() string {
	return f.name
}

// Add reflects a new component instance for the network.
func (f *Flow) Add(name, component string) error {
	ct, ok := registry[component]
	if !ok {
		return fmt.Errorf("component %q not registered", component)
	}
	ci := reflect.New(ct)
	f.components[name] = ci
	f.add.Func.Call([]reflect.Value{reflect.ValueOf(f.net), reflect.ValueOf(name), ci})
	return nil
}

// Connect connects two instance ports.
func (f *Flow) Connect(outName, outPort, inName, inPort string) error {
	return f.net.Connect(outName+"."+outPort, inName+"."+inPort)
}

func (f *Flow) channel(method reflect.Method, name, port string) (reflect.Value, error) {
	component, ok := f.components[name]
	if !ok {
		return reflect.Value{}, fmt.Errorf("component %q not found", name)
	}

	field := component.Elem().FieldByName(port)
	if !field.IsValid() {
		return reflect.Value{}, fmt.Errorf("component %q port %q not found", name, port)
	}

	kind := field.Kind()
	if kind != reflect.Chan {
		return reflect.Value{}, fmt.Errorf("invalid component %q port %q type %q", name, port, kind.String())
	}

	dc := reflect.MakeChan(reflect.ChanOf(reflect.BothDir, field.Type().Elem()), 1)
	method.Func.Call([]reflect.Value{reflect.ValueOf(f.net), reflect.ValueOf(name + "." + port), dc})
	return dc, nil
}

func (f *Flow) In(name string, port string) (Input, error) {
	cn := name + "." + port
	in, ok := f.inputs[cn]
	if !ok {
		c, err := f.channel(f.in, name, port)
		if err != nil {
			return nil, err
		}
		in = &input{channel: c}
		f.inputs[cn] = in
	}
	return in, nil
}

// Output creates an output channel.
func (f *Flow) Out(name, port string) (Output, error) {
	cn := name + "." + port
	out, ok := f.outputs[cn]
	if !ok {
		c, err := f.channel(f.out, name, port)
		if err != nil {
			return nil, err
		}
		out = &output{channel: c}
		f.outputs[cn] = out
	}
	return out, nil
}

// New creates a fresh flow network.
func New(name string) *Flow {
	net := protzi.New(name)
	elem := reflect.TypeOf(net)
	add, _ := elem.MethodByName("Add")
	in, _ := elem.MethodByName("In")
	out, _ := elem.MethodByName("Out")
	return &Flow{
		net:        net,
		add:        add,
		in:         in,
		out:        out,
		components: map[string]reflect.Value{},
		inputs:     map[string]Input{},
		outputs:    map[string]Output{},
	}
}
