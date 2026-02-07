# URL Tools: Extractor & Website Health Checker

A professional Golang toolkit to **extract URLs** from websites and **verify their status**. Designed for performance, safety, and ease of use.

> **Safety First:** This tool automatically respects `robots.txt` rules to ensure ethical crawling and prevent IP bans.

---

## Project Structure

```text
url-extractor/
├── cmd/
│   ├── url-extractor/   # Tool: Extract URLs from HTML
│   │   └── main.go
│   └── url-checker/     # Tool: Check if URLs are online
│       └── main.go
├── src/
│   ├── extractor/       # Extraction logic
│   ├── fetcher/         # HTTP request handling
│   ├── checker/         # Status checking logic
│   ├── robots/          # Robots.txt compliance engine
│   └── models/          # Shared data structures
├── urldata/             # Directory for saving JSON results
├── go.mod               # Go module definition
└── README.md            # Documentation
```

---

## Quick Start Guide

### Step 1: Installation
Open your terminal and build the tools:

```bash
# Build the Extractor
go build -o url-extractor.exe ./cmd/url-extractor

# Build the Checker
go build -o url-checker.exe ./cmd/url-checker
```

### Step 2: Extract URLs
Run the extractor on a website you want to analyze. This will verify `robots.txt` first, then save all found links to a file.

```bash
# Syntax: ./url-extractor.exe <URL> > <OUTPUT_FILE>
./url-extractor.exe https://example.com > urldata/my-scan.json
```

### Step 3: Check Status
Feed the result file into the checker to see which links are broken.

```bash
# Syntax: ./url-checker.exe <INPUT_FILE>
./url-checker.exe urldata/my-scan.json
```

---

## Key Features

| Feature | Description |
| :--- | :--- |
| **Smart Extraction** | Finds links in `href` and `src`, converting relative paths to absolute ones. |
| **Ethical Crawling** | Checks `robots.txt` before every request. Caches rules to protect servers. |
| **High Performance** | Uses Go concurrency to check hundreds of URLs in seconds. |
| **Clean Output** | Produces strict JSON for easy integration with other tools. |
| **Visual Reports** | Color-coded CLI output makes it easy to spot errors and broken links. |

---

## Advanced Usage

### Customizing the Checker
You can control how fast the checker works to avoid overloading servers.

```bash
# Usage: url-checker <file> [workers] [timeout]
./url-checker.exe urldata/my-scan.json 10 5
```
*   **workers**: Number of simultaneous checks (Default: 5). Higher = Faster but more load.
*   **timeout**: Seconds to wait for a response (Default: 10).

### Handling Errors & Logs
The extractor prints strict JSON to **Standard Output (stdout)** and logs/errors to **Standard Error (stderr)**. You can separate them easily:

```bash
# Save JSON to file, and errors to a log file
./url-extractor.exe https://google.com > urldata/google.json 2> urldata/debug.log
```

---

## JSON Output Format
The extractor produces JSON in the following format:

```json
{
  "source_url": "https://example.com",
  "extracted_urls": [
    "https://example.com/about",
    "https://example.com/contact"
  ],
  "count": 2
}
```

---

## Requirements
- **Go 1.21** or higher installed.
- Internet connection for fetching and checking live URLs.

---

## FAQ

**Q: Why does it say "disallowed by robots.txt"?**
A: The tool respects the website's rules. If a site owner blocks crawlers from a specific path (e.g., `/admin`), this tool will skip it to keep you safe from bans.

**Q: Can I check thousands of URLs?**
A: Yes! The checker is built for speed. Just be mindful of the `workers` setting so you don't accidentally attack a server (DDoS).

