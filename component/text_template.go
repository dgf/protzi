package component

import (
	"bytes"
	"text/template"
)

// TextTemplate component renders a template.
type TextTemplate struct {
	Template <-chan string
	Data     <-chan interface{}
	Output   chan<- string
	Error    chan<- string
}

// Run renders the template.
func (t *TextTemplate) Run() {
	for templateString := range t.Template {
		if parsedTemplate, err := template.New("text").Parse(templateString); err != nil {
			t.Error <- "Template error: " + err.Error()
		} else {
			writer := &bytes.Buffer{}
			if err := parsedTemplate.Execute(writer, <-t.Data); err != nil {
				t.Error <- "Execute error: " + err.Error()
			} else {
				t.Output <- writer.String()
			}
		}
	}
}
