package main

import (
	"fmt"
	"log"

	"github.com/bryanngzh/parklah-go/internal/config"
	"github.com/bryanngzh/parklah-go/internal/ura"
)

func main() {
	// Load environment configs
	cfg := config.Load()

	// Connect to DB

	
	// Init URA Client and Repo
	client, err := ura.NewClient(cfg.URAAccessKey)
	if (err != nil) {
		panic(err)
	}

	// Run ingestion job
	availability, err := client.FetchCarparkAvailability()
	if err != nil {
		log.Fatalf("Failed to fetch carpark availability: %v", err)
	}
	fmt.Printf("\n[Carpark Availability] Fetched %d records\n", len(availability))
	if len(availability) > 0 {
		fmt.Printf("Sample: %+v\n", availability[0])
	}

	details, err := client.FetchCarparkDetails()
	if err != nil {
		log.Fatalf("Failed to fetch carpark details: %v", err)
	}
	fmt.Printf("\n[Carpark Details] Fetched %d records\n", len(details))
	if len(details) > 0 {
		fmt.Printf("Sample: %+v\n", details[0])
	}

	seasonDetails, err := client.FetchCarparkSeasonDetails()
	if err != nil {
		log.Fatalf("Failed to fetch carpark season details: %v", err)
	}
	fmt.Printf("\n[Carpark Season Details] Fetched %d records\n", len(seasonDetails))
	if len(seasonDetails) > 0 {
		fmt.Printf("Sample: %+v\n", seasonDetails[0])
	}
}