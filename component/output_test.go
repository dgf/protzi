package component_test

import (
	"github.com/dgf/protzi/component"
)

func ExampleOutput() {
	messages := make(chan interface{})

	// create and process
	go (&component.Output{Message: messages}).Run()

	// output
	messages <- "test"

	// Output: test
}
