# Ujeebu API Go SDK

[![Go Reference](https://pkg.go.dev/badge/github.com/ujeebu/ujeebu-go.svg)](https://pkg.go.dev/github.com/ujeebu/ujeebu-go)
[![Go Report Card](https://goreportcard.com/badge/github.com/ujeebu/ujeebu-go)](https://goreportcard.com/report/github.com/ujeebu/ujeebu-go)

[Ujeebu](https://ujeebu.com) is a comprehensive API platform for web scraping, content extraction, and search engine results. This SDK provides a robust, production-ready interface for Go applications with support for:

- 🔄 **Context-based cancellation** - Full support for context.Context in all API calls
- ⚡ **Options pattern** - Flexible client configuration
- 🛡️ **Strong error handling** - Structured error types with helper methods
- 📊 **Credit tracking** - Monitor API usage in real-time
- 🧪 **Comprehensive testing** - Well-tested with extensive unit test coverage

## Table of Contents

- [Installation](#installation)
- [Quick Start](#quick-start)
- [Client Configuration](#client-configuration)
- [API Endpoints](#api-endpoints)
  - [Extract API](#extract-api)
  - [Card API](#card-api)
  - [Scrape API](#scrape-api)
  - [SERP API](#serp-api)
  - [Account API](#account-api)
- [Advanced Usage](#advanced-usage)
  - [Context Support](#context-support)
  - [Error Handling](#error-handling)
  - [Custom Headers](#custom-headers)
  - [Proxy Support](#proxy-support)
  - [Retry Configuration](#retry-configuration)
- [Examples](#examples)
- [Testing](#testing)
- [Contributing](#contributing)
- [License](#license)

## Installation

```bash
go get github.com/ujeebu/ujeebu-go
```

## Quick Start

```go
package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/ujeebu/ujeebu-go"
)

func main() {
	// Create a new client with your API key
	client, err := ujeebu.NewClient("YOUR-API-KEY")
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}

	// Extract article content
	params := ujeebu.ExtractParams{
		URL: "https://ujeebu.com/blog/web-scraping-in-2025",
	}

	article, credits, err := client.Extract(params)
	if err != nil {
		log.Fatalf("Extraction failed: %v", err)
	}

	fmt.Printf("Title: %s\n", article.Title)
	fmt.Printf("Author: %s\n", article.Author)
	fmt.Printf("Credits used: %d\n", credits)
}
```

## Client Configuration

The SDK uses an options pattern for flexible client configuration:

```go
import (
	"time"
	"github.com/ujeebu/ujeebu-go"
)

client, err := ujeebu.NewClient(
	"YOUR-API-KEY",
	ujeebu.WithTimeout(120 * time.Second),              // Custom timeout
	ujeebu.WithBaseURL("https://custom.api.url"),       // Custom API URL
	ujeebu.WithDebug(true),                             // Enable debug logging
	ujeebu.WithUserAgent("MyApp/1.0"),                  // Custom user agent
	ujeebu.WithRetry(3, time.Second, 5*time.Second),    // Retry configuration
)
if err != nil {
	log.Fatalf("Failed to create client: %v", err)
}
```

### Available Options

- `WithTimeout(duration)` - Set custom HTTP client timeout (default: 90 seconds)
- `WithBaseURL(url)` - Set custom API base URL (default: https://api.ujeebu.com)
- `WithDebug(bool)` - Enable/disable debug mode with detailed logging
- `WithLogger(logger)` - Set custom logger implementing the Logger interface
- `WithUserAgent(ua)` - Set custom User-Agent header
- `WithRetry(maxRetries, waitTime, maxWaitTime)` - Configure retry behavior

### Environment Variables

- `UJEEBU_BASE_URL` - Override the default API base URL

## API Endpoints

### Extract API

The Extract API automatically extracts clean, structured content from web pages.

#### Basic Usage

```go
params := ujeebu.ExtractParams{
	URL: "https://example.com/article",
}

article, credits, err := client.Extract(params)
if err != nil {
	log.Fatalf("Extraction failed: %v", err)
}

fmt.Printf("Title: %s\n", article.Title)
fmt.Printf("Text: %s\n", article.Text)
fmt.Printf("Author: %s\n", article.Author)
fmt.Printf("Publication Date: %s\n", article.PubDate)
```

#### Advanced Options

```go
params := ujeebu.ExtractParams{
	URL:           "https://example.com/article",
	JS:            true,                    // Enable JavaScript rendering
	JSTimeout:     10000,                   // JS execution timeout (ms)
	Text:          true,                    // Extract text content
	HTML:          true,                    // Include HTML
	Images:        true,                    // Extract images
	Author:        true,                    // Extract author information
	PubDate:       true,                    // Extract publication date
	Pagination:    true,                    // Handle multi-page articles
	PaginationMaxPages: 5,                  // Maximum pages to extract
	ProxyType:     "premium",               // Use premium proxies
	ProxyCountry:  "us",                    // Proxy location
	SessionID:     "my-session-123",        // Session for cookie persistence
	CustomHeaders: map[string]string{       // Custom headers (prefixed with UJB-)
		"Authorization": "Bearer token",
	},
}

article, credits, err := client.Extract(params)
```

#### Extract from Raw HTML

```go
params := ujeebu.ExtractParams{
	URL:    "https://example.com/article",
	RawHTML: "<html><body><h1>Title</h1><p>Content</p></body></html>",
}

article, credits, err := client.Extract(params)
```

### Card API

The Card API quickly retrieves metadata and preview information from URLs, optimized for social media cards and link previews.

#### Basic Usage

```go
params := ujeebu.CardParams{
	URL: "https://example.com/article",
}

card, credits, err := client.Card(params)
if err != nil {
	log.Fatalf("Card fetch failed: %v", err)
}

fmt.Printf("Title: %s\n", card.Title)
fmt.Printf("Summary: %s\n", card.Summary)
fmt.Printf("Image: %s\n", card.Image)
fmt.Printf("Favicon: %s\n", card.Favicon)
fmt.Printf("Site Name: %s\n", card.SiteName)
```

#### Advanced Options

```go
params := ujeebu.CardParams{
	URL:          "https://example.com",
	JS:           true,                    // Enable JavaScript for SPAs
	JSTimeout:    5000,                    // JS timeout (ms)
	Timeout:      30,                      // Request timeout (seconds)
	ProxyType:    "datacenter",            // Proxy type
	ProxyCountry: "us",                    // Proxy country
	CustomProxy:  "http://proxy:8080",    // Custom proxy URL
	AutoProxy:    true,                    // Automatic proxy selection
	SessionID:    "session-123",           // Session ID for cookies
	CustomHeaders: map[string]string{      // Custom headers
		"Accept-Language": "en-US",
	},
}

card, credits, err := client.Card(params)
```

### Scrape API

The Scrape API provides full control over web page scraping with JavaScript rendering, screenshots, PDFs, and custom extraction rules.

#### Basic Scraping

```go
params := ujeebu.ScrapeParams{
	URL: "https://example.com",
}

response, credits, err := client.Scrape(params)
if err != nil {
	log.Fatalf("Scraping failed: %v", err)
}

fmt.Printf("Success: %v\n", response.Success)
fmt.Printf("HTML: %s\n", response.HTML)
fmt.Printf("Status Code: %d\n", response.StatusCode)
fmt.Printf("Content-Type: %s\n", response.ResponseHeaders.Get("Content-Type"))
```

#### Screenshot Capture

```go
params := ujeebu.ScrapeParams{
	URL: "https://example.com",
}

// Full page screenshot
screenshot, credits, err := client.Screenshot(params, true, "")

// Partial screenshot with CSS selector
screenshot, credits, err := client.Screenshot(params, false, ".main-content")
```

#### PDF Generation

```go
params := ujeebu.ScrapeParams{
	URL: "https://example.com",
}

pdf, credits, err := client.PDF(params)
```

#### Advanced Scraping with JavaScript

```go
params := ujeebu.ScrapeParams{
	URL:            "https://example.com",
	JS:             true,                     // Enable JavaScript
	JSTimeout:      15000,                    // JS execution timeout
	WaitFor:        ".products-loaded",       // Wait for element
	WaitForTimeout: 10000,                    // Wait timeout
	ScrollDown:     true,                     // Auto-scroll
	ScrollWait:     2000,                     // Wait between scrolls
	ProgressiveScroll: true,                  // Progressive scrolling
	CustomJS:       "document.body.click()",  // Custom JavaScript
	Device:         "mobile",                 // Device type
	WindowWidth:    375,                      // Window width
	WindowHeight:   667,                      // Window height
	BlockAds:       true,                     // Block advertisements
	BlockResources: true,                     // Block images/fonts/etc
	ProxyType:      "premium",                // Premium proxies
	ProxyCountry:   "us",                     // Proxy location
	CustomHeaders: map[string]string{         // Custom headers
		"Accept-Language": "en-US",
	},
}

response, credits, err := client.Scrape(params)
```

#### Extraction Rules

```go
params := ujeebu.ScrapeParams{
	URL: "https://example.com/products",
	ExtractRules: map[string]any{
		"products": map[string]any{
			"_selector": ".product",
			"title":     ".product-title",
			"price":     ".product-price",
			"image":     map[string]string{
				"_selector":  ".product-image",
				"_attribute": "src",
			},
		},
	},
}

response, credits, err := client.Scrape(params)
// Access extracted data via response.Result
```

### SERP API

The SERP (Search Engine Results Page) API provides Google search results in structured format.

#### Google Web Search

```go
params := ujeebu.SerpParams{
	Search:       "golang web scraping",
	SearchType:   "text",                    // Optional: defaults to text
	Lang:         "en",                      // Language
	Location:     "us",                      // Location
	Device:       "desktop",                 // Device type
	ResultsCount: 10,                        // Results per page
	Page:         1,                         // Page number
}

results, credits, err := client.GoogleSearch(params)
if err != nil {
	log.Fatalf("Search failed: %v", err)
}

// Access organic results
for _, result := range results.OrganicResults {
	fmt.Printf("Title: %s\n", result.Title)
	fmt.Printf("Link: %s\n", result.Link)
	fmt.Printf("Description: %s\n", result.Description)
}

// Access knowledge graph
if results.KnowledgeGraph.Title != "" {
	fmt.Printf("Knowledge Graph: %s\n", results.KnowledgeGraph.Title)
}
```

#### Google News Search

```go
params := ujeebu.SerpParams{
	Search:   "tech news",
	Location: "us",
	Lang:     "en",
}

results, credits, err := client.GoogleNewsSearch(params)
for _, news := range results.News {
	fmt.Printf("Title: %s\n", news.Title)
	fmt.Printf("Source: %s\n", news.SiteName)
	fmt.Printf("Date: %s\n", news.Date)
	fmt.Printf("Link: %s\n", news.Link)
}
```

#### Google Video Search

```go
params := ujeebu.SerpParams{
	Search: "golang tutorial",
	Lang:   "en",
}

results, credits, err := client.GoogleVideoSearch(params)
for _, video := range results.Videos {
	fmt.Printf("Title: %s\n", video.Title)
	fmt.Printf("Provider: %s\n", video.Provider)
	fmt.Printf("URL: %s\n", video.Url)
}
```

#### Google Image Search

```go
params := ujeebu.SerpParams{
	Search: "golang logo",
	Lang:   "en",
}

results, credits, err := client.GoogleImageSearch(params)
for _, image := range results.Images {
	fmt.Printf("Title: %s\n", image.Title)
	fmt.Printf("Source: %s\n", image.Source)
	fmt.Printf("Image URL: %s\n", image.Image)
	fmt.Printf("Size: %dx%d\n", image.Width, image.Height)
}
```

#### Google Maps Search

```go
params := ujeebu.SerpParams{
	Search:   "restaurants near me",
	Location: "us",
	Lang:     "en",
}

results, credits, err := client.GoogleMapSearch(params)
for _, place := range results.Maps {
	fmt.Printf("Name: %s\n", place.Title)
	fmt.Printf("Address: %s\n", place.Address)
	fmt.Printf("Rating: %.1f (%d reviews)\n", place.Rating, place.Reviews)
	fmt.Printf("Category: %s\n", place.Category)
}
```

### Account API

Check your account usage and billing information.

```go
account, err := client.Account()
if err != nil {
	log.Fatalf("Failed to get account info: %v", err)
}

fmt.Printf("Plan: %s\n", account.Plan)
fmt.Printf("Quota: %s\n", account.Quota)
fmt.Printf("Used: %d (%.1f%%)\n", account.Used, account.UsedPercent)
fmt.Printf("Balance: %d credits\n", account.Balance)
fmt.Printf("Concurrent Requests: %d\n", account.ConcurrentRequests)

if account.NextBillingDate != nil {
	fmt.Printf("Next Billing: %s\n", *account.NextBillingDate)
	fmt.Printf("Days until billing: %d\n", account.DaysTillNextBilling)
}
```

## Advanced Usage

### Context Support

All API methods support context for cancellation and timeouts:

```go
import "context"

// Create a context with timeout
ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
defer cancel()

// Use context with API calls
article, credits, err := client.ExtractWithContext(ctx, params)
if err != nil {
	// Handle timeout or cancellation
	if ctx.Err() == context.DeadlineExceeded {
		log.Println("Request timed out")
	}
}
```

All endpoints have WithContext variants:
- `ExtractWithContext(ctx, params)`
- `CardWithContext(ctx, params)`
- `ScrapeWithContext(ctx, params)`
- `SerpWithContext(ctx, params)`
- `AccountWithContext(ctx)`

### Error Handling

The SDK provides structured error types for better error handling:

```go
article, credits, err := client.Extract(params)
if err != nil {
	// Check error type
	var apiErr *ujeebu.APIError
	var validationErr *ujeebu.ValidationError
	var netErr *ujeebu.NetworkError

	switch {
	case errors.As(err, &apiErr):
		// API error with status code and message
		fmt.Printf("API Error: %s (status: %d)\n", apiErr.Message, apiErr.StatusCode)

		// Use helper methods
		if apiErr.IsUnauthorized() {
			log.Println("Invalid API key")
		}
		if apiErr.IsNotFound() {
			log.Println("Resource not found")
		}
		if apiErr.IsRateLimited() {
			log.Println("Rate limit exceeded")
		}
		if apiErr.IsTimeout() {
			log.Println("Request timed out")
		}

	case errors.As(err, &validationErr):
		// Validation error (client-side)
		fmt.Printf("Validation error for %s: %s\n", validationErr.Field, validationErr.Message)

	case errors.As(err, &netErr):
		// Network error (connection issues, timeouts)
		fmt.Printf("Network error: %v\n", netErr.Err)
	}
}
```

### Custom Headers

Add custom headers to requests (they will be prefixed with `UJB-`):

```go
params := ujeebu.ExtractParams{
	URL: "https://example.com",
	CustomHeaders: map[string]string{
		"Authorization": "Bearer token",
		"X-Custom-ID":   "12345",
	},
}

// Headers sent to Ujeebu API:
// UJB-Authorization: Bearer token
// UJB-X-Custom-ID: 12345
```

### Proxy Support

The SDK supports various proxy types:

```go
params := ujeebu.ScrapeParams{
	URL:          "https://example.com",
	ProxyType:    "premium",           // premium, datacenter, residential
	ProxyCountry: "us",                // Proxy country code
	CustomProxy:  "http://user:pass@proxy.example.com:8080", // Custom proxy
	AutoProxy:    true,                // Automatic proxy selection
}
```

### Retry Configuration

Configure automatic retries for failed requests:

```go
client, err := ujeebu.NewClient(
	"YOUR-API-KEY",
	ujeebu.WithRetry(
		3,                      // Max retries
		1*time.Second,          // Initial wait time
		10*time.Second,         // Max wait time
	),
)
```

## Examples

Complete examples are available in the `examples/` directory:

- [Extract Article](examples/extract/main.go) - Basic article extraction
- [Card Preview](examples/card/main.go) - Get page metadata
- [Scrape with JavaScript](examples/scrape/main.go) - JavaScript rendering and scrolling
- [Screenshot](examples/screenshot/main.go) - Capture page screenshots
- [PDF Generation](examples/pdf/main.go) - Generate PDFs
- [Google Search](examples/serp/main.go) - Search Google and parse results
- [Account Info](examples/account/main.go) - Check account usage

## Testing

The SDK includes comprehensive unit tests with >90% coverage:

```bash
# Run all tests
go test ./...

# Run tests with coverage
go test -cover ./...

# Run tests with verbose output
go test -v ./...

# Run specific test
go test -run TestCard_Success
```

## Error Reference

### APIError

Server-side errors from the Ujeebu API:
- **StatusCode**: HTTP status code
- **Message**: Error message
- **ErrorCode**: API error code
- **URL**: Request URL

Helper methods:
- `IsUnauthorized()` - 401 status code
- `IsNotFound()` - 404 status code
- `IsRateLimited()` - 429 status code
- `IsTimeout()` - 408 or 504 status code

### ValidationError

Client-side validation errors:
- **Field**: Field name that failed validation
- **Message**: Validation error message

### NetworkError

Network-level errors:
- **Err**: Underlying error (timeout, connection refused, etc.)

## Best Practices

1. **Always check errors**: Handle errors appropriately for production use
2. **Use contexts**: Implement timeouts and cancellation for long-running requests
3. **Monitor credits**: Track credit usage via the returned credits value
4. **Reuse clients**: Create one client instance and reuse it across requests
5. **Handle rate limits**: Implement exponential backoff for rate-limited requests
6. **Set appropriate timeouts**: Adjust timeouts based on your use case
7. **Use sessions**: For multi-request workflows, use SessionID for cookie persistence

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## Support

- 📧 Email: support@ujeebu.com
- 📚 Documentation: https://ujeebu.com/docs
- 🐛 Issues: https://github.com/ujeebu/ujeebu-go/issues

## License

This SDK is distributed under the MIT License. See LICENSE file for more information.
