package main

import (
	"flag"
	"log"
	"os"
	"os/signal"

	"github.com/dgf/protzi/api/web"
	"github.com/dgf/protzi/api/web/socket"
)

var addr = flag.String("addr", "localhost:54321", "http service address")

func main() {
	flag.Parse()
	log.SetFlags(0)

	interrupt := make(chan os.Signal)
	signal.Notify(interrupt, os.Interrupt)

	log.Println("connect ", *addr)
	client := socket.Connect(*addr, func(response web.Message) {
		switch response.Type {
		case "error":
			log.Println("error:", string(response.Payload))
		case "output":
			log.Println("output:\n", string(response.Payload))
		case "command":
			log.Println("command:", string(response.Payload))
		case "flow":
			log.Println("flow:\n", string(response.Payload))
		case "flows":
			log.Println("flows:\n", string(response.Payload))
		case "components":
			log.Println("components:\n", string(response.Payload))
		default:
			log.Printf("unknown response type %q", response.Type)
		}
	})

	client.Components()

	flow := "client test"
	client.Flow(flow)
	client.Add(flow, "echo1", "Echo")
	client.Add(flow, "out1", "Print")
	client.Connect(flow, "echo1", "Pong", "out1", "Message")

	client.Flows()

	client.Receive(flow, "out1", "Printed")
	client.Send(flow, "echo1", "Ping", "echo one")
	client.Send(flow, "echo1", "Ping", "echo two")
	client.Send(flow, "echo1", "Ping", "echo three")

	<-interrupt
	log.Println("interrupted")
	client.Interrupt()
}
