package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

var apiKey string

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize dockershrink with your API key",
	Run: func(cmd *cobra.Command, args []string) {
		if apiKey == "" {
			fmt.Println("Error: --api-key is required")
			os.Exit(1)
		}

		config := map[string]string{
			"api_key": apiKey,
		}

		configPath := filepath.Join(os.Getenv("HOME"), ".dsconfig.json")
		file, err := os.Create(configPath)
		if err != nil {
			fmt.Printf("Error creating config file: %v\n", err)
			os.Exit(1)
		}
		defer file.Close()

		encoder := json.NewEncoder(file)
		if err := encoder.Encode(config); err != nil {
			fmt.Printf("Error writing to config file: %v\n", err)
			os.Exit(1)
		}

		fmt.Println("API key saved successfully.")
	},
}

func init() {
	initCmd.Flags().StringVar(&apiKey, "api-key", "", "Your Dockershrink API key")
	initCmd.MarkFlagRequired("api-key")
}
