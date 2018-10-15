package protzi_test

import (
	"fmt"
	"io/ioutil"
	"os"
	"sort"
	"testing"

	"github.com/dgf/protzi"
	"github.com/dgf/protzi/component"
)

func ExampleNetwork_echo() {
	in := make(chan interface{})
	out := make(chan interface{})

	net := protzi.New("passthru")
	net.Add("echo", &component.Echo{})
	net.In("echo.Ping", in)
	net.Out("echo.Pong", out)
	net.Run()

	in <- "echo"
	fmt.Println(<-out)
	// Output: echo
}

func ExampleNetwork_fileWordCounter() {

	// create temp file
	file, err := ioutil.TempFile(os.TempDir(), "example")
	if err != nil {
		panic(err)
	}

	// write test data string
	data := "\f\n\r\ttwo\fone\ntwo\r\t"
	if _, err := file.Write([]byte(data)); err != nil {
		panic(err)
	}

	// close temp file
	if err := file.Close(); err != nil {
		panic(err)
	}

	// create network
	in := make(chan string)
	out := make(chan map[string]int)

	network := protzi.New("file word counter")

	// add process component
	network.Add("read", &component.TextFileRead{})
	network.Add("count", &component.WordCount{})

	// connect component
	network.In("read.File", in)
	network.Connect("read.Text", "count.Text")
	network.Out("count.Counts", out)

	// run it
	network.Run()

	// process file
	in <- file.Name()
	countsByWord := <-out

	// stringify and sort word counts (needed for output assertion)
	wordCounts := []string{}
	for word := range countsByWord {
		wordCounts = append(wordCounts, fmt.Sprintf("%s: %d", word, countsByWord[word]))
	}
	sort.Strings(wordCounts)

	// print word counts
	fmt.Println(wordCounts)

	// delete temporary file
	if err := os.Remove(file.Name()); err != nil {
		panic(err)
	}
	// Output: [one: 1 two: 2]
}

func TestNetwork_Connect_valid(t *testing.T) {
	network := protzi.New("count out to in")
	network.Add("read", &component.TextFileRead{})
	network.Add("out", &component.Output{})
	network.Connect("read.Text", "out.Message")
}

func TestNetwork_Connect_invalidPanic(t *testing.T) {
	network := protzi.New("count out to in")
	network.Add("countOut", &component.WordCount{})
	network.Add("countIn", &component.WordCount{})

	defer func() {
		if r := recover(); r == nil {
			t.Errorf("should panic with invalid type")
		}
	}()
	network.Connect("countOut.Counts", "countIn.Text")
}
