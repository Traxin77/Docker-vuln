package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
)
var dockerizeCmd = &cobra.Command{
	Use:   "dock",
	Short: "Clone a GitHub repository and build a Docker image from its Dockerfile",
	Long: `The dockerize command clones a GitHub repository, builds a Docker image using the Dockerfile 
present in the repository, and optionally runs a Docker container from the built image.`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		repoURL := args[0]
		repoName := getRepoName(repoURL)
		cloneDir := filepath.Join(os.TempDir(), repoName)

		// Clone the repository
		fmt.Println("Cloning repository...")
		if err := cloneRepo(repoURL, cloneDir); err != nil {
			fmt.Printf("Error cloning repository: %v\n", err)
			return
		}

		defer func() {
			fmt.Println("Cleaning up...")
			if err := os.RemoveAll(cloneDir); err != nil {
				fmt.Printf("Error cleaning up cloned repository: %v\n", err)
			} else {
				fmt.Println("Cloned repository deleted successfully.")
			}
		}()

		// Build the Docker image
		imageName := repoName + ":latest"
		fmt.Println("Building Docker image...")
		if err := buildDockerImage(cloneDir, imageName); err != nil {
			fmt.Printf("Error building Docker image: %v\n", err)
			return
		}

		
		runContainer, _ := cmd.Flags().GetBool("run")
		if runContainer {
			fmt.Println("Running Docker container...")
			if err := runDockerContainer(imageName); err != nil {
				fmt.Printf("Error running Docker container: %v\n", err)
				return
			}
		}

		fmt.Println("Dockerize process completed successfully.")
	},
}

func init() {
	rootCmd.AddCommand(dockerizeCmd)
	dockerizeCmd.Flags().BoolP("run", "r", false, "Run the Docker container after building")
}


func getRepoName(repoURL string) string {
	parts := strings.Split(repoURL, "/")
	return strings.TrimSuffix(parts[len(parts)-1], ".git")
}

func cloneRepo(repoURL, cloneDir string) error {
	cmd := exec.Command("git", "clone", repoURL, cloneDir)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func buildDockerImage(cloneDir, imageName string) error {
	cmd := exec.Command("docker", "build", "-t", imageName, cloneDir)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func runDockerContainer(imageName string) error {
	cmd := exec.Command("docker", "run", "-it", imageName )
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}