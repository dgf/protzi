package main

import (
	"bufio"
	"os"

	"github.com/dgf/protzi"
	"github.com/dgf/protzi/component/core"
	"github.com/dgf/protzi/component/text"
)

func main() {
	in := make(chan string)
	net := protzi.New("word count")

	net.Add("read", &text.FileRead{})
	net.Add("count", &text.WordCount{})
	net.Add("output", &core.Output{})

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
