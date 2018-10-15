package main

import (
	"bufio"
	"os"

	"github.com/dgf/protzi"
	"github.com/dgf/protzi/component"
)

func main() {
	in := make(chan string)
	net := protzi.New("word count")

	net.Add("read", &component.TextFileRead{})
	net.Add("count", &component.WordCount{})
	net.Add("output", &component.Output{})

	net.Connect("read.Text", "count.Text")
	net.Connect("read.Error", "output.Message")
	net.Connect("count.Counts", "output.Message")

	net.In("read.File", in)
	net.Run()

	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		in <- scanner.Text()
	}
}
