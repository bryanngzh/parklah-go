package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/bryanngzh/parklah-go/internal/models"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

func InsertAvailabilityBatch(ctx context.Context, pool *pgxpool.Pool, avail []models.Availability) error {
	if len(avail) == 0 {
		return nil
	}
	rows := make([][]any, len(avail))
	for i, a := range avail {
		rows[i] = []any{a.CarparkCode, a.DataSource, a.VehicleType, a.LotsAvailable, a.TotalLots, a.SnapshotTime}
	}
	_, err := pool.CopyFrom(
		ctx,
		pgx.Identifier{"carpark_availability"},
		[]string{"carpark_code", "data_source", "vehicle_type", "lots_available", "total_lots", "snapshot_time"},
		pgx.CopyFromRows(rows),
	)
	return err
}

// --- Query types & functions ---

type AvailabilityRow struct {
	CarparkCode   string
	DataSource    string
	VehicleType   string
	LotsAvailable int
	TotalLots     *int
	SnapshotTime  time.Time
}

func GetLatestAvailability(ctx context.Context, pool *pgxpool.Pool, codes []string) ([]AvailabilityRow, error) {
	if len(codes) == 0 {
		return nil, nil
	}
	const sql = `
		SELECT DISTINCT ON (carpark_code, data_source, vehicle_type)
		    carpark_code, data_source, vehicle_type, lots_available, total_lots, snapshot_time
		FROM carpark_availability
		WHERE carpark_code = ANY($1)
		ORDER BY carpark_code, data_source, vehicle_type, snapshot_time DESC`

	rows, err := pool.Query(ctx, sql, codes)
	if err != nil {
		return nil, fmt.Errorf("query availability: %w", err)
	}
	defer rows.Close()

	var results []AvailabilityRow
	for rows.Next() {
		var r AvailabilityRow
		if err := rows.Scan(
			&r.CarparkCode, &r.DataSource, &r.VehicleType,
			&r.LotsAvailable, &r.TotalLots, &r.SnapshotTime,
		); err != nil {
			return nil, fmt.Errorf("scan availability: %w", err)
		}
		results = append(results, r)
	}
	return results, rows.Err()
}
