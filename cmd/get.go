package cmd

import (
	"fmt"

	"github.com/ElyKar/mpm/core"
	"github.com/atotto/clipboard"
	"github.com/spf13/cobra"
)

// get the password for the section and name, then copy it to the clipboard
var getCmd = &cobra.Command{
	Use:   "get",
	Short: "Copy a password to your clipboard",
	Run:   chainNodes(sectionAndNameRequired, storageExists, verifyPassphrase, getFunc),
}

// Get the password from the storage, decodes it, then copies it into the clipboard. It needs the storage and passphrase from the context.
func getFunc(context map[string]interface{}) (string, int) {
	var storage *core.Storage = (context["storage"]).(*core.Storage)
	var passphrase string = (context["passphrase"]).(string)

	decoder := core.NewTranscoder(passphrase)
	encoded, err := storage.Get(section, name)

	// Either section or name does not exist
	if err != nil {
		return err.Error(), 0
	}

	decoded, err := decoder.DecodePassword(encoded)
	if err != nil {
		return fmt.Sprintf("An error occurred:\n%s", err.Error()), 1
	}

	if err = clipboard.WriteAll(string(decoded)); err != nil {
		return fmt.Sprintf("Impossible to copy to clipboard.\n%s", err), 1
	} else {
		return "Your password has been successfully copied to your clipboard", 0
	}

}

func init() {
	getCmd.Flags().StringVar(&section, "section", "", "The section to get your password from")
	getCmd.Flags().StringVar(&name, "name", "", "The password you want")

	RootCmd.AddCommand(getCmd)
}
