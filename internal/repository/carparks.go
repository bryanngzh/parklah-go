package repository

import (
	"context"
	"fmt"

	"github.com/bryanngzh/parklah-go/internal/models"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

const upsertCarparkSQL = `
INSERT INTO carparks (carpark_code, data_source, carpark_name, carpark_type, parking_system, lat, lon, total_lots)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
ON CONFLICT (carpark_code, data_source) DO UPDATE SET
    carpark_name   = EXCLUDED.carpark_name,
    carpark_type   = COALESCE(EXCLUDED.carpark_type,   carparks.carpark_type),
    parking_system = COALESCE(EXCLUDED.parking_system, carparks.parking_system),
    lat            = COALESCE(EXCLUDED.lat,            carparks.lat),
    lon            = COALESCE(EXCLUDED.lon,            carparks.lon),
    total_lots     = COALESCE(EXCLUDED.total_lots,     carparks.total_lots),
    updated_at     = now()`

func UpsertCarparks(ctx context.Context, pool *pgxpool.Pool, carparks []models.Carpark) error {
	if len(carparks) == 0 {
		return nil
	}
	batch := &pgx.Batch{}
	for _, cp := range carparks {
		batch.Queue(upsertCarparkSQL,
			cp.CarparkCode, cp.DataSource, cp.CarparkName,
			cp.CarparkType, cp.ParkingSystem,
			cp.Lat, cp.Lon, cp.TotalLots,
		)
	}
	br := pool.SendBatch(ctx, batch)
	for i := range carparks {
		if _, err := br.Exec(); err != nil {
			br.Close()
			return fmt.Errorf("upsert carpark %d (%s): %w", i, carparks[i].CarparkCode, err)
		}
	}
	return br.Close()
}
