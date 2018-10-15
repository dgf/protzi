package main

import (
	"bufio"
	"os"

	"github.com/dgf/protzi"
	"github.com/dgf/protzi/component"
)

func main() {
	in := make(chan string)
	net := protzi.New("cat")

	net.Add("read", &component.TextFileRead{})
	net.Add("output", &component.Output{})

	net.Connect("read.Text", "output.Message")
	net.Connect("read.Error", "output.Message")

	net.In("read.File", in)
	net.Run()

	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		in <- scanner.Text()
	}
}
