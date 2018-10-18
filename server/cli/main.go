package main

import (
	"fmt"
	"log"

	"github.com/dgf/protzi/component/core"
	"github.com/dgf/protzi/component/text"
	"github.com/dgf/protzi/server"
)

func main() {

	// register available components
	server.Register("Echo", &core.Echo{})
	server.Register("Print", &core.Print{})
	server.Register("Ticker", &core.Ticker{})
	server.Register("Timer", &core.Timer{})
	server.Register("Read", &text.FileRead{})
	server.Register("Render", &text.Render{})
	server.Register("WordCount", &text.WordCount{})

	fmt.Println("registry:")
	for _, component := range server.Components() {
		fmt.Println(component)
	}

	echoFlow := server.New("echo")
	server.New("cat")
	server.New("wc")

	fmt.Println("flows:")
	fmt.Println(server.Flows())

	// flow an echo to print
	echoFlow.Add("echo1", "Echo")
	echoFlow.Add("out1", "Print")
	echoFlow.Connect("echo1", "Pong", "out1", "Message")
	in := echoFlow.In("echo1", "Ping")
	out := echoFlow.Out("out1", "Printed")

	for _, ping := range []string{"one", "two"} {
		go in.Send(ping)
		_, ok := out.Receive()
		if !ok {
			log.Println(ping, "failed!")
		}
	}
}
