package server

import (
	"reflect"
	"sort"

	"github.com/dgf/protzi"
)

var registry = map[string]reflect.Type{}

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
