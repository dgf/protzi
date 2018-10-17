package core_test

import (
	"fmt"

	"github.com/dgf/protzi/component/core"
)

func ExampleEcho_Run() {
	ping := make(chan interface{})
	pong := make(chan interface{})

	echoer := &core.Echo{Ping: ping, Pong: pong}
	go echoer.Run()

	ping <- "echo"
	fmt.Println(<-pong)
	// Output: echo
}
