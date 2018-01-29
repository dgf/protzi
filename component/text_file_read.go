package component

import "io/ioutil"

// TextFileRead component reads text files.
type TextFileRead struct {
	File <-chan string
	Text chan<- string
}

func (r *TextFileRead) Run() {
	for f := range r.File {
		t, _ := ioutil.ReadFile(f)
		r.Text <- string(t)
	}
}
