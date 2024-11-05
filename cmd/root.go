package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "dockershrink",
	Short: "Dockershrink optimizes your NodeJS Docker images.",
	Long: `Dockershrink is a CLI tool that helps you reduce the size of your NodeJS Docker images
by applying best practices and optimizations to your Dockerfile and related files.`,
	// Updated Run function to show help message
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

// Execute executes the root command.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.AddCommand(initCmd)
	rootCmd.AddCommand(optimizeCmd)
}
