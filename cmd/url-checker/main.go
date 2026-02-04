package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"url-extractor/src/checker"
	"url-extractor/src/models"
)

// Color codes for CLI output
const (
	ColorReset  = "\033[0m"
	ColorRed    = "\033[31m"
	ColorGreen  = "\033[32m"
	ColorYellow = "\033[33m"
	ColorBlue   = "\033[34m"
	ColorCyan   = "\033[36m"
)

// printResults prints the results in a formatted way
func printResults(results []models.URLStatus) {
	onlineCount := 0
	offlineCount := 0

	fmt.Printf("\n%s╔════════════════════════════════════════════════════════════════════════════╗%s\n", ColorCyan, ColorReset)
	fmt.Printf("%s║                         URL STATUS CHECK RESULTS                          ║%s\n", ColorCyan, ColorReset)
	fmt.Printf("%s╚════════════════════════════════════════════════════════════════════════════╝%s\n\n", ColorCyan, ColorReset)

	for i, result := range results {
		fmt.Printf("%s[%d]%s ", ColorBlue, i+1, ColorReset)

		if result.Online {
			fmt.Printf("%s✓ ONLINE%s  ", ColorGreen, ColorReset)
			fmt.Printf("[%s%d%s] ", ColorGreen, result.StatusCode, ColorReset)
			onlineCount++
		} else {
			fmt.Printf("%s✗ OFFLINE%s ", ColorRed, ColorReset)
			if result.Error != nil {
				fmt.Printf("[%sERROR%s] ", ColorRed, ColorReset)
			} else {
				fmt.Printf("[%s%d%s] ", ColorYellow, result.StatusCode, ColorReset)
			}
			offlineCount++
		}

		// Truncate long URLs for display
		displayURL := result.URL
		if len(displayURL) > 70 {
			displayURL = displayURL[:67] + "..."
		}
		fmt.Printf("%s\n", displayURL)

		// Show error if exists
		if result.Error != nil {
			fmt.Printf("        %sError: %v%s\n", ColorRed, result.Error, ColorReset)
		}
	}

	// Summary
	total := len(results)
	fmt.Printf("\n%s╔════════════════════════════════════════════════════════════════════════════╗%s\n", ColorCyan, ColorReset)
	fmt.Printf("%s║                                  SUMMARY                                   ║%s\n", ColorCyan, ColorReset)
	fmt.Printf("%s╚════════════════════════════════════════════════════════════════════════════╝%s\n", ColorCyan, ColorReset)
	fmt.Printf("Total URLs Checked: %s%d%s\n", ColorBlue, total, ColorReset)
	fmt.Printf("Online:             %s%d%s (%.1f%%)\n", ColorGreen, onlineCount, ColorReset, float64(onlineCount)/float64(total)*100)
	fmt.Printf("Offline:            %s%d%s (%.1f%%)\n", ColorRed, offlineCount, ColorReset, float64(offlineCount)/float64(total)*100)
	fmt.Println()
}

func main() {
	if len(os.Args) < 2 {
		fmt.Printf("%sUsage:%s %s <json-filename-in-urldata> [concurrency] [timeout-seconds]%s\n", ColorYellow, ColorReset, os.Args[0], ColorReset)
		fmt.Printf("%sExample:%s %s reddittest.json 10 5%s\n", ColorYellow, ColorReset, os.Args[0], ColorReset)
		os.Exit(1)
	}

	// Find the file
	filename := os.Args[1]
	jsonPath := filename
	if _, err := os.Stat(jsonPath); os.IsNotExist(err) {
		// Try looking in urldata/
		jsonPath = filepath.Join("urldata", filename)
	}

	concurrency := 5
	timeout := 10 * time.Second

	if len(os.Args) >= 3 {
		fmt.Sscanf(os.Args[2], "%d", &concurrency)
	}

	if len(os.Args) >= 4 {
		var timeoutSec int
		fmt.Sscanf(os.Args[3], "%d", &timeoutSec)
		timeout = time.Duration(timeoutSec) * time.Second
	}

	// Read JSON file
	fmt.Printf("%s📂 Reading JSON file: %s%s\n", ColorCyan, jsonPath, ColorReset)
	data, err := os.ReadFile(jsonPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s❌ Error reading file: %v%s\n", ColorRed, err, ColorReset)
		os.Exit(1)
	}

	// Strip BOM
	data = checker.StripBOM(data)

	// Parse JSON
	var urlData models.URLResult
	err = json.Unmarshal(data, &urlData)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s❌ Error parsing JSON: %v%s\n", ColorRed, err, ColorReset)
		os.Exit(1)
	}

	if len(urlData.ExtractedURLs) == 0 {
		fmt.Fprintf(os.Stderr, "%s⚠ No URLs found in JSON file%s\n", ColorYellow, ColorReset)
		os.Exit(1)
	}

	fmt.Printf("%s✓ Found %d URLs to check%s\n", ColorGreen, len(urlData.ExtractedURLs), ColorReset)
	fmt.Printf("%s🔍 Checking URLs with %d concurrent workers (timeout: %v)...%s\n\n", ColorCyan, concurrency, timeout, ColorReset)

	// Start checking
	startTime := time.Now()
	results := checker.CheckURLsConcurrently(urlData.ExtractedURLs, concurrency, timeout)
	duration := time.Since(startTime)

	// Print results
	printResults(results)

	fmt.Printf("⏱  Time elapsed: %s%.2f seconds%s\n", ColorBlue, duration.Seconds(), ColorReset)
}
