package repository

import (
	"context"
	"fmt"

	"github.com/bryanngzh/parklah-go/internal/models"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

const upsertShortTermRateSQL = `
INSERT INTO carpark_short_term_rates
    (carpark_code, data_source, vehicle_type, day_type, start_time, end_time, rate_per_30min, min_duration)
VALUES ($1, $2, $3, $4, NULLIF($5, '')::time, NULLIF($6, '')::time, $7, $8)
ON CONFLICT (carpark_code, data_source, vehicle_type, day_type, start_time) DO UPDATE SET
    end_time      = EXCLUDED.end_time,
    rate_per_30min = EXCLUDED.rate_per_30min,
    min_duration  = EXCLUDED.min_duration,
    updated_at    = now()`

const upsertSeasonRateSQL = `
INSERT INTO carpark_season_rates
    (carpark_code, data_source, vehicle_type, ticket_type, parking_hrs, monthly_rate)
VALUES ($1, $2, $3, $4, $5, $6)
ON CONFLICT (carpark_code, data_source, vehicle_type, ticket_type) DO UPDATE SET
    parking_hrs  = EXCLUDED.parking_hrs,
    monthly_rate = EXCLUDED.monthly_rate,
    updated_at   = now()`

func UpsertShortTermRates(ctx context.Context, pool *pgxpool.Pool, rates []models.ShortTermRate) error {
	if len(rates) == 0 {
		return nil
	}
	batch := &pgx.Batch{}
	for _, r := range rates {
		batch.Queue(upsertShortTermRateSQL,
			r.CarparkCode, r.DataSource, r.VehicleType, r.DayType,
			r.StartTime, r.EndTime, r.RatePer30Min, r.MinDuration,
		)
	}
	br := pool.SendBatch(ctx, batch)
	for i := range rates {
		if _, err := br.Exec(); err != nil {
			br.Close()
			return fmt.Errorf("upsert short-term rate %d (%s): %w", i, rates[i].CarparkCode, err)
		}
	}
	return br.Close()
}

func UpsertSeasonRates(ctx context.Context, pool *pgxpool.Pool, rates []models.SeasonRate) error {
	if len(rates) == 0 {
		return nil
	}
	batch := &pgx.Batch{}
	for _, r := range rates {
		batch.Queue(upsertSeasonRateSQL,
			r.CarparkCode, r.DataSource, r.VehicleType, r.TicketType,
			r.ParkingHrs, r.MonthlyRate,
		)
	}
	br := pool.SendBatch(ctx, batch)
	for i := range rates {
		if _, err := br.Exec(); err != nil {
			br.Close()
			return fmt.Errorf("upsert season rate %d (%s): %w", i, rates[i].CarparkCode, err)
		}
	}
	return br.Close()
}
