package cmd

import (
	"context"
	"fmt"

	"github.com/docker/docker/api/types/image"
	"github.com/docker/docker/client"
	"github.com/spf13/cobra"
)

// listImagesCmd represents the command to list Docker images
var listImagesCmd = &cobra.Command{
	Use:   "list",
	Short: "List all Docker images on the system",
	Long:  `This command lists all the Docker images available on the system, showing repository, tags, and image IDs.`,
	Run: func(cmd *cobra.Command, args []string) {
		listImages()
	},
}

func init() {
	rootCmd.AddCommand(listImagesCmd)
}

func listImages() {
	// Initialize Docker client
	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		fmt.Println("Error initializing Docker client:", err)
		return
	}
	cli.NegotiateAPIVersion(context.Background())

	// List all Docker images
	images, err := cli.ImageList(context.Background(), image.ListOptions{})
	if err != nil {
		fmt.Println("Error fetching Docker images:", err)
		return
	}

	fmt.Println("Listing all Docker images:")
	for _, image := range images {
	
		repoTags := "none"
		if len(image.RepoTags) > 0 {
			repoTags = image.RepoTags[0] 
		}
		fmt.Printf("Repository/Tag: %s | Image ID: %s | Size: %.2f MB\n", repoTags, image.ID[:12], float64(image.Size)/(1024*1024))
	}
}
