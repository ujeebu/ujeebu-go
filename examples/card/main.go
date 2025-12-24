package main

import (
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/ujeebu/ujeebu-go"
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

	// Basic card/preview fetch
	fmt.Println("=== Basic Card Preview ===")
	basicCard(client)

	// Card fetch with JavaScript (for SPAs)
	fmt.Println("\n=== Card with JavaScript Rendering ===")
	cardWithJS(client)

	// Card fetch with error handling
	fmt.Println("\n=== Card with Error Handling ===")
	cardWithErrorHandling(client)
}

func basicCard(client *ujeebu.Client) {
	params := ujeebu.CardParams{
		URL: "https://ujeebu.com/blog/web-scraping-in-2025-state-of-the-art-and-trends/",
	}

	card, credits, err := client.Card(params)
	if err != nil {
		log.Printf("Card fetch failed: %v", err)
		return
	}

	fmt.Printf("URL: %s\n", card.URL)
	fmt.Printf("Title: %s\n", card.Title)
	fmt.Printf("Summary: %s\n", truncate(card.Summary, 150))
	fmt.Printf("Author: %s\n", card.Author)
	fmt.Printf("Site Name: %s\n", card.SiteName)
	fmt.Printf("Image: %s\n", card.Image)
	fmt.Printf("Favicon: %s\n", card.Favicon)
	fmt.Printf("Published: %s\n", card.DatePublished)
	fmt.Printf("Modified: %s\n", card.DateModified)
	fmt.Printf("Language: %s\n", card.Lang)
	fmt.Printf("Keywords: %v\n", card.Keywords)
	fmt.Printf("Credits used: %d\n", credits)
}

func cardWithJS(client *ujeebu.Client) {
	params := ujeebu.CardParams{
		URL:       "https://ujeebu.com/blog/web-scraping-in-2025-state-of-the-art-and-trends/",
		JS:        true,
		JSTimeout: 5000,
	}

	card, credits, err := client.Card(params)
	if err != nil {
		log.Printf("Card fetch failed: %v", err)
		return
	}

	fmt.Printf("Title: %s\n", card.Title)
	fmt.Printf("Summary: %s\n", truncate(card.Summary, 100))
	fmt.Printf("Credits used: %d\n", credits)
}

func cardWithErrorHandling(client *ujeebu.Client) {
	params := ujeebu.CardParams{
		URL:     "https://example.com/does-not-exist",
		Timeout: 20,
	}

	card, credits, err := client.Card(params)
	if err != nil {
		// Type-specific error handling using errors.As
		var apiErr *ujeebu.APIError
		var validationErr *ujeebu.ValidationError
		var netErr *ujeebu.NetworkError

		switch {
		case errors.As(err, &apiErr):
			fmt.Printf("API Error: %s (status: %d)\n", apiErr.Message, apiErr.StatusCode)
			if apiErr.IsNotFound() {
				fmt.Println("URL not found or unable to fetch")
			}
			if apiErr.IsUnauthorized() {
				fmt.Println("Invalid API key")
			}

		case errors.As(err, &validationErr):
			fmt.Printf("Validation error for %s: %s\n", validationErr.Field, validationErr.Message)

		case errors.As(err, &netErr):
			fmt.Printf("Network error: %v\n", netErr.Err)

		default:
			fmt.Printf("Unknown error: %v\n", err)
		}
		return
	}

	fmt.Printf("Title: %s\n", card.Title)
	fmt.Printf("Credits used: %d\n", credits)
}

func truncate(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen] + "..."
}
