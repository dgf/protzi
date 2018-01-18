package protzi_test

import (
	"fmt"
	"io/ioutil"
	"os"
	"sort"

	"github.com/dgf/protzi"
	"github.com/dgf/protzi/components"
)

func ExampleNetwork_echo() {
	network := protzi.New("echo net")
	network.Add("proxy", &components.Echo{})
	network.In("proxy.Ping")
	network.Out("proxy.Pong")
	network.Run()

	fmt.Println(network.Process("echo"))
	// Output: echo
}

func ExampleNetwork_fileWordCounter() {
	data := "\f\n\r\ttwo\fone\ntwo\r\t"

	// create temp file
	file, err := ioutil.TempFile(os.TempDir(), "example")
	if err != nil {
		panic(err)
	}

	// write test data string
	if _, err := file.Write([]byte(data)); err != nil {
		panic(err)
	}

	// close temp file
	if err := file.Close(); err != nil {
		panic(err)
	}

	// create network
	network := protzi.New("file word counter")

	// add process components
	network.Add("read", &components.TextFileRead{})
	network.Add("count", &components.WordCount{})

	// connect components
	network.In("read.File")
	network.Connect("read.Text", "count.Text")
	network.Out("count.Counts")

	// run it
	network.Run()

	// process file
	result := network.Process(file.Name())
	countsByWord, ok := result.(map[string]int)
	if !ok {
		panic(fmt.Sprintf("invalid result: %s", result))
	}

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
