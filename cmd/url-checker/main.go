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
// printHeader prints the header for the output
func printHeader() {
	fmt.Printf("\n%s╔════════════════════════════════════════════════════════════════════════════╗%s\n", ColorCyan, ColorReset)
	fmt.Printf("%s║                         URL STATUS CHECK RESULTS                          ║%s\n", ColorCyan, ColorReset)
	fmt.Printf("%s╚════════════════════════════════════════════════════════════════════════════╝%s\n\n", ColorCyan, ColorReset)
}
// printSummary prints the final summary
func printSummary(total, online, offline int) {
	fmt.Printf("\n%s╔════════════════════════════════════════════════════════════════════════════╗%s\n", ColorCyan, ColorReset)
	fmt.Printf("%s║                                  SUMMARY                                   ║%s\n", ColorCyan, ColorReset)
	fmt.Printf("%s╚════════════════════════════════════════════════════════════════════════════╝%s\n", ColorCyan, ColorReset)
	fmt.Printf("Total URLs Checked: %s%d%s\n", ColorBlue, total, ColorReset)
	fmt.Printf("Online:             %s%d%s (%.1f%%)\n", ColorGreen, online, ColorReset, float64(online)/float64(total)*100)
	fmt.Printf("Offline:            %s%d%s (%.1f%%)\n", ColorRed, offline, ColorReset, float64(offline)/float64(total)*100)
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
	resultsChan := checker.CheckURLsConcurrently(urlData.ExtractedURLs, concurrency, timeout)
	printHeader()
	onlineCount := 0
	offlineCount := 0
	count := 0
	// Process results in real-time
	for result := range resultsChan {
		count++
		fmt.Printf("%s[%d]%s ", ColorBlue, count, ColorReset)
		if result.Online {
			onlineCount++
			fmt.Printf("%s✓ ONLINE%s  ", ColorGreen, ColorReset)
			fmt.Printf("[%s%d%s] ", ColorGreen, result.StatusCode, ColorReset)
		} else {
			offlineCount++
			fmt.Printf("%s✗ OFFLINE%s ", ColorRed, ColorReset)
			if result.Error != nil {
				fmt.Printf("[%sERROR%s] ", ColorRed, ColorReset)
			} else {
				fmt.Printf("[%s%d%s] ", ColorYellow, result.StatusCode, ColorReset)
			}
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
	duration := time.Since(startTime)
	// Print summary
	printSummary(count, onlineCount, offlineCount)
	fmt.Printf("⏱  Time elapsed: %s%.2f seconds%s\n", ColorBlue, duration.Seconds(), ColorReset)
}
