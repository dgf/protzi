package main

import (
	"bufio"
	"os"

	"github.com/dgf/protzi"
	"github.com/dgf/protzi/component"
)

func main() {
	in := make(chan string)

	net := protzi.New("display")
	net.Add("output", &component.Output{})
	net.In("output.Message", in)
	net.Run()

	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		in <- scanner.Text()
	}
}
