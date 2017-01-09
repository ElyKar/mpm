package cmd

import (
	"fmt"
	"os"

	"github.com/ElyKar/mpm/core"
	"github.com/howeyc/gopass"
	"github.com/spf13/cobra"
)

// nodeFunc is the type of the functions that will be executed. As there are lots of operations which will be the same for the different commands (verifying passphrase, the master file, etc...), nodeFuncs are designed to be as small and as reusable as possible. Each function has the possibility to throw an error aswell as an error code through its return type. For example, if the passphrase is wrong, or the storage file does not exist, then the CLI should go into error. To communicate, node have access to a shared context.
type nodeFunc func(map[string]interface{}) (string, int)

// Here are stored the two flags used by the CLI
// Name of the section
var section string

// Name of the password inside the section
var name string

// Node asserting the section and name exist and are not empty
func sectionAndNameRequired(context map[string]interface{}) (string, int) {
	if section == "" || name == "" {
		return "You need to provide a section and a name for your entry !", 1
	}

	return "", 0
}

// Node asserting the storage exists and is valid. On success, it stores the storage in the context.
func storageExists(context map[string]interface{}) (string, int) {

	if storage, err := core.GetStorage(); err != nil {
		return fmt.Sprintf(`No previous storage found, error is:
%s

Are you sure the file is created and you have proper access rights ?
You can initialize your mpm file with the command 'mpm init'
	`, err.Error()), 1
	} else {
		context["storage"] = storage
		return "", 0
	}
}

// Prompts the user for the passphrase and verifies it. On success, it stores the passphrase in the context (will be used for encoding/decoding)
func verifyPassphrase(context map[string]interface{}) (string, int) {
	var storage *core.Storage = (context["storage"]).(*core.Storage)
	success := false
	var passphrase string
	for i := 0; i < 3 && !success; i++ {

		fmt.Printf("Enter your passphrase: ")
		passB, _ := gopass.GetPasswd()
		passphrase = string(passB)

		if err := storage.CheckPassphrase(passphrase); err != nil {
			fmt.Printf("Wrong passphrase\n\n")
		} else {
			success = true
		}
	}

	if !success {
		return "Try again later !", 1
	} else {
		context["passphrase"] = passphrase
		return "", 0
	}
}

// Updates the storage on the disk. It retrieves the storage from the context, and tries to dump it on the disk.
func updateStore(context map[string]interface{}) (string, int) {
	var storage *core.Storage = (context["storage"]).(*core.Storage)

	if err := storage.DumpOnDisk(); err != nil {
		return fmt.Sprintf("Something went wrong, your changes haven't been saved. Try again later !\n%s", err), 1
	}

	return "\nEverything went well !", 0
}

// Prompts for a new passphrase. On success, it is stored on the context under 'newPassword'.
func createPassphrase(context map[string]interface{}) (string, int) {
	fmt.Printf("Enter your new passphrase: ")
	pass1, err := gopass.GetPasswd()
	if err != nil {
		return fmt.Sprintf("An error occurred !\n%s", err), 1
	}

	fmt.Printf("Re-enter your passphrase: ")
	pass2, err := gopass.GetPasswd()
	if err != nil {
		return fmt.Sprintf("An error occurred !\n%s", err), 1
	}

	if string(pass1) != string(pass2) {
		return "Passphrases do not match !", 1
	}

	context["newPassphrase"] = pass1
	return "", 0
}

// Simple function to chain nodes and create the actual Run function for *cobra.Command
func chainNodes(nodes ...nodeFunc) func(*cobra.Command, []string) {

	return func(cmd *cobra.Command, args []string) {
		context := make(map[string]interface{})
		var msg string
		var code int

		for _, fptr := range nodes {
			msg, code = fptr(context)
			if msg != "" {
				fmt.Println(msg)
				os.Exit(code)
			}
		}
	}
}
