package api

import (
	"fmt"
	"reflect"
	"sort"

	"github.com/dgf/protzi"
)

// Flow declares a network on registered components.
type Flow interface {
	Add(name, component string)
	Connect(outName, outPort, inName, inPort string)
	In(name, port string) Input
	Out(name, port string) Output
}

type flow struct {
	net        protzi.Network
	add        reflect.Method
	in         reflect.Method
	out        reflect.Method
	components map[string]reflect.Value
}

// Add reflects a new component instance for the network.
func (f *flow) Add(name, component string) {
	ct, ok := registry[component]
	if !ok {
		panic(fmt.Sprintf("component %q not registered", component))
	}

	ci := reflect.New(ct)
	f.components[name] = ci
	f.add.Func.Call([]reflect.Value{reflect.ValueOf(f.net), reflect.ValueOf(name), ci})
}

// Connect connects two instance ports.
func (f *flow) Connect(outName, outPort, inName, inPort string) {
	f.net.Connect(outName+"."+outPort, inName+"."+inPort)
}

func (f *flow) channel(method reflect.Method, name, port string) reflect.Value {
	component, ok := f.components[name]
	if !ok {
		panic(fmt.Sprintf("component %q not found.", name))
	}

	field := component.Elem().FieldByName(port)
	if !field.IsValid() {
		panic(fmt.Sprintf("Component %q port %q not found.", name, port))
	}

	kind := field.Kind()
	if kind != reflect.Chan {
		panic(fmt.Sprintf("Invalid Component %q port %q type %q.", name, port, kind.String()))
	}

	dc := reflect.MakeChan(reflect.ChanOf(reflect.BothDir, field.Type().Elem()), 1)
	method.Func.Call([]reflect.Value{reflect.ValueOf(f.net), reflect.ValueOf(name + "." + port), dc})
	return dc
}

// In creates an input channel.
func (f *flow) In(name, port string) Input {
	return &input{channel: f.channel(f.in, name, port)}
}

// Out creates an output channel.
func (f *flow) Out(name, port string) Output {
	return &output{channel: f.channel(f.out, name, port)}
}

var flows = map[string]Flow{}

// Flows returns flow instance names.
func Flows() []string {
	var f []string
	for name := range flows {
		f = append(f, name)
	}
	sort.Strings(f)
	return f
}

// New creates a fresh flow network.
func New(name string) Flow {
	net := protzi.New(name)
	elem := reflect.TypeOf(net)
	add, _ := elem.MethodByName("Add")
	in, _ := elem.MethodByName("In")
	out, _ := elem.MethodByName("Out")
	f := &flow{
		net:        net,
		add:        add,
		in:         in,
		out:        out,
		components: map[string]reflect.Value{},
	}
	flows[name] = f
	return f
}
