package main

import (
	"context"
	"log"

	"github.com/bryanngzh/parklah-go/internal/config"
	"github.com/bryanngzh/parklah-go/internal/db"
	"github.com/bryanngzh/parklah-go/internal/hdb"
	"github.com/bryanngzh/parklah-go/internal/models"
	"github.com/bryanngzh/parklah-go/internal/repository"
	"github.com/bryanngzh/parklah-go/internal/ura"
)

func main() {
	ctx := context.Background()

	cfg := config.Load()

	pool, err := db.Connect(ctx, cfg.DSN())
	if err != nil {
		log.Fatalf("[main] DB connect: %v", err)
	}
	defer pool.Close()
	log.Println("[main] Connected to database")

	uraClient, err := ura.NewClient(cfg.URAAccessKey)
	if err != nil {
		log.Fatalf("[main] URA client: %v", err)
	}
	hdbClient := hdb.NewClient(cfg.DataGovAPIKey)

	// --- URA static ---
	staticData, err := fetchStaticCarparkData(uraClient)
	if err != nil {
		log.Fatalf("[main] fetch URA static: %v", err)
	}

	uraCarparks, uraShortRates := ura.TransformDetails(staticData.Details)
	log.Printf("[ura] details → %d carparks, %d short-term rates", len(uraCarparks), len(uraShortRates))

	uraSeasonCarparks, uraSeasonRates := ura.TransformSeasonDetails(staticData.SeasonDetails)
	log.Printf("[ura] season → %d carparks, %d season rates", len(uraSeasonCarparks), len(uraSeasonRates))

	if err := repository.UpsertCarparks(ctx, pool, uraCarparks); err != nil {
		log.Fatalf("[main] upsert URA carparks (details): %v", err)
	}
	if err := repository.UpsertCarparks(ctx, pool, uraSeasonCarparks); err != nil {
		log.Fatalf("[main] upsert URA carparks (season): %v", err)
	}
	if err := repository.UpsertShortTermRates(ctx, pool, uraShortRates); err != nil {
		log.Fatalf("[main] upsert URA short-term rates: %v", err)
	}
	if err := repository.UpsertSeasonRates(ctx, pool, uraSeasonRates); err != nil {
		log.Fatalf("[main] upsert URA season rates: %v", err)
	}
	log.Println("[ura] static data persisted")

	// --- HDB static ---
	hdbInfo, err := fetchHDBStaticData(hdbClient)
	if err != nil {
		log.Fatalf("[main] fetch HDB static: %v", err)
	}

	hdbCarparks, hdbFeatures := hdb.TransformCarparkInfo(hdbInfo)
	log.Printf("[hdb] info → %d carparks, %d feature rows", len(hdbCarparks), len(hdbFeatures))

	if err := repository.UpsertCarparks(ctx, pool, hdbCarparks); err != nil {
		log.Fatalf("[main] upsert HDB carparks: %v", err)
	}
	if err := repository.UpsertFeaturesBatch(ctx, pool, hdbFeatures); err != nil {
		log.Fatalf("[main] upsert HDB features: %v", err)
	}

	hdbRates := deriveHDBRates(hdbCarparks)
	log.Printf("[hdb] derived %d short-term rate rows", len(hdbRates))
	if err := repository.UpsertShortTermRates(ctx, pool, hdbRates); err != nil {
		log.Fatalf("[main] upsert HDB short-term rates: %v", err)
	}
	log.Println("[hdb] static data persisted")

	// --- URA availability ---
	uraAvail, err := fetchCarparkAvailability(uraClient)
	if err != nil {
		log.Fatalf("[main] fetch URA availability: %v", err)
	}
	uraAvailRows := ura.TransformAvailability(uraAvail)
	log.Printf("[ura] %d availability rows", len(uraAvailRows))
	if err := repository.InsertAvailabilityBatch(ctx, pool, uraAvailRows); err != nil {
		log.Fatalf("[main] insert URA availability: %v", err)
	}

	// --- HDB availability ---
	hdbAvailResp, err := fetchHDBAvailability(hdbClient)
	if err != nil {
		log.Fatalf("[main] fetch HDB availability: %v", err)
	}
	hdbAvailRows := hdb.TransformHDBAvailability(hdbAvailResp)
	log.Printf("[hdb] %d availability rows", len(hdbAvailRows))
	if err := repository.InsertAvailabilityBatch(ctx, pool, hdbAvailRows); err != nil {
		log.Fatalf("[main] insert HDB availability: %v", err)
	}

	log.Println("[main] Ingestion complete")
}

func deriveHDBRates(carparks []models.Carpark) []models.ShortTermRate {
	var rates []models.ShortTermRate
	for _, cp := range carparks {
		rates = append(rates, hdb.DeriveShortTermRates(cp.CarparkCode)...)
	}
	return rates
}
