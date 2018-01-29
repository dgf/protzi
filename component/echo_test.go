package component_test

import (
	"fmt"

	"github.com/dgf/protzi/component"
)

func ExampleEcho() {
	ping := make(chan interface{})
	pong := make(chan interface{})

	// create and process
	go (&component.Echo{Ping: ping, Pong: pong}).Run()

	// ping pong
	ping <- "echo"
	fmt.Println(<-pong)

	// Output: echo
}
