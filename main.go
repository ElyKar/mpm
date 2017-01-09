package main

import (
	"fmt"
	"os"

	"github.com/ElyKar/mpm/cmd"
)

func main() {
	// Just delegate everything to the custom commands
	if err := cmd.RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
}
