package repository

import (
	"context"
	"errors"
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

// --- Query types & functions ---

type NearbyCarpark struct {
	CarparkCode   string
	CarparkName   string
	DataSource    string
	Lat           float64
	Lon           float64
	DistanceM     float64
	ParkingSystem *string
	TotalLots     *int
}

type CarparkRow struct {
	CarparkCode   string
	CarparkName   string
	DataSource    string
	CarparkType   *string
	ParkingSystem *string
	Lat           *float64
	Lon           *float64
	TotalLots     *int
}

const nearbySQL = `
SELECT carpark_code, carpark_name, data_source, lat, lon, parking_system, total_lots, distance_m
FROM (
    SELECT
        carpark_code, carpark_name, data_source, lat, lon, parking_system, total_lots,
        (6371000 * acos(LEAST(1.0,
            cos(radians($1)) * cos(radians(lat)) * cos(radians(lon) - radians($2))
            + sin(radians($1)) * sin(radians(lat))
        ))) AS distance_m
    FROM carparks
    WHERE lat IS NOT NULL
      AND lat BETWEEN $1 - ($3 / 111320.0) AND $1 + ($3 / 111320.0)
      AND lon BETWEEN $2 - ($3 / 78710.0)  AND $2 + ($3 / 78710.0)
) sub
WHERE distance_m <= $3
ORDER BY distance_m
LIMIT $4`

func GetNearby(ctx context.Context, pool *pgxpool.Pool, lat, lon, radiusM float64, limit int) ([]NearbyCarpark, error) {
	rows, err := pool.Query(ctx, nearbySQL, lat, lon, radiusM, limit)
	if err != nil {
		return nil, fmt.Errorf("query nearby: %w", err)
	}
	defer rows.Close()

	var results []NearbyCarpark
	for rows.Next() {
		var r NearbyCarpark
		if err := rows.Scan(
			&r.CarparkCode, &r.CarparkName, &r.DataSource,
			&r.Lat, &r.Lon, &r.ParkingSystem, &r.TotalLots, &r.DistanceM,
		); err != nil {
			return nil, fmt.Errorf("scan nearby: %w", err)
		}
		results = append(results, r)
	}
	return results, rows.Err()
}

func GetByCodes(ctx context.Context, pool *pgxpool.Pool, codes []string, lat, lon float64) ([]NearbyCarpark, error) {
	const sql = `
	SELECT carpark_code, carpark_name, data_source, lat, lon, parking_system, total_lots,
	    (6371000 * acos(LEAST(1.0,
	        cos(radians($1)) * cos(radians(lat)) * cos(radians(lon) - radians($2))
	        + sin(radians($1)) * sin(radians(lat))
	    ))) AS distance_m
	FROM carparks
	WHERE carpark_code = ANY($3) AND lat IS NOT NULL
	ORDER BY distance_m`

	rows, err := pool.Query(ctx, sql, lat, lon, codes)
	if err != nil {
		return nil, fmt.Errorf("query batch carparks: %w", err)
	}
	defer rows.Close()

	var results []NearbyCarpark
	for rows.Next() {
		var r NearbyCarpark
		if err := rows.Scan(
			&r.CarparkCode, &r.CarparkName, &r.DataSource,
			&r.Lat, &r.Lon, &r.ParkingSystem, &r.TotalLots, &r.DistanceM,
		); err != nil {
			return nil, fmt.Errorf("scan batch carpark: %w", err)
		}
		results = append(results, r)
	}
	return results, rows.Err()
}

func GetByCode(ctx context.Context, pool *pgxpool.Pool, code, source string) (*CarparkRow, error) {
	const sql = `
		SELECT carpark_code, carpark_name, data_source, carpark_type, parking_system, lat, lon, total_lots
		FROM carparks
		WHERE carpark_code = $1 AND data_source = $2`

	var r CarparkRow
	err := pool.QueryRow(ctx, sql, code, source).Scan(
		&r.CarparkCode, &r.CarparkName, &r.DataSource, &r.CarparkType,
		&r.ParkingSystem, &r.Lat, &r.Lon, &r.TotalLots,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("query carpark %s/%s: %w", source, code, err)
	}
	return &r, nil
}
