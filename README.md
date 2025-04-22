# Ujeebu API Go SDK

[Ujeebu](https://ujeebu.com) is a set of powerful APIs for Web data scraping and automatic content extraction. This SDK provides an easy-to-use interface for interacting with Ujeebu API in Go applications.

## Installation

You can install the SDK using go get:

```bash
go get github.com/ujeebu/ujeebu-go-sdk
```

## Usage

To use the SDK, you first need to create an instance of it with your API credentials:

```go
package main

import (
	"fmt"
	"log"

	"github.com/ujeebu/ujeebu-go-sdk"
)

func main() {
	client, err := ujeebu.NewClient("YOUR-API-KEY")
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}

	params := ujeebu.ExtractParams{
		URL: "https://ujeebu.com/blog/scraping-javascript-heavy-pages-using-puppeteer/",
	}

	article, credits, err := client.Extract(params)
	if err != nil {
		log.Fatalf("Extraction failed: %v", err)
	}

	fmt.Printf("Article title: %s\n", article.Title)
	fmt.Printf("Credits used: %d\n", credits)
}
```

## APIs

The SDK provides the following methods:

### Scrape API

- `Scrape(params ScrapeParams) (response *ScrapeResponse, credits int, err error)`
  - `params`: A struct containing the scrape API parameters.
  
  The method returns a pointer to a `ScrapeResponse` struct containing the scraped data, the number of credits used, and any error that occurred.

- Helper methods for common scraping operations:
  - `Screenshot(params ScrapeParams, fullPage bool, selector string) (string, int, error)`
  - `PDF(params ScrapeParams) (string, int, error)`
  - `HTML(params ScrapeParams) (string, int, error)`
  - `Raw(params ScrapeParams) (string, int, error)`

### Extract API

- `Extract(params ExtractParams) (article *Article, credits int, err error)`
  - `params`: A struct containing the extract API parameters.
  
  The method returns a pointer to an `Article` struct containing the extracted data, the number of credits used, and any error that occurred.

### SERP API

The SDK also includes support for Google SERP (Search Engine Results Page) APIs:
- Google Search
    - `GoogleSearch(params SerpParams) (GoogleSearchResult, int, error)`
        - `params`: A struct containing parameters such as query, location, language, and pagination options.
        - Returns a `GoogleSearchResult` struct with organic results, knowledge graph, and related information, the number of credits used, and any error that occurred.

- Google News
    - `GoogleNewsSearch(params SerpParams) (GoogleNewsResult, int, error)`
        - `params`: A struct containing parameters such as query, location, language, and pagination options.
        - Returns a `GoogleNewsResult` struct with the fetched news data, the number of credits used, and any error.

- Google Videos
    - `GoogleVideoSearch(params SerpParams) (GoogleVideosResult, int, error)`
        - `params`: A struct containing parameters such as query, location, language, and pagination options.
        - Returns a `GoogleVideosResult` struct with video details, the number of credits used, and any error.

- Google Images
    - `GoogleImageSearch(params SerpParams) (GoogleImagesResult, int, error)`
        - `params`: A struct containing parameters such as query, location, language, and pagination options.
        - Returns a `GoogleImagesResult` struct with image search results, the number of credits used, and any error.

- Google Maps
    - `GoogleMapSearch(params SerpParams) (GoogleMapsResult, int, error)`
        - `params`: A struct containing parameters such as query, location, language, and pagination options.
        - Returns a `GoogleMapsResult` struct with map search results, the number of credits used, and any error.


## Examples

### Scraping a Page with Infinite Scroll

```go
package main

import (
	"fmt"
	"log"

	"github.com/ujeebu/ujeebu-go-sdk"
)

func main() {
	client, err := ujeebu.NewClient("YOUR-API-KEY")
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}

	params := ujeebu.ScrapeParams{
		URL:           "https://ujeebu.com/docs/scrape-me/load-more",
		ResponseType:  "html",
		JSONOutput:    true,
		UserAgent:     "Ujeebu-Go",
		JS:            true,
		WaitFor:       ".products-list",
		WaitForTimeout: 5000,
		ScrollDown:    true,
		ScrollWait:    2000,
		ScrollToSelector: ".load-more-section",
		ScrollCallback: "() => (document.querySelector('.no-more-products') === null)",
		ProxyType:     "premium",
		ProxyCountry:  "US",
		Device:        "desktop",
		WindowWidth:   1200,
		WindowHeight:  900,
		BlockAds:      true,
		CustomHeaders: map[string]string{
			"Authorization": "Basic XXXX",
		},
	}

	response, credits, err := client.Scrape(params)
	if err != nil {
		log.Fatalf("Scraping failed: %v", err)
	}

	fmt.Printf("Scraping successful: %v\n", response.Success)
	fmt.Printf("HTML length: %d\n", len(response.HTML))
	fmt.Printf("Credits used: %d\n", credits)
}
```

### Extracting Article Content

```go
package main

import (
	"fmt"
	"log"

	"github.com/ujeebu/ujeebu-go-sdk"
)

func main() {
	client, err := ujeebu.NewClient("YOUR-API-KEY")
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}

	params := ujeebu.ExtractParams{
		URL: "https://ujeebu.com/blog/web-scraping-in-2025-state-of-the-art-and-trends/",
		JS:  true,
	}

	article, credits, err := client.Extract(params)
	if err != nil {
		log.Fatalf("Extraction failed: %v", err)
	}

	fmt.Printf("Title: %s\n", article.Title)
	fmt.Printf("Author: %s\n", article.Author)
	fmt.Printf("Publication Date: %s\n", article.PubDate)
	fmt.Printf("Summary: %s\n", article.Summary)
	fmt.Printf("Credits used: %d\n", credits)
}
```

### Checking Account Usage

```go
package main

import (
	"fmt"
	"log"

	"github.com/ujeebu/ujeebu-go-sdk"
)

func main() {
	client, err := ujeebu.NewClient("YOUR-API-KEY")
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}

	accountInfo, err := client.Account()
	if err != nil {
		log.Fatalf("Failed to get account information: %v", err)
	}

	fmt.Printf("Plan: %s\n", accountInfo.Plan)
	fmt.Printf("Credits used: %d\n", accountInfo.Used)
}
```


## Configuration

You can customize the API client:

```go
client, err := ujeebu.NewClient("YOUR-API-KEY")
if err != nil {
	log.Fatalf("Failed to create client: %v", err)
}

// Set a custom timeout
client.SetTimeout(120 * time.Second)

// Change the API key dynamically
client.SetAPIKey("NEW-API-KEY")
```

## Environment Variables

The SDK supports the following environment variables:

- `UJEEBU_BASE_URL`: Override the default API base URL (default: https://api.ujeebu.com)

## Error Handling

All API methods return an error as the last return value, which can be used to handle errors:

```go
response, credits, err := client.Scrape(params)
if err != nil {
	// Handle error
	log.Fatalf("API error: %v", err)
}
```


