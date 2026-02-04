# URL Extractor

A Golang tool to extract all URLs from a given webpage and output them as JSON.

## Project Structure

```
url-extractor/
├── cmd/
│   └── url-extractor/
│       └── main.go      # Entry point
├── src/
│   ├── extractor/       # URL extraction logic
│   │   └── extractor.go
│   ├── fetcher/         # HTTP content fetching
│   │   └── fetcher.go
│   └── models/          # Data structures
│       └── models.go
├── urldata/             # Directory for JSON results
├── go.mod               # Go module file
└── README.md            # Documentation
```

## Features

- ✅ Fetches content from any URL
- ✅ Extracts all URLs from HTML content (href and src attributes)
- ✅ Supports both absolute and relative URLs
- ✅ Resolves relative URLs to absolute URLs
- ✅ Removes duplicate URLs
- ✅ Outputs results as formatted JSON
- ✅ Clean architecture with `cmd/` and `src/` separation

## Installation

```bash
# Build the executable
go build -o url-extractor.exe ./cmd/url-extractor
```

## Usage

```bash
# Run directly with go
go run ./cmd/url-extractor/main.go <URL>

# Or use the built binary
./url-extractor.exe <URL>
```

### Examples

```bash
# Extract URLs from a website
go run ./cmd/url-extractor/main.go https://example.com

# Save output to urldata directory
go run ./cmd/url-extractor/main.go https://example.com > urldata/example.json 2>urldata/error.log
```

## Output Format

```json
  "source_url": "https://example.com",
  "extracted_urls": [
    "https://example.com/page1",
    "https://example.com/page2",
    "https://example.com/image.jpg"
  ],
  "count": 3
}
```

## How It Works

1. **Fetch Content**: Downloads the HTML content from the specified URL
2. **Extract URLs**: Uses regular expressions to find all URLs in the content
3. **Resolve URLs**: Converts relative URLs to absolute URLs based on the source URL
4. **Output JSON**: Formats the results as JSON with the source URL, extracted URLs, and count

## Requirements

- Go 1.21 or higher
