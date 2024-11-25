package cmd

import (
	"fmt"
	"os"

	"github.com/gookit/color"
	"github.com/spf13/cobra"
)


var rootCmd = &cobra.Command{
	Use:   "docker-vuln",
	Short: "A CLI tool for scanning Docker images for vulnerabilities",
	Long:  `docker-vuln-scanner is a CLI tool to scan Docker images layer by layer to detect vulnerabilities using cve-bin-tool.`,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {		
		fmt.Println()
		printArt()	
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func printArt(){
	response := ` 
 _____                _                                       _        
|  __ \              | |                                     | |       
| |  | |  ___    ___ | | __  ___  _ __  ______ __   __ _   _ | | _ __  
| |  | | / _ \  / __|| |/ / / _ \| '__||______|\ \ / /| | | || || '_ \ 
| |__| || (_) || (__ |   < |  __/| |            \ V / | |_| || || | | |
|_____/  \___/  \___||_|\_\ \___||_|             \_/   \__,_||_||_| |_|
                                                                                                                                             
 `
	color.Style{color.FgWhite,color.OpBold}.Println(response)
}
func init() {
	
}
