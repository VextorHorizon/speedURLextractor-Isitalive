package main

import (
	"encoding/json"
	"fmt"
	"net/url"
	"os"

	"url-extractor/src/extractor"
	"url-extractor/src/fetcher"
	"url-extractor/src/models"
	"url-extractor/src/robots"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: url-extractor <URL>")
		fmt.Println("Example: url-extractor https://example.com")
		os.Exit(1)
	}

	targetURL := os.Args[1]

	// Validate URL
	_, err := url.ParseRequestURI(targetURL)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: Invalid URL: %v\n", err)
		os.Exit(1)
	}

	// Check robots.txt
	fmt.Fprintf(os.Stderr, "Checking robots.txt for: %s\n", targetURL)
	allowed, err := robots.IsAllowed(targetURL)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Warning: Could not check robots.txt: %v. Proceeding with caution.\n", err)
	} else if !allowed {
		fmt.Fprintf(os.Stderr, "Error: Access to this URL is disallowed by robots.txt\n")
		os.Exit(1)
	}

	// Fetch content
	fmt.Fprintf(os.Stderr, "Fetching content from: %s\n", targetURL)
	content, err := fetcher.FetchURL(targetURL)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error fetching URL: %v\n", err)
		os.Exit(1)
	}

	// Extract URLs
	fmt.Fprintf(os.Stderr, "Extracting URLs...\n")
	extractedURLs := extractor.ExtractURLs(content, targetURL)

	// Create result
	result := models.URLResult{
		SourceURL:     targetURL,
		ExtractedURLs: extractedURLs,
		Count:         len(extractedURLs),
	}

	// Output as JSON
	jsonOutput, err := json.MarshalIndent(result, "", "  ")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error creating JSON: %v\n", err)
		os.Exit(1)
	}

	fmt.Println(string(jsonOutput))
	fmt.Fprintf(os.Stderr, "\nTotal URLs found: %d\n", result.Count)
}
