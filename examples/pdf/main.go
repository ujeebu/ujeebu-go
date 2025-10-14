package main

import (
	"encoding/base64"
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

	// Basic PDF generation
	fmt.Println("=== Basic PDF Generation ===")
	basicPDF(client)

	// PDF with custom viewport
	fmt.Println("\n=== PDF with Custom Viewport ===")
	customPDF(client)
}

func basicPDF(client *ujeebu.Client) {
	params := ujeebu.ScrapeParams{
		URL: "https://ujeebu.com/blog/web-scraping-in-2025-state-of-the-art-and-trends/",
		JS:  true,
	}

	pdf, credits, err := client.PDF(params)
	if err != nil {
		log.Printf("PDF generation failed: %v", err)
		return
	}

	// Save PDF to file
	if err := savePDF(pdf, "article.pdf"); err != nil {
		log.Printf("Failed to save PDF: %v", err)
		return
	}

	fmt.Println("PDF saved to article.pdf")
	fmt.Printf("Credits used: %d\n", credits)
}

func customPDF(client *ujeebu.Client) {
	params := ujeebu.ScrapeParams{
		URL:          "https://ujeebu.com",
		JS:           true,
		Device:       "desktop",
		WindowWidth:  1920,
		WindowHeight: 1080,
	}

	pdf, credits, err := client.PDF(params)
	if err != nil {
		log.Printf("PDF generation failed: %v", err)
		return
	}

	// Save PDF to file
	if err := savePDF(pdf, "custom.pdf"); err != nil {
		log.Printf("Failed to save PDF: %v", err)
		return
	}

	fmt.Println("Custom PDF saved to custom.pdf")
	fmt.Printf("Credits used: %d\n", credits)
}

func savePDF(base64Data, filename string) error {
	// Decode base64 string
	data, err := base64.StdEncoding.DecodeString(base64Data)
	if err != nil {
		return fmt.Errorf("failed to decode base64: %w", err)
	}

	// Write to file
	if err := os.WriteFile(filename, data, 0644); err != nil {
		return fmt.Errorf("failed to write file: %w", err)
	}

	return nil
}
