package core_test

import "github.com/dgf/protzi/component/core"

func ExamplePrint_Run() {
	messages := make(chan interface{})
	printed := make(chan bool)

	printer := &core.Print{
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
