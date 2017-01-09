package cmd

import (
	"fmt"
	"os"
	"text/template"

	"github.com/ElyKar/mpm/core"
	"github.com/spf13/cobra"
)

// Template for listing sections and their entries.
const listAllTmpl = `{{ range $k, $v := .All }}    - {{ $k }}
{{ range $idx, $name := $v }}        * {{ $name }}
{{ end }}{{ end }}
`

// Template for listing sections alone.
const listSectionsTmpl = `{{ range $idx, $v := .Sections }}    - {{ $v }}
{{ end }}
`

// Root List command
var listCmd = &cobra.Command{
	Use:   "list [all|sections|passwords]",
	Short: "List the sections and passwords stored",
}

// Lists all the sections and their content
var listAllCmd = &cobra.Command{
	Use:   "all",
	Short: "List all sections and passwords",
	Run:   chainNodes(storageExists, listAllFunc),
}

// Node for listing all the content of the storage. It requires the storage from the context.
func listAllFunc(context map[string]interface{}) (string, int) {
	storage := (context["storage"]).(*core.Storage)

	entries := storage.ListAll()
	tmpl := template.Must(template.New("allTmpl").Parse(listAllTmpl))

	fmt.Println("Here are the passwords stored:")
	tmpl.Execute(os.Stdout, struct{ All map[string][]string }{entries})
	return "", 0
}

// Lists only the section
var listSectionsCmd = &cobra.Command{
	Use:   "sections",
	Short: "List all sections stored",
	Run:   chainNodes(storageExists, listSectionsFunc),
}

// Node for listing all the sections of the storage. It requires the storage from the context.
func listSectionsFunc(context map[string]interface{}) (string, int) {
	storage := (context["storage"]).(*core.Storage)

	entries := storage.ListSections()
	tmpl := template.Must(template.New("sectionsTmpl").Parse(listSectionsTmpl))

	fmt.Println("Here are the sections stored:")
	tmpl.Execute(os.Stdout, struct{ Sections []string }{entries})
	return "", 0
}

// Lists passwords for a section. If the section does not exist, no password will appear
var listPasswordCmd = &cobra.Command{
	Use:   "passwords --section <section>",
	Short: "List all passwords of a section",
	Run:   chainNodes(sectionRequired, storageExists, listPasswordFunc),
}

// Internal node, only the section flag is required to be non-empty
func sectionRequired(context map[string]interface{}) (string, int) {
	if section == "" {
		return "You need to provide a section !", 1
	}

	return "", 0
}

// Node for listing all the passwords of a given section. It requires the storage from the context.
func listPasswordFunc(context map[string]interface{}) (string, int) {
	storage := (context["storage"]).(*core.Storage)

	entries := storage.ListPasswords(section)
	tmpl := template.Must(template.New("allTmpl").Parse(listAllTmpl))

	fmt.Println("Here are the passwords stored:")
	tmpl.Execute(os.Stdout, struct{ All map[string][]string }{map[string][]string{section: entries}})
	return "", 0
}

func init() {
	listPasswordCmd.Flags().StringVar(&section, "section", "", "The section to list passwords for")
	listCmd.AddCommand(listPasswordCmd)
	listCmd.AddCommand(listAllCmd)
	listCmd.AddCommand(listSectionsCmd)

	RootCmd.AddCommand(listCmd)
}
