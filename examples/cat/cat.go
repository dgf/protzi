package main

import (
	"bufio"
	"fmt"
	"os"

	"github.com/dgf/protzi"
	"github.com/dgf/protzi/component/core"
	"github.com/dgf/protzi/component/text"
)

// go run cat.go
// go run cat.go /etc/issue /unknown/file
// echo -e "one\ntwo" | go run cat.go
// echo -e "/etc/issue\n/unknown/file" | xargs go run cat.go
func main() {

	// pipe or char device?
	if info, err := os.Stdin.Stat(); err != nil {
		panic(err)
	} else if len(os.Args) == 1 && info.Mode()&(os.ModeNamedPipe|os.ModeCharDevice) != 0 {
		scanner := bufio.NewScanner(os.Stdin)
		for scanner.Scan() {
			fmt.Println(scanner.Text())
		}
		return
	}

	// create network
	net := protzi.New("cat")
	net.Add("read", &text.FileRead{})
	net.Add("output", &core.Print{})

	// bind channels
	in := make(chan string)
	out := make(chan bool)
	net.In("read.File", in)
	net.Out("output.Printed", out)

	// connect and run it
	net.Connect("read.Text", "output.Message")
	net.Connect("read.Error", "output.Message")
	net.Run()

	// flow the arguments
	for _, arg := range os.Args[1:] {
		in <- arg
		<-out
	}
}
