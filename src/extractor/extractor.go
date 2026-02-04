package extractor

import (
	"net/url"
	"strings"

	"golang.org/x/net/html"
)

// ExtractURLs extracts all URLs from the given content using an HTML tokenizer
// It handles both absolute and relative URLs, resolving relative URLs based on the baseURL
func ExtractURLs(content string, baseURL string) []string {
	// Use a map to store unique URLs
	urlMap := make(map[string]bool)

	// Parse base URL for resolving relative URLs
	base, err := url.Parse(baseURL)
	if err != nil {
		base = nil
	}

	// Create a tokenizer for the HTML content
	tokenizer := html.NewTokenizer(strings.NewReader(content))

	for {
		tt := tokenizer.Next()

		switch tt {
		case html.ErrorToken:
			// End of document or error
			return mapToSlice(urlMap)
		case html.StartTagToken, html.SelfClosingTagToken:
			token := tokenizer.Token()

			// Check common attributes that contain URLs
			for _, attr := range token.Attr {
				if attr.Key == "href" || attr.Key == "src" {
					processURL(attr.Val, base, urlMap)
				}
				// You can add more attributes here if needed (e.g., "action", "cite", "data-src")
			}
		}
	}
}

// processURL standardizes, resolves, and adds the URL to the map
func processURL(rawURL string, base *url.URL, urlMap map[string]bool) {
	// Clean up the URL
	foundURL := strings.TrimSpace(rawURL)
	if foundURL == "" || foundURL == "#" || strings.HasPrefix(foundURL, "#") {
		return
	}

	// Try to resolve relative URLs
	if base != nil {
		// Handle protocol-relative URLs manually if needed, or rely on ResolveReference
		// url.Parse handles most cases, but let's be careful with parsing the foundURL
		parsedFound, err := url.Parse(foundURL)
		if err == nil {
			resolvedURL := base.ResolveReference(parsedFound)
			foundURL = resolvedURL.String()
		} else {
			// If parsing fails, skip it or attempt raw concatenation (unsafe)
			return
		}
	}

	// Filter protocols: we typically only want http and https
	// But `ResolveReference` might output a clean URL, let's check the scheme
	// This helps filter out javascript: uri schemes etc. if they weren't caught by ResolveReference logic efficiently
	if !strings.HasPrefix(foundURL, "http://") && !strings.HasPrefix(foundURL, "https://") {
		return
	}

	urlMap[foundURL] = true
}

// mapToSlice converts the unique map keys to a slice
func mapToSlice(m map[string]bool) []string {
	urls := make([]string, 0, len(m))
	for u := range m {
		urls = append(urls, u)
	}
	return urls
}
