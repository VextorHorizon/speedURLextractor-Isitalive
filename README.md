# URL Extractor

A simple Golang tool to extract all URLs from a given webpage and output them as JSON.

## Features

- ✅ Fetches content from any URL
- ✅ Extracts all URLs from HTML content (href and src attributes)
- ✅ Supports both absolute and relative URLs
- ✅ Resolves relative URLs to absolute URLs
- ✅ Removes duplicate URLs
- ✅ Outputs results as formatted JSON

## Installation

```bash
go build -o url-extractor main.go
```

## Usage

```bash
# Run directly with go
go run main.go <URL>

# Or use the built binary
./url-extractor <URL>
```

### Examples

```bash
# Extract URLs from a website
go run main.go https://example.com

# Save output to file
go run main.go https://example.com > output.json 2>/dev/null
```

## Output Format

```json
{
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
