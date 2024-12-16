package cmd

import (
	"context"
	"bytes"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
	"github.com/spf13/cobra"
)

var scanCmd = &cobra.Command{
	Use:   "scan",
	Short: "Scan the Docker image for vulnerabilities",
	Long:  "This command accepts a Docker image as input and lists the vulnerabilities found.",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		imageName := args[0]
		outputFile, _ := cmd.Flags().GetString("output")
		fmt.Printf("Scanning Docker image: %s\n", imageName)
		scanImage(imageName, outputFile)
	},
}

func init() {
	scanCmd.Flags().StringP("output", "o", "", "Output file to store vulnerabilities ")
	rootCmd.AddCommand(scanCmd)
}

// Function to scan a Docker image
func scanImage(imageName string, outputFile string) {
	sbomFile := "sbom.json"

	// Generate the SBOM for the image
	err := generateSBOM(imageName, sbomFile)
	if err != nil {
		log.Fatalf("Error generating SBOM: %v", err)
		return
	}

	// If no output file is provided, default to printing the output to the terminal
	if outputFile == "" {
		err = scanSBOM(sbomFile)
	} else {
		err = scanSBOMToFile(sbomFile, outputFile)
	}
	if err != nil {
		log.Fatalf("Error scanning SBOM: %v", err)
		return
	}
}

func generateSBOM(imageName string, sbomFile string) error {
	cmd := exec.CommandContext(context.Background(), "syft", imageName, "-o", "cyclonedx-json")

	file, err := os.Create(sbomFile)
	if err != nil {
		return fmt.Errorf("failed to create SBOM file: %v", err)
	}
	defer file.Close()

	// Set the command's output to the file
	cmd.Stdout = file

	// Start the command
	if err := cmd.Start(); err != nil {
		return fmt.Errorf("failed to start SBOM generation command: %v", err)
	}

	// Wait for the command to finish
	if err := cmd.Wait(); err != nil {
		return fmt.Errorf("command finished with error: %v", err)
	}

	fmt.Println("SBOM generated:", sbomFile)
	return nil
}

// Function to scan the SBOM and print results directly to the terminal
func scanSBOM(sbomFile string) error {
	cmd := exec.Command("osv-scanner", "--sbom", sbomFile)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("error running OSV-Scanner: %v\nOutput: %s", err, string(output))
	}
	fmt.Printf("Scan results:\n%s\n", string(output))
	return nil
}

func scanSBOMToFile(sbomFile, outputFile string) error {
	var outputBuffer bytes.Buffer

	// Run the OSV-Scanner command and capture its output
	cmd := exec.Command("osv-scanner", "--sbom", sbomFile)
	cmd.Stdout = &outputBuffer
	cmd.Stderr = &outputBuffer

	err := cmd.Run()
	if err != nil {
		fmt.Printf("Warning: error running OSV-Scanner: %v\n", err)
	}

	// Get the content of the output buffer as a string
	output := outputBuffer.String()

	// Extract the table portion of the output
	start := strings.Index(output, "+--------------------------------+")
	end := strings.LastIndex(output, "+--------------------------------+")
	var table string
	if start != -1 && end != -1 && end > start {
		table = output[start : end+len("+--------------------------------+")]
	} else {
		table =  output
	}

	// Open or create the output file for writing
	file, err := os.Create(outputFile)
	if err != nil {
		return fmt.Errorf("failed to create output file: %v", err)
	}
	defer file.Close()

	// Write the extracted table to the file
	_, err = file.WriteString(table)
	if err != nil {
		return fmt.Errorf("failed to write table to file: %v", err)
	}

	// Confirm that the results were saved
	fmt.Printf("Scan results saved to %s\n", outputFile)
	return nil
}
