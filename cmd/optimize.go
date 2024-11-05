package cmd

import (
	"bytes" // Added this import
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

var (
	dockerfilePath   string
	dockerignorePath string
	packageJSONPath  string
	openaiAPIKey     string
)

type Config struct {
	APIKey string `json:"api_key"`
}

type APIRequest struct {
	Dockerfile   string `json:"Dockerfile,omitempty"`
	Dockerignore string `json:".dockerignore,omitempty"`
	PackageJSON  string `json:"package.json,omitempty"`
	OpenAIAPIKey string `json:"openai_api_key,omitempty"`
}

type Action struct {
	Filename    string `json:"filename"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Rule        string `json:"rule"`
}

type APIResponse struct {
	ModifiedProject map[string]string `json:"modified_project"`
	ActionsTaken    []Action          `json:"actions_taken"`
	Recommendations []Action          `json:"recommendations"`
	Error           string            `json:"error"`
}

var optimizeCmd = &cobra.Command{
	Use:   "optimize",
	Short: "Modifies code files to reduce the size of your Nodejs project's docker Image",
	Long: `Optimize your NodeJS Docker project by sending your Dockerfile, .dockerignore, and package.json
to the Dockershrink backend. You can provide your OpenAI API key using the --openai-api-key flag or by setting
the OPENAI_API_KEY environment variable to enable AI features (highly recommended).`,
	Run: func(cmd *cobra.Command, args []string) {
		// Load API key from config
		configPath := filepath.Join(os.Getenv("HOME"), ".dsconfig.json")
		configFile, err := os.Open(configPath)
		if err != nil {
			fmt.Println("Error: API key not configured. Please run 'dockershrink init --api-key <api key>' first.")
			os.Exit(1)
		}
		defer configFile.Close()

		var config Config
		decoder := json.NewDecoder(configFile)
		if err := decoder.Decode(&config); err != nil {
			fmt.Printf("Error reading config file: %v\n", err) // Updated to include error
			os.Exit(1)
		}

		// Collect files
		dockerfileContent := readFile(dockerfilePath, "Dockerfile")
		dockerignoreContent := readFile(dockerignorePath, ".dockerignore")
		packageJSONContent := readFile(packageJSONPath, "package.json")

		// Prepare API request
		apiReq := APIRequest{
			Dockerfile:   dockerfileContent,
			Dockerignore: dockerignoreContent,
			PackageJSON:  packageJSONContent,
		}

		if openaiAPIKey == "" {
			openaiAPIKey = os.Getenv("OPENAI_API_KEY")
		}
		if openaiAPIKey != "" {
			apiReq.OpenAIAPIKey = openaiAPIKey
		}

		// Send API request
		serverURL := os.Getenv("SERVER_URL")
		if serverURL == "" {
			serverURL = "https://dockershrink.com"
		}
		apiURL := serverURL + "/api/v1/optimize"

		reqBody, err := json.Marshal(apiReq)
		if err != nil {
			fmt.Printf("Error preparing API request: %v\n", err) // Updated
			os.Exit(1)
		}

		client := &http.Client{}
		request, err := http.NewRequest("POST", apiURL, bytes.NewReader(reqBody)) // Updated
		if err != nil {
			fmt.Printf("Error creating API request: %v\n", err) // Updated
			os.Exit(1)
		}
		request.Header.Set("Content-Type", "application/json")
		request.Header.Set("Authorization", config.APIKey)

		response, err := client.Do(request)
		if err != nil {
			fmt.Printf("Error sending API request: %v\n", err) // Updated
			os.Exit(1)
		}
		defer response.Body.Close()

		// Handle API response
		body, err := ioutil.ReadAll(response.Body)
		if err != nil {
			fmt.Printf("Error reading API response: %v\n", err) // Updated
			os.Exit(1)
		}

		var apiResp APIResponse
		if err := json.Unmarshal(body, &apiResp); err != nil {
			fmt.Printf("Error parsing API response: %v\n", err) // Updated
			os.Exit(1)
		}

		if response.StatusCode != http.StatusOK {
			fmt.Printf("Error: %s\n", apiResp.Error)
			os.Exit(1)
		}

		// Write modified files
		outputDir := filepath.Join(".", "dockershrink.optimised")
		if err := os.MkdirAll(outputDir, os.ModePerm); err != nil { // Added error handling
			fmt.Printf("Error creating output directory: %v\n", err)
			os.Exit(1)
		}

		for filename, content := range apiResp.ModifiedProject {
			outputPath := filepath.Join(outputDir, filename)
			err := ioutil.WriteFile(outputPath, []byte(content), 0644)
			if err != nil {
				fmt.Printf("Error writing file %s: %v\n", filename, err) // Updated
				os.Exit(1)
			}
		}

		// Display actions and recommendations
		printActions(apiResp.ActionsTaken, "Actions Taken")
		printActions(apiResp.Recommendations, "Recommendations")
	},
}

func init() {
	optimizeCmd.Flags().StringVar(&dockerfilePath, "dockerfile", "", "Path to Dockerfile")
	optimizeCmd.Flags().StringVar(&dockerignorePath, "dockerignore", "", "Path to .dockerignore")
	optimizeCmd.Flags().StringVar(&packageJSONPath, "package-json", "", "Path to package.json")
	optimizeCmd.Flags().StringVar(&openaiAPIKey, "openai-api-key", "", "Your OpenAI API key (can also be set from the OPENAI_API_KEY environment variable)")
}

func readFile(path string, defaultName string) string {
	if path == "" {
		path = filepath.Join(".", defaultName)
		if _, err := os.Stat(path); os.IsNotExist(err) && defaultName == "package.json" {
			path = filepath.Join(".", "src", defaultName)
		}
	}
	content, err := ioutil.ReadFile(path)
	if err != nil {
		return ""
	}
	return string(content)
}

func printActions(actions []Action, title string) {
	if len(actions) == 0 {
		return
	}
	fmt.Println()
	fmt.Println("========== " + title + " ==========")
	for _, action := range actions {
		filenameColor := color.New(color.FgCyan).SprintFunc()
		titleColor := color.New(color.FgGreen).SprintFunc()
		descColor := color.New(color.FgWhite).SprintFunc()
		ruleColor := color.New(color.FgYellow).SprintFunc()

		fmt.Printf("File: %s\n", filenameColor(action.Filename))
		fmt.Printf("Title: %s\n", titleColor(action.Title))
		fmt.Printf("Description: %s\n", descColor(action.Description))
		fmt.Printf("Rule: %s\n", ruleColor(action.Rule))
		fmt.Println("-----------------------------------")
	}
}
