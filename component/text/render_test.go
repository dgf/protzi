package text_test

import (
	"fmt"

	"github.com/dgf/protzi/component/text"
)

func ExampleRender_Run() {
	templates := make(chan string)
	input := make(chan interface{})
	output := make(chan string)
	failures := make(chan string)

	renderer := &text.Render{
		Template: templates,
		Data:     input,
		Output:   output,
		Error:    failures,
	}
	go renderer.Run()

	t := `{{if eq "de" .Lang}}Hallo{{else if eq "es" .Lang}}Hola{{else}}Hello{{end}} {{.Name}}!`
	for _, data := range []struct{ Lang, Name string }{
		{"es", "mundo"},
		{"en", "world"},
		{"de", "Welt"},
	} {
		templates <- t
		input <- data
		select {
		case o := <-output:
			fmt.Println(o)
		case f := <-failures:
			fmt.Println(f)
		}
	}
	// Output:
	// Hola mundo!
	// Hello world!
	// Hallo Welt!
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
	// Output: template: 13cbf0569d57e335d0cea50cf32c32c013720e17:1:2: executing "13cbf0569d57e335d0cea50cf32c32c013720e17" at <.Test>: can't evaluate field Test in type struct {}
}

func ExampleRender_Run_invalidTemplate() {
	templates := make(chan string)
	failures := make(chan string)

	renderer := &text.Render{
		Template: templates,
		Error:    failures,
	}
	go renderer.Run()

	templates <- "{{Invalid?}}"
	fmt.Println(<-failures)
	// Output: template: 464efc87b5ca3af0e58ac2c469c4df507ef16ddc:1: unexpected bad character U+003F '?' in command
}
