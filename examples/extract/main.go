package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/ujeebu/ujeebu-go"
)

func main() {
	// Get API key from environment variable
	apiKey := os.Getenv("UJEEBU_API_KEY")
	if apiKey == "" {
		log.Fatal("UJEEBU_API_KEY environment variable is required")
	}

	// Create a new client with custom options
	client, err := ujeebu.NewClient(
		apiKey,
		ujeebu.WithTimeout(120*time.Second),
	)
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}

	// Basic extraction
	fmt.Println("=== Basic Article Extraction ===")
	basicExtraction(client)

	// Extraction with JavaScript
	fmt.Println("\n=== Extraction with JavaScript Rendering ===")
	extractionWithJS(client)

	// Extraction with context and timeout
	fmt.Println("\n=== Extraction with Context Timeout ===")
	extractionWithContext(client)

	// Extract from raw HTML
	fmt.Println("\n=== Extract from Raw HTML ===")
	extractFromRawHTML(client)
}

func basicExtraction(client *ujeebu.Client) {
	params := ujeebu.ExtractParams{
		URL: "https://ujeebu.com/blog/web-scraping-in-2025-state-of-the-art-and-trends/",
	}

	article, credits, err := client.Extract(params)
	if err != nil {
		log.Printf("Extraction failed: %v", err)
		return
	}

	fmt.Printf("Title: %s\n", article.Title)
	fmt.Printf("Author: %s\n", article.Author)
	fmt.Printf("Publication Date: %s\n", article.PubDate)
	fmt.Printf("Summary: %s\n", truncate(article.Summary, 200))
	fmt.Printf("Text length: %d characters\n", len(article.Text))
	fmt.Printf("Images: %d\n", len(article.Images))
	fmt.Printf("Credits used: %d\n", credits)
}

func extractionWithJS(client *ujeebu.Client) {
	params := ujeebu.ExtractParams{
		URL:       "https://ujeebu.com/blog/web-scraping-in-2025-state-of-the-art-and-trends/",
		JS:        true,
		Timeout:   120,
		JSTimeout: 60,
		Text:      true,
		Images:    true,
		Author:    true,
		PubDate:   true,
	}

	article, credits, err := client.Extract(params)
	if err != nil {
		log.Printf("Extraction failed: %v", err)
		return
	}

	fmt.Printf("Title: %s\n", article.Title)
	fmt.Printf("Text length: %d characters\n", len(article.Text))
	fmt.Printf("Credits used: %d\n", credits)
}

func extractionWithContext(client *ujeebu.Client) {
	// Create a context with 30-second timeout
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	params := ujeebu.ExtractParams{
		URL: "https://ujeebu.com/blog/web-scraping-in-2025-state-of-the-art-and-trends/",
	}

	article, credits, err := client.ExtractWithContext(ctx, params)
	if err != nil {
		if ctx.Err() == context.DeadlineExceeded {
			log.Println("Request timed out")
		} else {
			log.Printf("Extraction failed: %v", err)
		}
		return
	}

	fmt.Printf("Title: %s\n", article.Title)
	fmt.Printf("Credits used: %d\n", credits)
}

func extractFromRawHTML(client *ujeebu.Client) {
	rawHTML := `
		<html>
			<head>
				<title>Example Article</title>
				<meta name="author" content="John Doe">
			</head>
			<body>
				<article>
					<h1>Example Article Title</h1>
					<p>This is the main content of the article.</p>
					<p>It contains multiple paragraphs.</p>
				</article>
			</body>
		</html>
	`

	params := ujeebu.ExtractParams{
		URL:     "https://example.com/",
		RawHTML: rawHTML,
	}

	article, credits, err := client.Extract(params)
	if err != nil {
		log.Printf("Extraction failed: %v", err)
		return
	}

	fmt.Printf("Title: %s\n", article.Title)
	fmt.Printf("Author: %s\n", article.Author)
	fmt.Printf("Text: %s\n", article.Text)
	fmt.Printf("Credits used: %d\n", credits)
}

func truncate(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen] + "..."
}
