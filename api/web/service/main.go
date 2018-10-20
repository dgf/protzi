package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/dgf/protzi/api"
	"github.com/dgf/protzi/api/web/socket"
)

var addr string

func init() {
	flag.StringVar(&addr, "addr", ":54321", "TCP address to listen on")
}

func main() {
	flag.Parse()

	server := socket.Serve(addr)
	log.Printf("Running on: " + server.Endpoint())

	fmt.Println("registry:")
	for _, component := range api.Components() {
		fmt.Println(component)
	}

	// flow an echo to print
	echoFlow := api.New("Echo example")
	echoFlow.Add("echo1", "Echo")
	echoFlow.Add("out1", "Print")
	echoFlow.Connect("echo1", "Pong", "out1", "Message")
	in := echoFlow.In("echo1", "Ping")
	out := echoFlow.Out("out1", "Printed")

	fmt.Println("flows:")
	fmt.Println(api.Flows())

	for _, ping := range []string{"one", "two"} {
		go in.Send(ping)
		_, ok := out.Receive()
		if !ok {
			log.Fatal(ping, "failed!")
		}
	}

	if err := server.Start(); err != nil {
		log.Fatal("something went wrong:", err.Error())
	}
}
