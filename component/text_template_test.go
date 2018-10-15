package component_test

import (
	"fmt"
	"text/template"

	"github.com/dgf/protzi/component"
)

func ExampleTextTemplate_Run() {
	templates := make(chan *template.Template)
	data := make(chan interface{})
	output := make(chan string)

	go (&component.TextTemplate{
		Template: templates,
		Data:     data,
		Output:   output,
	}).Run()

	helloTemplate, err := template.New("hello").Parse("Hello {{.Name}}!")
	if err != nil {
		panic(err)
	}

	templates <- helloTemplate
	data <- struct{ Name string }{Name: "World"}
	fmt.Println(<-output)
	// Output: Hello World!
}

func ExampleTextTemplate_Run_invalidData() {
	templates := make(chan *template.Template)
	data := make(chan interface{})
	failures := make(chan string)

	go (&component.TextTemplate{
		Template: templates,
		Data:     data,
		Error:    failures,
	}).Run()

	invalidTemplate, err := template.New("invalid").Parse("{{.Test}}")
	if err != nil {
		panic(err)
	}

	templates <- invalidTemplate
	data <- struct{}{}
	fmt.Println(<-failures)
	// Output:
	// Error: template: invalid:1:2: executing "invalid" at <.Test>: can't evaluate field Test in type struct {}
}
