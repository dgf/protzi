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

	if err := server.Start(); err != nil {
		log.Fatal("something went wrong:", err.Error())
	}
}
