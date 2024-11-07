package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print current CLI version",
	Long:  `Prints the current version of the Dockershrink CLI.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("Dockershrink CLI version %s\n", version)
	},
}
