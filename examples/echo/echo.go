package main

import (
	"bufio"
	"os"

	"github.com/dgf/protzi"
	"github.com/dgf/protzi/component/core"
)

func main() {
	in := make(chan string)

	net := protzi.New("display")
	net.Add("output", &core.Print{})
	net.In("output.Message", in)

	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		in <- scanner.Text()
	}
}
