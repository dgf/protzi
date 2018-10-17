package main

import (
	"os"

	"github.com/dgf/protzi"
	"github.com/dgf/protzi/component/core"
	"github.com/dgf/protzi/component/text"
)

// go run wc.go /etc/group /unknown/file
func main() {
	if len(os.Args) == 1 {
		return
	}

	// create
	net := protzi.New("word count")
	net.Add("read", &text.FileRead{})
	net.Add("count", &text.WordCount{})
	net.Add("output", &core.Print{})

	// bind
	in := make(chan string)
	out := make(chan bool)
	net.In("read.File", in)
	net.Out("output.Printed", out)

	// connect and run
	net.Connect("read.Text", "count.Text")
	net.Connect("read.Error", "output.Message")
	net.Connect("count.Counts", "output.Message")
	net.Run()

	// flow the arguments
	for _, arg := range os.Args[1:] {
		in <- arg
		<-out
	}
}
