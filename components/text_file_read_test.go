package components_test

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/dgf/protzi/components"
)

func ExampleTextFileRead() {
	data := "one\ntwo"

	// create temp file
	file, err := ioutil.TempFile(os.TempDir(), "reader_test")
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

	// create and run reader
	files := make(chan string)
	text := make(chan string)
	go (&components.TextFileRead{File: files, Text: text}).Run()

	// send file name
	files <- file.Name()

	// print file content
	fmt.Println(<-text)

	// delete temporary file
	if err := os.Remove(file.Name()); err != nil {
		panic(err)
	}

	// Output:
	// one
	// two
}
