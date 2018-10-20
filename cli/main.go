package main

import (
	"flag"
	"log"
	"os"
	"os/signal"

	"github.com/dgf/protzi/api/web/socket"
)

var addr = flag.String("addr", "localhost:54321", "http service address")

func main() {
	flag.Parse()
	log.SetFlags(0)

	interrupt := make(chan os.Signal)
	signal.Notify(interrupt, os.Interrupt)

	log.Println("connect ", *addr)
	client := socket.Connect(*addr, func(response socket.Message) {
		switch response.Type {
		case "error":
			log.Println("error:", string(response.Payload))
		case "flows":
			log.Println("flows:\n", string(response.Payload))
		case "components":
			log.Println("components:\n", string(response.Payload))
		default:
			log.Printf("unknown response type %q", response.Type)
		}
	})

	client.Send(socket.Command{Call: "flows"})
	client.Send(socket.Command{Call: "unknown"})
	client.Send(socket.Command{Call: "components"})

	<-interrupt
	log.Println("interrupted")
	client.Interrupt()
}
