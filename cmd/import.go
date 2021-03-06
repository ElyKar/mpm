package cmd

import (
	"github.com/ElyKar/mpm/core"
	"github.com/spf13/cobra"
)

// Imports an existing password
var importCmd = &cobra.Command{
	Use:   "import --section <section> --name <name>",
	Short: "Imports an existing password in the storage",
	Run:   chainNodes(sectionAndNameRequired, storageExists, verifyPassphrase, verifyErase, createPassword, importFunc, updateStore),
}

// importFunc requires the storage and passphrase from the context, and also non-empty name and section
func importFunc(context map[string]interface{}) (string, int) {
	var storage *core.Storage = (context["storage"]).(*core.Storage)

	// Assert we won't erase a password
	password := string(context["newPass"].([]byte))
	encoder := core.NewTranscoder(context["passphrase"].(string))
	encoded, _ := encoder.EncodePassword(password)

	storage.Set(section, name, string(encoded))
	return "", 0
}

func init() {
	importCmd.Flags().StringVar(&section, "section", "", "The section to add the imported password")
	importCmd.Flags().StringVar(&name, "name", "", "A name for your the imported password")

	RootCmd.AddCommand(importCmd)

}
