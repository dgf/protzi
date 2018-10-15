package component

import (
	"bytes"
	"text/template"
)

// TextTemplate component renders a template.
type TextTemplate struct {
	Template <-chan *template.Template
	Data     <-chan interface{}
	Output   chan<- string
	Error    chan<- string
}

// Run renders the template.
func (t *TextTemplate) Run() {
	tmpl := <-t.Template
	for data := range t.Data {
		b := &bytes.Buffer{}
		if err := tmpl.Execute(b, data); err != nil {
			t.Error <- "Error: " + err.Error()
		} else {
			t.Output <- b.String()
		}
	}
}
