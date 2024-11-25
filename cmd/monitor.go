package cmd

import (
	"bytes"
	"fmt"
	"net/http"
	"os"

	"github.com/spf13/cobra"
)


var serverURL string

	// Root command
var monitor = &cobra.Command{
	Use:   "repo-sender",
	Short: "CLI tool to send GitHub repo links to a server",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 1 {
			fmt.Println("Error: No GitHub repository link provided.")
			os.Exit(1)
		}

		repoLink := args[0]
		if err := sendRepoLink(serverURL, repoLink); err != nil {
			fmt.Printf("Failed to send repo link: %v\n", err)
			os.Exit(1)
		}

		fmt.Println("Repository link sent successfully!")
	},
}
func init(){
	// Add a flag for the server URL
	rootCmd.Flags().StringVarP(&serverURL, "server", "s", "http://localhost:8080", "Server URL to send the repo link")
	rootCmd.AddCommand(monitor)
}

// Function to send the repository link to the server
func sendRepoLink(serverURL, repoLink string) error {
	// Create the JSON payload
	payload := fmt.Sprintf(`{"repo": "%s"}`, repoLink)

	// Send POST request to the server
	resp, err := http.Post(serverURL+"/repos", "application/json", bytes.NewBuffer([]byte(payload)))
	if err != nil {
		return fmt.Errorf("error sending POST request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("server responded with status: %s", resp.Status)
	}

	return nil
}
