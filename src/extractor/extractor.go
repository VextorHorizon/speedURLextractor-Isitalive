package extractor

import (
	"net/url"
	"regexp"
	"strings"
)

// ExtractURLs extracts all URLs from the given content
// It handles both absolute and relative URLs, resolving relative URLs based on the baseURL
func ExtractURLs(content string, baseURL string) []string {
	// Regular expression to match URLs
	// This matches http://, https://, and relative URLs from href/src attributes
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
