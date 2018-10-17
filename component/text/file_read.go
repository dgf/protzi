package text

import "io/ioutil"

// FileRead component reads text files.
type FileRead struct {
	File  <-chan string
	Text  chan<- string
	Error chan<- string
}

// Run tries to read the file and outputs the text or an error.
func (r *FileRead) Run() {
	for f := range r.File {
		if t, err := ioutil.ReadFile(f); err != nil {
			r.Error <- "Error: file " + f + " not found."
		} else {
			r.Text <- string(t)
		}
	}
}
