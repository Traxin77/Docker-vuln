package cmd

import (
	"fmt"
	"os"
	"github.com/spf13/cobra"
)


var rootCmd = &cobra.Command{
	Use:   "docker-vuln",
	Short: "A CLI tool for scanning Docker images for vulnerabilities",
	Long:  `docker-vuln-scanner is a CLI tool to scan Docker images layer by layer to detect vulnerabilities using cve-bin-tool.`,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	
}
