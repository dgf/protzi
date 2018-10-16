package component_test

import "github.com/dgf/protzi/component"

func ExampleOutput() {
	messages := make(chan interface{})
	printed := make(chan bool)

	printer := &component.Output{
		Message: messages,
		Printed: printed,
	}
	go printer.Run()

	for _, m := range []string{"one", "two", "three"} {
		messages <- m
		<-printed
	}

	// Output:
	// one
	// two
	// three
}
