package robots

import (
	"bufio"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"url-extractor/src/fetcher"
)

// Policy represents the robots.txt policy for a domain
type Policy struct {
	DisallowedPaths []string
	IsFetched       bool
}

var (
	cache = make(map[string]*Policy)
	mu    sync.RWMutex
)

// IsAllowed checks if a URL is allowed by robots.txt
func IsAllowed(targetURL string) (bool, error) {
	u, err := url.Parse(targetURL)
	if err != nil {
		return false, err
	}

	domain := u.Host
	path := u.Path
	if path == "" {
		path = "/"
	}

	policy := getPolicy(domain)
	if policy == nil {
		policy, err = fetchPolicy(u.Scheme, domain)
		if err != nil {
			// If we can't fetch robots.txt, we assume it's allowed (or you can be stricter)
			// A friendly crawler usually allows if robots.txt is missing (404)
			return true, nil
		}
		savePolicy(domain, policy)
	}

	for _, disallowed := range policy.DisallowedPaths {
		if disallowed == "" {
			continue
		}
		// Strict check: if path starts with disallowed path
		if strings.HasPrefix(path, disallowed) {
			return false, nil
		}
	}

	return true, nil
}

func getPolicy(domain string) *Policy {
	mu.RLock()
	defer mu.RUnlock()
	return cache[domain]
}

func savePolicy(domain string, policy *Policy) {
	mu.Lock()
	defer mu.Unlock()
	cache[domain] = policy
}

func fetchPolicy(scheme, domain string) (*Policy, error) {
	robotsURL := fmt.Sprintf("%s://%s/robots.txt", scheme, domain)

	client := &http.Client{}
	req, err := http.NewRequest("GET", robotsURL, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", fetcher.DefaultUserAgent)

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return &Policy{IsFetched: true}, nil
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("status: %d", resp.StatusCode)
	}

	policy := &Policy{IsFetched: true}
	scanner := bufio.NewScanner(resp.Body)

	relevantUserAgent := false
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		parts := strings.SplitN(line, ":", 2)
		if len(parts) < 2 {
			continue
		}

		key := strings.ToLower(strings.TrimSpace(parts[0]))
		val := strings.TrimSpace(parts[1])

		if key == "user-agent" {
			// We check for wildcard or our specific user-agent (if we had a specific name)
			// Since we use a Chrome UA, we mostly care about '*'
			if val == "*" {
				relevantUserAgent = true
			} else {
				relevantUserAgent = false
			}
		}

		if relevantUserAgent && key == "disallow" {
			policy.DisallowedPaths = append(policy.DisallowedPaths, val)
		}
	}

	return policy, nil
}
