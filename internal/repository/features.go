package repository

import (
	"context"
	"fmt"

	"github.com/bryanngzh/parklah-go/internal/models"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

const upsertFeaturesSQL = `
INSERT INTO carpark_features
    (carpark_code, data_source, short_term_parking, free_parking, night_parking,
     car_park_decks, gantry_height, car_park_basement, is_central_area, is_peak_hour_carpark)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
ON CONFLICT (carpark_code, data_source) DO UPDATE SET
    short_term_parking   = EXCLUDED.short_term_parking,
    free_parking         = EXCLUDED.free_parking,
    night_parking        = EXCLUDED.night_parking,
    car_park_decks       = EXCLUDED.car_park_decks,
    gantry_height        = EXCLUDED.gantry_height,
    car_park_basement    = EXCLUDED.car_park_basement,
    is_central_area      = EXCLUDED.is_central_area,
    is_peak_hour_carpark = EXCLUDED.is_peak_hour_carpark,
    updated_at           = now()`

func UpsertFeaturesBatch(ctx context.Context, pool *pgxpool.Pool, features []models.Features) error {
	if len(features) == 0 {
		return nil
	}
	batch := &pgx.Batch{}
	for _, f := range features {
		batch.Queue(upsertFeaturesSQL,
			f.CarparkCode, f.DataSource, f.ShortTermParking, f.FreeParking,
			f.NightParking, f.CarParkDecks, f.GantryHeight, f.CarParkBasement,
			f.IsCentralArea, f.IsPeakHourCarpark,
		)
	}
	br := pool.SendBatch(ctx, batch)
	for i := range features {
		if _, err := br.Exec(); err != nil {
			br.Close()
			return fmt.Errorf("upsert features %d (%s): %w", i, features[i].CarparkCode, err)
		}
	}
	return br.Close()
}
