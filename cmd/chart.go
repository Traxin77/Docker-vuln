package cmd

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"

	"github.com/spf13/cobra"
	chart "github.com/wcharczuk/go-chart/v2"
)

var scanWithChartCmd = &cobra.Command{
	Use:   "scan-chart",
	Short: "Scan the Docker image and generate a severity count with visualization",
	Long:  "This command scans a Docker image for vulnerabilities and generates a pie chart visualization of severity levels.",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		inputFile := args[0]
		outputChart, _ := cmd.Flags().GetString("chart")
		fmt.Printf("Parsing scan results from file: %s\n", inputFile)
		criticalityCounts := scanFileWithCriticalityCount(inputFile)

		// Generate pie chart
		if outputChart != "" {
			generatePieChart(criticalityCounts, outputChart)
		}
	},
}

func init() {
	scanWithChartCmd.Flags().StringP("chart", "c", "chart.png", "Output file for pie chart visualization")
	rootCmd.AddCommand(scanWithChartCmd)
}

// Function to parse a scan result file and count vulnerabilities by criticality levels
func scanFileWithCriticalityCount(inputFile string) map[string]int {
	file, err := os.Open(inputFile)
	if err != nil {
		log.Fatalf("Error opening file: %v", err)
	}
	defer file.Close()

	// Regular expressions to extract CVSS scores
	cvssRegex := regexp.MustCompile(`\| (\d+(\.\d+)?)\s+\|`) // Match CVSS scores
	criticalityCounts := map[string]int{
		"low":       0,
		"medium":    0,
		"high":      0,
		"critical":  0,
		"undefined": 0,
	}

	// Read the file line by line
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()

		// Find CVSS score in the line
		match := cvssRegex.FindStringSubmatch(line)
		if len(match) > 1 {
			cvss, err := strconv.ParseFloat(match[1], 64)
			if err != nil {
				log.Printf("Warning: failed to parse CVSS score %v\n", match[1])
				criticalityCounts["undefined"]++
				continue
			}

			// Map CVSS score to criticality
			switch {
			case cvss >= 0.1 && cvss <= 3.9:
				criticalityCounts["low"]++
			case cvss >= 4.0 && cvss <= 6.9:
				criticalityCounts["medium"]++
			case cvss >= 7.0 && cvss <= 8.9:
				criticalityCounts["high"]++
			case cvss >= 9.0 && cvss <= 10.0:
				criticalityCounts["critical"]++
			default:
				criticalityCounts["undefined"]++
			}
		} else if strings.Contains(line, "|") {
			// If line has no CVSS score but contains a vulnerability entry
			criticalityCounts["undefined"]++
		}
	}

	// Check for scanning errors
	if err := scanner.Err(); err != nil {
		log.Fatalf("Error reading file: %v", err)
	}
	return criticalityCounts
}

// Function to generate a pie chart from criticality counts
func generatePieChart(criticalityCounts map[string]int, outputFile string) {
	pie := chart.PieChart{
		Values: []chart.Value{
			{Value: float64(criticalityCounts["low"]), Label: "Low"},
			{Value: float64(criticalityCounts["medium"]), Label: "Medium"},
			{Value: float64(criticalityCounts["high"]), Label: "High"},
			{Value: float64(criticalityCounts["critical"]), Label: "Critical"},
			{Value: float64(criticalityCounts["undefined"]), Label: "Undefined"},
		},
	}

	// Create the chart and save it to a file
	file, err := os.Create(outputFile)
	if err != nil {
		log.Fatalf("Error creating chart file: %v", err)
	}
	defer file.Close()

	err = pie.Render(chart.PNG, file)
	if err != nil {
		log.Fatalf("Error rendering chart: %v", err)
	}

	fmt.Printf("Pie chart saved to %s\n", outputFile)
}
