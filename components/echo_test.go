package components_test

import (
	"fmt"

	"github.com/dgf/protzi/components"
)

func ExampleEcho() {
	ping := make(chan interface{})
	pong := make(chan interface{})

	// create and run process
	go (&components.Echo{Ping: ping, Pong: pong}).Run()

	// ping pong
	ping <- "echo"
	fmt.Println(<-pong)

	// Output: echo
}
