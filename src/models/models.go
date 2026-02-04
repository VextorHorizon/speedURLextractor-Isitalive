package models

// URLResult represents the result of URL extraction
type URLResult struct {
	SourceURL     string   `json:"source_url"`
	ExtractedURLs []string `json:"extracted_urls"`
	Count         int      `json:"count"`
}

// URLStatus represents the status of a URL check
type URLStatus struct {
	URL        string
	StatusCode int
	Status     string
	Online     bool
	Error      error
}
