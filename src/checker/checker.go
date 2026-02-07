package checker
import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"sync"
	"time"
	"url-extractor/src/fetcher"
	"url-extractor/src/models"
	"url-extractor/src/robots"
)
// CheckURL checks the status of a single URL
func CheckURL(url string, timeout time.Duration) models.URLStatus {
	// Check robots.txt first
	allowed, err := robots.IsAllowed(url)
	if err == nil && !allowed {
		return models.URLStatus{
			URL:    url,
			Online: false,
			Error:  fmt.Errorf("disallowed by robots.txt"),
		}
	}
	client := &http.Client{
		Timeout: timeout,
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			// Allow redirects but limit to 10
			if len(via) >= 10 {
				return fmt.Errorf("too many redirects")
			}
			return nil
		},
	}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return models.URLStatus{
			URL:    url,
			Online: false,
			Error:  err,
		}
	}
	// Add User-Agent to avoid being blocked
	req.Header.Set("User-Agent", fetcher.DefaultUserAgent)
	resp, err := client.Do(req)
	if err != nil {
		return models.URLStatus{
			URL:    url,
			Online: false,
			Error:  err,
		}
	}
	defer resp.Body.Close()
	// Read and discard body to reuse connection
	io.Copy(io.Discard, resp.Body)
	online := resp.StatusCode >= 200 && resp.StatusCode < 400
	return models.URLStatus{
		URL:        url,
		StatusCode: resp.StatusCode,
		Status:     resp.Status,
		Online:     online,
		Error:      nil,
	}
}
// CheckURLsConcurrently checks multiple URLs concurrently and returns a channel of results
func CheckURLsConcurrently(urls []string, concurrency int, timeout time.Duration) <-chan models.URLStatus {
	var wg sync.WaitGroup
	urlChan := make(chan string, len(urls))
	resultChan := make(chan models.URLStatus, len(urls))
	// Worker pool
	for i := 0; i < concurrency; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for url := range urlChan {
				result := CheckURL(url, timeout)
				resultChan <- result
			}
		}()
	}
	// Send URLs to workers
	for _, url := range urls {
		urlChan <- url
	}
	close(urlChan)
	// Wait for all workers to finish
	go func() {
		wg.Wait()
		close(resultChan)
	}()
	return resultChan
}
// StripBOM and convert UTF-16 to UTF-8 if necessary
func StripBOM(data []byte) []byte {
	// Remove UTF-8 BOM if present
	if bytes.HasPrefix(data, []byte("\xef\xbb\xbf")) {
		return data[3:]
	}
	// Remove UTF-16 LE BOM and convert to UTF-8
	if bytes.HasPrefix(data, []byte("\xff\xfe")) {
		return utf16ToUtf8(data[2:])
	}
	// Detect if it's UTF-16 LE without BOM (check for null bytes in common patterns)
	if len(data) > 2 && data[1] == 0 && data[3] == 0 {
		return utf16ToUtf8(data)
	}
	return data
}
// utf16ToUtf8 converts UTF-16 LE bytes to UTF-8
func utf16ToUtf8(data []byte) []byte {
	if len(data)%2 != 0 {
		return data // Not valid UTF-16
	}
	result := make([]byte, 0, len(data)/2)
	for i := 0; i < len(data); i += 2 {
		// This is a simple conversion for basic Latin characters (common in JSON)
		// For a full implementation, one would use unicode/utf16 package
		if data[i+1] == 0 {
			result = append(result, data[i])
		} else {
			// If it's not a basic Latin char, we'd need more complex logic
			// But for URLs and JSON structure, this is usually enough
			// Or we just return the original if it gets too complex
			result = append(result, '?')
		}
	}
	return result
}
