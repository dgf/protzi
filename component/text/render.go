package text

import (
	"bytes"
	"text/template"
)

// Render component renders a template.
type Render struct {
	Template <-chan string
	Data     <-chan interface{}
	Output   chan<- string
	Error    chan<- string
}

// Run renders the template.
func (r *Render) Run() {
	for t := range r.Template {
		if parsedTemplate, err := template.New("text").Parse(t); err != nil {
			r.Error <- "Render error: " + err.Error()
		} else {
			writer := &bytes.Buffer{}
			if err := parsedTemplate.Execute(writer, <-r.Data); err != nil {
				r.Error <- "Execute error: " + err.Error()
			} else {
				r.Output <- writer.String()
			}
		}
	}
}
