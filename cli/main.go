package main

import (
	"fmt"
	"log"

	"github.com/dgf/protzi/api"
	"github.com/dgf/protzi/component/core"
	"github.com/dgf/protzi/component/text"
)

func main() {

	// register available components
	api.Register("Echo", &core.Echo{})
	api.Register("Print", &core.Print{})
	api.Register("Tick", &core.Tick{})
	api.Register("Time", &core.Time{})
	api.Register("Read", &text.FileRead{})
	api.Register("Render", &text.Render{})
	api.Register("WordCount", &text.WordCount{})

	fmt.Println("registry:")
	for _, component := range api.Components() {
		fmt.Println(component)
	}

	echoFlow := api.New("echo")
	api.New("cat")
	api.New("wc")

	fmt.Println("flows:")
	fmt.Println(api.Flows())

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
