package main

import (
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/ujeebu/ujeebu-go-sdk"
)

func main() {
	// Get API key from environment variable
	apiKey := os.Getenv("UJEEBU_API_KEY")
	if apiKey == "" {
		log.Fatal("UJEEBU_API_KEY environment variable is required")
	}

	// Create a new client
	client, err := ujeebu.NewClient(apiKey)
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}

	// Basic scraping
	fmt.Println("=== Basic Scraping ===")
	basicScrape(client)

	// Scraping with JavaScript
	fmt.Println("\n=== Scraping with JavaScript ===")
	scrapeWithJS(client)

	// Scraping with infinite scroll
	fmt.Println("\n=== Scraping with Infinite Scroll ===")
	scrapeWithScroll(client)

	// Scraping with extraction rules
	fmt.Println("\n=== Scraping with Extraction Rules ===")
	scrapeWithRules(client)
}

func basicScrape(client *ujeebu.Client) {
	params := ujeebu.ScrapeParams{
		URL: "https://ujeebu.com",
	}

	response, credits, err := client.Scrape(params)
	if err != nil {
		log.Printf("Scraping failed: %v", err)
		return
	}

	fmt.Printf("Success: %v\n", response.Success)
	fmt.Printf("HTML length: %d bytes\n", len(response.HTML))
	fmt.Printf("Credits used: %d\n", credits)
}

func scrapeWithJS(client *ujeebu.Client) {
	params := ujeebu.ScrapeParams{
		URL:            "https://example.com/",
		JS:             true,
		JSTimeout:      10000,
		Device:         "desktop",
		BlockAds:       true,
		BlockResources: true,
	}

	response, credits, err := client.Scrape(params)
	if err != nil {
		var apiErr *ujeebu.APIError
		if errors.As(err, &apiErr) {
			fmt.Printf("API Error: %s (status: %d)\n", apiErr.Message, apiErr.StatusCode)
		} else {
			log.Printf("Scraping failed: %v", err)
		}
		return
	}

	fmt.Printf("Success: %v\n", response.Success)
	fmt.Printf("HTML length: %d bytes\n", len(response.HTML))
	fmt.Printf("Credits used: %d\n", credits)
}

func scrapeWithScroll(client *ujeebu.Client) {
	params := ujeebu.ScrapeParams{
		URL:               "https://example.com/",
		JS:                true,
		Timeout:           60,
		JSTimeout:         30,
		ScrollDown:        true,
		ScrollWait:        1000,
		ProgressiveScroll: true,
		Device:            "desktop",
		WindowWidth:       1920,
		WindowHeight:      1080,
	}

	response, credits, err := client.Scrape(params)
	if err != nil {
		log.Printf("Scraping failed: %v", err)
		return
	}

	fmt.Printf("Success: %v\n", response.Success)
	fmt.Printf("HTML length: %d bytes\n", len(response.HTML))
	fmt.Printf("Credits used: %d\n", credits)
}

func scrapeWithRules(client *ujeebu.Client) {
	params := ujeebu.ScrapeParams{
		URL: "https://books.toscrape.com/",
		ExtractRules: map[string]any{
			"products": map[string]any{
				"_selector": ".product_pod",
				"title": map[string]string{
					"_selector":  "h3 a",
					"_attribute": "title",
				},
				"price": ".price_color",
				"image": map[string]string{
					"_selector":  ".image_container img",
					"_attribute": "src",
				},
				"rating": map[string]string{
					"_selector":  ".star-rating",
					"_attribute": "class",
				},
			},
		},
	}

	response, credits, err := client.Scrape(params)
	if err != nil {
		log.Printf("Scraping failed: %v", err)
		return
	}

	fmt.Printf("Success: %v\n", response.Success)
	fmt.Printf("Extracted data: %+v\n", response.Result)
	fmt.Printf("Credits used: %d\n", credits)
}
