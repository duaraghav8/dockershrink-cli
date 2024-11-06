package cmd

import (
	"bytes"
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
	Dockerfile   string          `json:"Dockerfile,omitempty"`
	Dockerignore string          `json:".dockerignore,omitempty"`
	PackageJSON  json.RawMessage `json:"package.json,omitempty"`
	OpenAIAPIKey string          `json:"openai_api_key,omitempty"`
}

type Action struct {
	Filename    string `json:"filename"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Rule        string `json:"rule"`
}

type APIResponse struct {
	ModifiedProject map[string]json.RawMessage `json:"modified_project"`
	ActionsTaken    []Action                   `json:"actions_taken"`
	Recommendations []Action                   `json:"recommendations"`
	Error           string                     `json:"error"`
}

var optimizeCmd = &cobra.Command{
	Use:   "optimize",
	Short: "Optimize your NodeJS Docker project",
	Long: `Optimize your NodeJS Docker project by sending your Dockerfile, .dockerignore, and package.json
to the Dockershrink API. You can provide your OpenAI API key using the --openai-api-key flag or by setting
the OPENAI_API_KEY environment variable.`,
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
			fmt.Printf("Error reading config file: %v\n", err)
			os.Exit(1)
		}

		// Collect files and their statuses
		dockerfileContent, dockerfileUsedPath, dockerfileFound := readFile(dockerfilePath, "Dockerfile")
		dockerignoreContent, dockerignoreUsedPath, dockerignoreFound := readFile(dockerignorePath, ".dockerignore")
		packageJSONContent, packageJSONUsedPath, packageJSONFound := readJSONFile(packageJSONPath, "package.json")

		// Inform the user about which files were picked up
		fmt.Println()
		if dockerfileFound {
			fmt.Printf("- Using %s\n", dockerfileUsedPath)
		} else {
			fmt.Printf("- No Dockerfile found in the default paths\n")
		}

		if dockerignoreFound {
			fmt.Printf("- Using %s\n", dockerignoreUsedPath)
		} else {
			fmt.Printf("- No .dockerignore found in the default paths\n")
		}

		if packageJSONFound {
			fmt.Printf("- Using %s\n", packageJSONUsedPath)
		} else {
			fmt.Printf("- No package.json found in the default paths\n")
		}
		fmt.Println()

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
			fmt.Printf("Error preparing API request: %v\n", err)
			os.Exit(1)
		}

		client := &http.Client{}
		request, err := http.NewRequest("POST", apiURL, bytes.NewReader(reqBody))
		if err != nil {
			fmt.Printf("Error creating API request: %v\n", err)
			os.Exit(1)
		}
		request.Header.Set("Content-Type", "application/json")
		request.Header.Set("Authorization", config.APIKey)

		response, err := client.Do(request)
		if err != nil {
			fmt.Printf("Error sending API request: %v\n", err)
			os.Exit(1)
		}
		defer response.Body.Close()

		// Handle API response
		body, err := ioutil.ReadAll(response.Body)
		if err != nil {
			fmt.Printf("Error reading API response: %v\n", err)
			os.Exit(1)
		}

		var apiResp APIResponse
		if err := json.Unmarshal(body, &apiResp); err != nil {
			fmt.Printf("Error parsing API response: %v\n", err)
			os.Exit(1)
		}

		if response.StatusCode != http.StatusOK {
			fmt.Printf("Error: %s\n", apiResp.Error)
			os.Exit(1)
		}

		// Write modified files
		outputDir := filepath.Join(".", "dockershrink.optimised")
		if err := os.MkdirAll(outputDir, os.ModePerm); err != nil {
			fmt.Printf("Error creating output directory: %v\n", err)
			os.Exit(1)
		}

		for filename, content := range apiResp.ModifiedProject {
			outputPath := filepath.Join(outputDir, filename)
			if filename == "package.json" {
				// Handle package.json as JSON object
				var formattedJSON bytes.Buffer
				err := json.Indent(&formattedJSON, content, "", "  ")
				if err != nil {
					fmt.Printf("Error formatting JSON for %s: %v\n", filename, err)
					os.Exit(1)
				}
				err = ioutil.WriteFile(outputPath, formattedJSON.Bytes(), 0644)
				if err != nil {
					fmt.Printf("Error writing file %s: %v\n", filename, err)
					os.Exit(1)
				}
			} else {
				// Handle other files as strings
				var fileContent string
				err := json.Unmarshal(content, &fileContent)
				if err != nil {
					fmt.Printf("Error parsing content for %s: %v\n", filename, err)
					os.Exit(1)
				}
				err = ioutil.WriteFile(outputPath, []byte(fileContent), 0644)
				if err != nil {
					fmt.Printf("Error writing file %s: %v\n", filename, err)
					os.Exit(1)
				}
			}
		}

		// Check if there are no actions taken and no recommendations
		if len(apiResp.ActionsTaken) == 0 && len(apiResp.Recommendations) == 0 {
			fmt.Println("Docker image is already optimized, no further actions were taken by dockershrink.")
			return
		}

		// Display actions and recommendations
		printActions(apiResp.ActionsTaken, "Actions Taken")
		printActions(apiResp.Recommendations, "Recommendations")
	},
}

func init() {
	optimizeCmd.Flags().StringVar(&dockerfilePath, "dockerfile", "", "Path to Dockerfile (default: ./Dockerfile)")
	optimizeCmd.Flags().StringVar(&dockerignorePath, "dockerignore", "", "Path to .dockerignore (default: ./.dockerignore)")
	optimizeCmd.Flags().StringVar(&packageJSONPath, "package-json", "", "Path to package.json (default: ./package.json or ./src/package.json)")
	optimizeCmd.Flags().StringVar(&openaiAPIKey, "openai-api-key", "", "Your OpenAI API key (alternatively, set the OPENAI_API_KEY environment variable)")
}

func readFile(path string, defaultName string) (string, string, bool) {
	var usedPath string
	var found bool
	if path != "" {
		usedPath = path
	} else {
		usedPath = filepath.Join(".", defaultName)
	}

	content, err := ioutil.ReadFile(usedPath)
	if err != nil {
		found = false
		return "", usedPath, found
	}
	found = true
	return string(content), usedPath, found
}

func readJSONFile(path string, defaultName string) (json.RawMessage, string, bool) {
	var usedPath string
	var found bool
	if path != "" {
		usedPath = path
	} else {
		usedPath = filepath.Join(".", defaultName)
		if _, err := os.Stat(usedPath); os.IsNotExist(err) && defaultName == "package.json" {
			usedPath = filepath.Join(".", "src", defaultName)
		}
	}

	content, err := ioutil.ReadFile(usedPath)
	if err != nil {
		found = false
		return nil, usedPath, found
	}

	// Validate that it's valid JSON
	var temp interface{}
	if err := json.Unmarshal(content, &temp); err != nil {
		fmt.Printf("Error parsing JSON file %s: %v\n", usedPath, err)
		os.Exit(1)
	}
	found = true
	return json.RawMessage(content), usedPath, found
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

		fmt.Printf("File: %s\n", filenameColor(action.Filename))
		fmt.Printf("Title: %s\n", titleColor(action.Title))
		fmt.Printf("Description: %s\n", descColor(action.Description))
		fmt.Println("-----------------------------------")
	}
}
