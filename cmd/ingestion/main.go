package main

import (
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

	// Run daily ingestion (details and season details)
	staticData, err := fetchStaticCarparkData(client)
	if err != nil {
		log.Fatalf("Failed to fetch static carpark data: %v", err)
	}
	logStaticCarparkDataSummary(staticData)

	// Run periodic ingestion (availability every 5 mins)
	availability, err := fetchCarparkAvailability(client)
	if err != nil {
		log.Fatalf("Failed to fetch carpark availability: %v", err)
	}
	logAvailabilitySummary(availability)
}