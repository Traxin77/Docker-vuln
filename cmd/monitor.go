package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/spf13/cobra"
)
var serverURL string

var monitorCmd = &cobra.Command{
	Use:   "monitor [repo link]",
	Short: "Send GitHub repo links to the server for monitoring",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		repoLink := args[0]
		if err := sendRepoLink(serverURL, repoLink); err != nil {
			fmt.Printf("Failed to send repo link: %v\n", err)
			os.Exit(1)
		}

		fmt.Println("Repository link sent successfully!")
	},
}

func init() {
	// Add server flag as a persistent flag
	monitorCmd.Flags().StringVarP(&serverURL, "server", "s", "http://localhost:8080", "Server URL to send the repo link")
	// Register monitor command to the root command
	rootCmd.AddCommand(monitorCmd)
}

// Sends the repository link to the server
func sendRepoLink(serverURL, repoLink string) error {
	payload, err := json.Marshal(map[string]string{"repo": repoLink})
	if err != nil {
		return fmt.Errorf("error creating JSON payload: %w", err)
	}

	resp, err := http.Post(serverURL+"/repos", "application/json", bytes.NewBuffer(payload))
	if err != nil {
		return fmt.Errorf("error sending POST request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("server responded with status: %s, message: %s", resp.Status, string(body))
	}

	return nil
}
