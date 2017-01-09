package cmd

import (
	"github.com/ElyKar/mpm/core"
	"github.com/spf13/cobra"
)

// Rudimentary command, it simply initializes an empty storage with a passphrase.
var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize an empty store for mpm",
	Run:   chainNodes(createPassphrase, initFunc, updateStore),
}

// Requires the newPassphrase from the context, and creates a new storage with it.
func initFunc(context map[string]interface{}) (string, int) {
	var passphrase []byte = (context["newPassphrase"]).([]byte)

	if _, err := core.GetStorage(); err == nil {
		return "There is already a store !", 1
	}

	context["storage"] = core.InitPassphrase(passphrase)
	return "", 0
}

func init() {
	RootCmd.AddCommand(initCmd)
}
