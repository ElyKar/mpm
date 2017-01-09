package cmd

// The packages contains all commands and helpers for CLI interaction

import "github.com/spf13/cobra"

// The root command of this interface
var RootCmd = &cobra.Command{
	Use:   "mpm <action>",
	Short: "mpm is a sweet and tiny password manager written in Go.",
	Long: `mpm is a CLI password manager made to handle all of your passwords.
You can customise each password by choosing which characters it may contain and its total length.`,
}
