package component_test

import "github.com/dgf/protzi/component"

func ExampleOutput() {
	messages := make(chan interface{})
	done := make(chan bool)

	go func() {
		(&component.Output{Message: messages}).Run()
		done <- true
	}()

	messages <- "one"
	messages <- "two"
	messages <- "three"

	close(messages)
	<-done

	// Output:
	// one
	// two
	// three
}
