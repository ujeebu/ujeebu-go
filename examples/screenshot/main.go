package main

import (
	"encoding/base64"
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

	// Full page screenshot
	fmt.Println("=== Full Page Screenshot ===")
	fullPageScreenshot(client)

	// Partial screenshot with selector
	fmt.Println("\n=== Partial Screenshot (CSS Selector) ===")
	partialScreenshot(client)

	// Screenshot with custom viewport
	fmt.Println("\n=== Screenshot with Custom Viewport ===")
	customViewportScreenshot(client)
}

func fullPageScreenshot(client *ujeebu.Client) {
	params := ujeebu.ScrapeParams{
		URL: "https://ujeebu.com",
		JS:  true,
	}

	screenshot, credits, err := client.Screenshot(params, true, "")
	if err != nil {
		log.Printf("Screenshot failed: %v", err)
		return
	}

	// Save screenshot to file
	if err := saveScreenshot(screenshot, "fullpage.png"); err != nil {
		log.Printf("Failed to save screenshot: %v", err)
		return
	}

	fmt.Println("Full page screenshot saved to fullpage.png")
	fmt.Printf("Credits used: %d\n", credits)
}

func partialScreenshot(client *ujeebu.Client) {
	params := ujeebu.ScrapeParams{
		URL: "https://ujeebu.com",
		JS:  true,
	}

	// Capture only the main content area
	screenshot, credits, err := client.Screenshot(params, false, "main")
	if err != nil {
		log.Printf("Screenshot failed: %v", err)
		return
	}

	// Save screenshot to file
	if err := saveScreenshot(screenshot, "partial.png"); err != nil {
		log.Printf("Failed to save screenshot: %v", err)
		return
	}

	fmt.Println("Partial screenshot saved to partial.png")
	fmt.Printf("Credits used: %d\n", credits)
}

func customViewportScreenshot(client *ujeebu.Client) {
	params := ujeebu.ScrapeParams{
		URL:          "https://ujeebu.com",
		JS:           true,
		Device:       "mobile",
		WindowWidth:  375,
		WindowHeight: 667,
	}

	screenshot, credits, err := client.Screenshot(params, true, "")
	if err != nil {
		log.Printf("Screenshot failed: %v", err)
		return
	}

	// Save screenshot to file
	if err := saveScreenshot(screenshot, "mobile.png"); err != nil {
		log.Printf("Failed to save screenshot: %v", err)
		return
	}

	fmt.Println("Mobile viewport screenshot saved to mobile.png")
	fmt.Printf("Credits used: %d\n", credits)
}

func saveScreenshot(base64Data, filename string) error {
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
