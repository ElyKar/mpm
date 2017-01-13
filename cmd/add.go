package cmd

import (
	"fmt"

	"github.com/ElyKar/mpm/core"
	"github.com/spf13/cobra"
)

// Adds a new password to the master file
var addCmd = &cobra.Command{
	Use:   "add --section <section> --name <name>",
	Short: "Generates a new password for the section and name",
	Long:  `Interacts with the user to generate a new password`,
	Run:   chainNodes(sectionAndNameRequired, storageExists, verifyPassphrase, verifyErase, addFunc, updateStore),
}

// addFunc requires the storage and passphrase from the context, and also non-empty name and section
func addFunc(context map[string]interface{}) (string, int) {
	var storage *core.Storage = (context["storage"]).(*core.Storage)

	// Prompt user for alphabet to choose
	var choice int
	if err := chooseAlphabet(&choice); err != nil {
		return err.Error(), 1
	}

	// Prompt user for password length
	length, min, max := 0, 8, 1000
	interactI("Length of your password:  ", &length)
	if length < min || length > max {
		return fmt.Sprintf("Length must be comprised between %d and %d, received %d\n", min, max, length), 1
	}

	// Generates, encrypts, encode and save changes
	password := core.Alphas[choice].GenPassword(length)
	encoder := core.NewTranscoder(context["passphrase"].(string))
	encoded, _ := encoder.EncodePassword(password)

	storage.Set(section, name, string(encoded))
	return "", 0
}

func init() {
	addCmd.Flags().StringVar(&section, "section", "", "The section to add the newly-generated password")
	addCmd.Flags().StringVar(&name, "name", "", "A name for your the newly-generated password")

	RootCmd.AddCommand(addCmd)

}
