package text

import (
	"bytes"
	"crypto/sha1"
	"fmt"
	"text/template"
)

// Render component renders a template.
type Render struct {
	Template <-chan string
	Data     <-chan interface{}
	Output   chan<- string
	Error    chan<- string
}

var templates = map[string]*template.Template{}

func hash(s string) string {
	h := sha1.New()
	h.Write([]byte(s))
	bs := h.Sum(nil)
	return fmt.Sprintf("%x", bs)
}

// Run parses and executes a template.
func (r *Render) Run() {
	for tmpl := range r.Template {
		hash := hash(tmpl)
		instance, ok := templates[hash]
		if !ok {
			if parsed, err := template.New(hash).Parse(tmpl); err != nil {
				r.Error <- err.Error()
				continue
			} else {
				templates[hash] = parsed
				instance = parsed
			}
		}

		data := <-r.Data
		writer := &bytes.Buffer{}
		if err := instance.Execute(writer, data); err != nil {
			r.Error <- err.Error()
		} else {
			r.Output <- writer.String()
		}
	}
}
