package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"strings"
)

type URLResult struct {
	SourceURL     string   `json:"source_url"`
	ExtractedURLs []string `json:"extracted_urls"`
	Count         int      `json:"count"`
}

// extractURLs extracts all URLs from the given content
func extractURLs(content string, baseURL string) []string {
	// Regular expression to match URLs
	// This matches http://, https://, and relative URLs
	urlRegex := regexp.MustCompile(`(?i)(?:href|src)=["']([^"']+)["']|(?:http|https)://[^\s<>"{}|\\^` + "`" + `\[\]]+`)

	matches := urlRegex.FindAllStringSubmatch(content, -1)

	// Use a map to store unique URLs
	urlMap := make(map[string]bool)

	// Parse base URL for resolving relative URLs
	base, err := url.Parse(baseURL)
	if err != nil {
		base = nil
	}

	for _, match := range matches {
		var foundURL string

		// Check if it's from href/src attribute (group 1)
		if len(match) > 1 && match[1] != "" {
			foundURL = match[1]
		} else {
			// Otherwise it's a direct match (group 0)
			foundURL = match[0]
		}

		// Clean up the URL
		foundURL = strings.TrimSpace(foundURL)
		if foundURL == "" || foundURL == "#" {
			continue
		}

		// Try to resolve relative URLs
		if base != nil && !strings.HasPrefix(foundURL, "http://") && !strings.HasPrefix(foundURL, "https://") {
			if strings.HasPrefix(foundURL, "//") {
				// Protocol-relative URL
				foundURL = base.Scheme + ":" + foundURL
			} else if strings.HasPrefix(foundURL, "/") {
				// Absolute path
				foundURL = base.Scheme + "://" + base.Host + foundURL
			} else if !strings.HasPrefix(foundURL, "mailto:") && !strings.HasPrefix(foundURL, "tel:") && !strings.HasPrefix(foundURL, "javascript:") {
				// Relative path
				resolvedURL, err := base.Parse(foundURL)
				if err == nil {
					foundURL = resolvedURL.String()
				}
			}
		}

		// Skip non-http URLs
		if !strings.HasPrefix(foundURL, "http://") && !strings.HasPrefix(foundURL, "https://") {
			continue
		}

		urlMap[foundURL] = true
	}

	// Convert map to slice
	urls := make([]string, 0, len(urlMap))
	for u := range urlMap {
		urls = append(urls, u)
	}

	return urls
}

// fetchURL fetches content from the given URL
func fetchURL(targetURL string) (string, error) {
	client := &http.Client{}

	req, err := http.NewRequest("GET", targetURL, nil)
	if err != nil {
		return "", fmt.Errorf("failed to create request: %v", err)
	}

	// Add User-Agent to avoid being blocked
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36")

	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to fetch URL: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("HTTP error: %s", resp.Status)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response: %v", err)
	}

	return string(body), nil
}

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

	// Fetch content
	fmt.Fprintf(os.Stderr, "Fetching content from: %s\n", targetURL)
	content, err := fetchURL(targetURL)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error fetching URL: %v\n", err)
		os.Exit(1)
	}

	// Extract URLs
	fmt.Fprintf(os.Stderr, "Extracting URLs...\n")
	extractedURLs := extractURLs(content, targetURL)

	// Create result
	result := URLResult{
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
