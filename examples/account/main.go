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

	// Get account information
	fmt.Println("=== Account Information ===")
	getAccountInfo(client)
}

func getAccountInfo(client *ujeebu.Client) {
	account, err := client.Account()
	if err != nil {
		log.Fatalf("Failed to get account info: %v", err)
	}

	// Display account information
	fmt.Printf("User ID: %s\n", account.UserID)
	fmt.Printf("Plan: %s\n", account.Plan)
	fmt.Printf("Quota: %s credits\n", account.Quota)
	fmt.Println()

	// Usage information
	fmt.Println("Usage:")
	fmt.Printf("  Used: %d credits (%.1f%%)\n", account.Used, account.UsedPercent)
	fmt.Printf("  Remaining: %d credits\n", account.Balance)
	fmt.Printf("  Total Requests: %d\n", account.TotalRequests)
	fmt.Printf("  Concurrent Requests Allowed: %d\n", account.ConcurrentRequests)
	fmt.Println()

	// Billing information
	if account.NextBillingDate != nil {
		fmt.Println("Billing:")
		fmt.Printf("  Next Billing Date: %s\n", *account.NextBillingDate)
		fmt.Printf("  Days Until Next Billing: %d\n", account.DaysTillNextBilling)
	} else {
		fmt.Println("Billing: No billing information available")
	}

	// Usage bar visualization
	fmt.Println()
	printUsageBar(account.UsedPercent)
}

func printUsageBar(percent float64) {
	const barWidth = 50
	filled := int(percent / 100.0 * float64(barWidth))
	if filled > barWidth {
		filled = barWidth
	}

	fmt.Print("Usage: [")
	for i := 0; i < barWidth; i++ {
		if i < filled {
			fmt.Print("=")
		} else {
			fmt.Print(" ")
		}
	}
	fmt.Printf("] %.1f%%\n", percent)
}
