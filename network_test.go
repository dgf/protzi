package protzi_test

import (
	"fmt"
	"io/ioutil"
	"os"
	"sort"
	"testing"
	"time"

	"github.com/dgf/protzi"
	"github.com/dgf/protzi/component/core"
	"github.com/dgf/protzi/component/text"
)

func ExampleNetwork_echo() {
	in := make(chan interface{})
	out := make(chan interface{})

	net := protzi.New("passthru")
	net.Add("echo", &core.Echo{})

	if err := net.In("echo.Ping", in); err != nil {
		fmt.Println(err)
	} else if err := net.Out("echo.Pong", out); err != nil {
		fmt.Println(err)
	}

	in <- "one"
	fmt.Println(<-out)
	in <- "two"
	fmt.Println(<-out)
	// Output:
	// one
	// two
}

func ExampleNetwork_split() {
	in := make(chan interface{})
	out1 := make(chan interface{})
	out2 := make(chan interface{})

	net := protzi.New("split")
	net.Add("echoIn", &core.Echo{})
	net.Add("echoOut1", &core.Echo{})
	net.Add("echoOut2", &core.Echo{})

	if err := net.In("echoIn.Ping", in); err != nil {
		fmt.Println(err)
	} else if err := net.Connect("echoIn.Pong", "echoOut1.Ping"); err != nil {
		fmt.Println(err)
	} else if err := net.Connect("echoIn.Pong", "echoOut2.Ping"); err != nil {
		fmt.Println(err)
	} else if err := net.Out("echoOut1.Pong", out1); err != nil {
		fmt.Println(err)
	} else if err := net.Out("echoOut2.Pong", out2); err != nil {
		fmt.Println(err)
	}

	in <- "echo"
	twice := 0
	for twice != 2 {
		select {
		case o := <-out1:
			fmt.Println(o)
			twice++
		case o := <-out2:
			fmt.Println(o)
			twice++
		}
	}
	// Output:
	// echo
	// echo
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
	network.Add("read", &text.FileRead{})
	network.Add("count", &text.WordCount{})

	// connect component
	if err := network.In("read.File", in); err != nil {
		fmt.Println(err)
	} else if err := network.Connect("read.Text", "count.Text"); err != nil {
		fmt.Println(err)
	} else if err := network.Out("count.Counts", out); err != nil {
		fmt.Println(err)
	}

	// process file
	in <- file.Name()
	countsByWord := <-out

	// stringify and sort word counts (needed for output assertion)
	var wordCounts []string
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
	network := protzi.New("valid read text to output interface")
	network.Add("read", &text.FileRead{})
	network.Add("out", &core.Print{})
	if err := network.Connect("read.Text", "out.Message"); err != nil {
		t.Error(err)
	}
}

func TestNetwork_Connect_invalidPanic(t *testing.T) {
	network := protzi.New("invalid count output map to input text")
	network.Add("out", &text.WordCount{})
	network.Add("in", &text.WordCount{})

	if err := network.Connect("out.Counts", "in.Text"); err == nil {
		t.Errorf("should fail with invalid type mapping")
	}
}

func TestNetwork_Init_timerTwice(t *testing.T) {
	stamps := make(chan time.Time)

	network := protzi.New("endless initialized timer")
	network.Add("timer", &core.Time{})
	if err := network.Out("timer.Stamp", stamps); err != nil {
		t.Error(err)
	} else if err := network.Init("timer.Duration", time.Nanosecond); err != nil {
		t.Error(err)
	}

	count := 0
	for count < 2 {
		select {
		case <-stamps:
			count++
		case <-time.After(1 * time.Millisecond):
			t.Error("Timed out!")
		}
	}
}
