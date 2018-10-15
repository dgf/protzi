package component_test

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/dgf/protzi/component"
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
	go (&component.TextFileRead{File: files, Text: text}).Run()

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

func ExampleTextFileRead_Run_fileNotFound() {
	files := make(chan string)
	failures := make(chan string)
	go (&component.TextFileRead{File: files, Error: failures}).Run()

	unknown := "/unknown/test/file/name"
	files <- unknown
	fmt.Println(<-failures)

	// Output: Error: file /unknown/test/file/name not found.
}
