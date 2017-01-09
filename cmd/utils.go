package cmd

import (
	"fmt"
	"os"
	"text/template"

	"github.com/ElyKar/mpm/core"
)

// Template for choosing alphabet
const alphaTmpl = `Choose your alphabet:
{{ range $idx, $elt := .Possibilities }}    [{{$idx}}]  {{$elt.Display}}
{{ end }}`

// Prompts the user for a string
func interactS(question string, answer *string) {
	fmt.Printf(question)
	fmt.Scanf("%s\n", answer)
}

// Prompts the user for an integer
func interactI(question string, answer *int) {
	fmt.Printf(question)
	fmt.Scanf("%d\n", answer)
}

// Display alphabets and let the user choose. It checks that the alphabet exists.
func chooseAlphabet(choice *int) error {
	tmpl := template.Must(template.New("alphabets").Parse(alphaTmpl))
	err := tmpl.Execute(os.Stdout, struct{ Possibilities []core.Alphabet }{core.Alphas})
	if err != nil {
		panic(err)
	}

	interactI("\nWhat's your choice ?   ", choice)

	if *choice < 0 || *choice >= len(core.Alphas) {
		return fmt.Errorf("Invalid choice: %d", *choice)
	}

	return nil
}
