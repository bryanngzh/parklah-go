package repository

import (
	"context"

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
