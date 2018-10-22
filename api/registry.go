package api

import (
	"reflect"
	"sort"

	"github.com/dgf/protzi"
)

var (
	registry  = map[string]reflect.Type{}
	inventory = map[string]*Flow{}
)

// Register adds named component.
func Register(name string, c protzi.Component) {
	registry[name] = reflect.TypeOf(c).Elem()
}

// Components returns component names.
func Components() []string {
	var c []string
	for name, rt := range registry {
		c = append(c, name+" > "+rt.String())
	}
	sort.Strings(c)
	return c
}

func Instance(name string) *Flow {
	if i, ok := inventory[name]; ok {
		return i
	}
	i := New(name)
	inventory[name] = i
	return i
}

// Flows returns flow instance names.
func Flows() []string {
	var f []string
	for name := range inventory {
		f = append(f, name)
	}
	sort.Strings(f)
	return f
}
