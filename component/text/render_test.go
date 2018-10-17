package text_test

import (
	"fmt"

	"github.com/dgf/protzi/component/text"
)

func ExampleRender_Run() {
	templates := make(chan string)
	data := make(chan interface{})
	output := make(chan string)

	renderer := &text.Render{
		Template: templates,
		Data:     data,
		Output:   output,
	}
	go renderer.Run()

	templates <- "Hello {{.Name}}!"
	data <- struct{ Name string }{Name: "World"}
	fmt.Println(<-output)
	// Output: Hello World!
}

func ExampleRender_Run_invalidData() {
	templates := make(chan string)
	data := make(chan interface{})
	failures := make(chan string)

	renderer := &text.Render{
		Template: templates,
		Data:     data,
		Error:    failures,
	}
	go renderer.Run()

	templates <- "{{.Test}}"
	data <- struct{}{}
	fmt.Println(<-failures)
	// Output:
	// Execute error: template: text:1:2: executing "text" at <.Test>: can't evaluate field Test in type struct {}
}

func ExampleTextTemplate_Run_invalidTemplate() {
	templates := make(chan string)
	failures := make(chan string)

	renderer := &text.Render{
		Template: templates,
		Error:    failures,
	}
	go renderer.Run()

	templates <- "{{Invalid?}}"
	fmt.Println(<-failures)
	// Output:
	// Render error: template: text:1: unexpected bad character U+003F '?' in command
}
