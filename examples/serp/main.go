package main

import (
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

	// Google Search
	fmt.Println("=== Google Web Search ===")
	googleSearch(client)

	// Google News
	fmt.Println("\n=== Google News Search ===")
	googleNews(client)

	// Google Videos
	fmt.Println("\n=== Google Video Search ===")
	googleVideos(client)

	// Google Images
	fmt.Println("\n=== Google Image Search ===")
	googleImages(client)

	// Google Maps
	fmt.Println("\n=== Google Maps Search ===")
	googleMaps(client)
}

func googleSearch(client *ujeebu.Client) {
	params := ujeebu.SerpParams{
		Search:       "golang web scraping",
		Lang:         "en",
		Location:     "us",
		Device:       "desktop",
		ResultsCount: 10,
		Page:         1,
	}

	results, credits, err := client.GoogleSearch(params)
	if err != nil {
		log.Printf("Search failed: %v", err)
		return
	}

	fmt.Printf("Query: %s\n", results.Metadata.QueryDisplayed)
	fmt.Printf("Total results: %d\n", results.Metadata.NumberOfResults)
	fmt.Printf("Results time: %s\n", results.Metadata.ResultsTime)

	// Knowledge Graph
	if results.KnowledgeGraph.Title != "" {
		fmt.Printf("\nKnowledge Graph:\n")
		fmt.Printf("  Title: %s\n", results.KnowledgeGraph.Title)
		fmt.Printf("  Type: %s\n", results.KnowledgeGraph.Type)
	}

	// Organic Results
	fmt.Printf("\nOrganic Results:\n")
	for i, result := range results.OrganicResults {
		if i >= 5 {
			break
		}
		fmt.Printf("%d. %s\n", result.Position, result.Title)
		fmt.Printf("   %s\n", result.Link)
		fmt.Printf("   %s\n", truncate(result.Description, 100))
	}

	fmt.Printf("\nCredits used: %d\n", credits)
}

func googleNews(client *ujeebu.Client) {
	params := ujeebu.SerpParams{
		Search:   "technology news",
		Location: "us",
		Lang:     "en",
	}

	results, credits, err := client.GoogleNewsSearch(params)
	if err != nil {
		log.Printf("News search failed: %v", err)
		return
	}

	fmt.Printf("News Results:\n")
	for i, news := range results.News {
		if i >= 5 {
			break
		}
		fmt.Printf("%d. %s\n", news.Position, news.Title)
		fmt.Printf("   Source: %s\n", news.SiteName)
		fmt.Printf("   Date: %s\n", news.Date)
		fmt.Printf("   Link: %s\n", news.Link)
		fmt.Println()
	}

	fmt.Printf("Credits used: %d\n", credits)
}

func googleVideos(client *ujeebu.Client) {
	params := ujeebu.SerpParams{
		Search: "golang tutorial",
		Lang:   "en",
	}

	results, credits, err := client.GoogleVideoSearch(params)
	if err != nil {
		log.Printf("Video search failed: %v", err)
		return
	}

	fmt.Printf("Video Results:\n")
	for i, video := range results.Videos {
		if i >= 5 {
			break
		}
		fmt.Printf("%d. %s\n", video.Position, video.Title)
		fmt.Printf("   Provider: %s\n", video.Provider)
		fmt.Printf("   Author: %s\n", video.Author)
		fmt.Printf("   URL: %s\n", video.Url)
		fmt.Println()
	}

	fmt.Printf("Credits used: %d\n", credits)
}

func googleImages(client *ujeebu.Client) {
	params := ujeebu.SerpParams{
		Search: "golang logo",
		Lang:   "en",
	}

	results, credits, err := client.GoogleImageSearch(params)
	if err != nil {
		log.Printf("Image search failed: %v", err)
		return
	}

	fmt.Printf("Image Results:\n")
	for i, image := range results.Images {
		if i >= 5 {
			break
		}
		fmt.Printf("%d. %s\n", image.Position, image.Title)
		fmt.Printf("   Source: %s\n", image.Source)
		fmt.Printf("   Size: %dx%d\n", image.Width, image.Height)
		fmt.Printf("   Link: %s\n", image.Link)
		fmt.Println()
	}

	fmt.Printf("Credits used: %d\n", credits)
}

func googleMaps(client *ujeebu.Client) {
	params := ujeebu.SerpParams{
		Search:   "restaurants near me",
		Location: "us",
		Lang:     "en",
	}

	results, credits, err := client.GoogleMapSearch(params)
	if err != nil {
		log.Printf("Maps search failed: %v", err)
		return
	}

	fmt.Printf("Map Results:\n")
	for i, place := range results.Maps {
		if i >= 5 {
			break
		}
		fmt.Printf("%d. %s\n", place.Position, place.Title)
		fmt.Printf("   Address: %s\n", place.Address)
		fmt.Printf("   Category: %s\n", place.Category)
		fmt.Printf("   Rating: %.1f (%d reviews)\n", place.Rating, place.Reviews)
		fmt.Println()
	}

	fmt.Printf("Credits used: %d\n", credits)
}

func truncate(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen] + "..."
}
