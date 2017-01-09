package cmd

import (
	"github.com/ElyKar/mpm/core"
	"github.com/spf13/cobra"
)

// Changes the passphrase for the storage
var changeCmd = &cobra.Command{
	Use:   "change",
	Short: "Change the master password",
	Run:   chainNodes(storageExists, verifyPassphrase, createPassphrase, changeFunc, updateStore),
}

// changeFunc requires the storage, current passphrase and new passphrase to be stored in the context
func changeFunc(context map[string]interface{}) (string, int) {
	old := (context["passphrase"]).(string)
	new := string((context["newPassphrase"]).([]byte))
	var storage *core.Storage = (context["storage"]).(*core.Storage)

	storage.SetNewPassphrase(old, new)
	return "", 0
}

func init() {
	RootCmd.AddCommand(changeCmd)
}
