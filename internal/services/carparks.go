package services

import (
	"context"
	"math"
	"time"

	"github.com/bryanngzh/parklah-go/internal/repository"
	"github.com/jackc/pgx/v5/pgxpool"
)

type NearbyCarparkResult struct {
	CarparkCode   string     `json:"carpark_code"`
	CarparkName   string     `json:"carpark_name"`
	DataSource    string     `json:"data_source"`
	Lat           float64    `json:"lat"`
	Lon           float64    `json:"lon"`
	DistanceM     float64    `json:"distance_m"`
	ParkingSystem *string    `json:"parking_system"`
	TotalLots     *int       `json:"total_lots"`
	LotsAvailable *int       `json:"lots_available"`
	SnapshotTime  *time.Time `json:"snapshot_time"`
}

type NearbyMeta struct {
	Count   int     `json:"count"`
	RadiusM float64 `json:"radius_m"`
}

type CarparkDetail struct {
	CarparkCode   string               `json:"carpark_code"`
	CarparkName   string               `json:"carpark_name"`
	DataSource    string               `json:"data_source"`
	CarparkType   *string              `json:"carpark_type"`
	ParkingSystem *string              `json:"parking_system"`
	Lat           *float64             `json:"lat"`
	Lon           *float64             `json:"lon"`
	TotalLots     *int                 `json:"total_lots"`
	Features      *FeaturesDetail      `json:"features"`
	Availability  []AvailabilityDetail `json:"availability"`
}

type FeaturesDetail struct {
	ShortTermParking  string  `json:"short_term_parking"`
	FreeParking       string  `json:"free_parking"`
	NightParking      bool    `json:"night_parking"`
	CarParkDecks      int     `json:"car_park_decks"`
	GantryHeight      float64 `json:"gantry_height"`
	CarParkBasement   bool    `json:"car_park_basement"`
	IsCentralArea     bool    `json:"is_central_area"`
	IsPeakHourCarpark bool    `json:"is_peak_hour_carpark"`
}

type AvailabilityDetail struct {
	VehicleType   string    `json:"vehicle_type"`
	LotsAvailable int       `json:"lots_available"`
	TotalLots     *int      `json:"total_lots"`
	SnapshotTime  time.Time `json:"snapshot_time"`
}

type RatesResult struct {
	ShortTerm []ShortTermRateResult `json:"short_term"`
	Season    []SeasonRateResult    `json:"season"`
}

type ShortTermRateResult struct {
	VehicleType  string  `json:"vehicle_type"`
	DayType      string  `json:"day_type"`
	StartTime    string  `json:"start_time"`
	EndTime      string  `json:"end_time"`
	RatePer30Min float64 `json:"rate_per_30min"`
	MinDuration  string  `json:"min_duration,omitempty"`
}

type SeasonRateResult struct {
	VehicleType string  `json:"vehicle_type"`
	TicketType  string  `json:"ticket_type"`
	ParkingHrs  string  `json:"parking_hrs"`
	MonthlyRate float64 `json:"monthly_rate"`
}

func GetNearby(ctx context.Context, pool *pgxpool.Pool, lat, lon, radiusM float64, vehicleType string, limit int) ([]NearbyCarparkResult, NearbyMeta, error) {
	nearby, err := repository.GetNearby(ctx, pool, lat, lon, radiusM, limit)
	if err != nil {
		return nil, NearbyMeta{}, err
	}
	if len(nearby) == 0 {
		return []NearbyCarparkResult{}, NearbyMeta{Count: 0, RadiusM: radiusM}, nil
	}

	codes := make([]string, len(nearby))
	for i, cp := range nearby {
		codes[i] = cp.CarparkCode
	}

	avail, err := repository.GetLatestAvailability(ctx, pool, codes)
	if err != nil {
		return nil, NearbyMeta{}, err
	}

	availMap := make(map[string]repository.AvailabilityRow)
	for _, a := range avail {
		if a.VehicleType == vehicleType {
			key := a.CarparkCode + ":" + a.DataSource
			availMap[key] = a
		}
	}

	results := make([]NearbyCarparkResult, len(nearby))
	for i, cp := range nearby {
		r := NearbyCarparkResult{
			CarparkCode:   cp.CarparkCode,
			CarparkName:   cp.CarparkName,
			DataSource:    cp.DataSource,
			Lat:           cp.Lat,
			Lon:           cp.Lon,
			DistanceM:     math.Round(cp.DistanceM),
			ParkingSystem: cp.ParkingSystem,
			TotalLots:     cp.TotalLots,
		}
		if a, ok := availMap[cp.CarparkCode+":"+cp.DataSource]; ok {
			r.LotsAvailable = &a.LotsAvailable
			r.SnapshotTime = &a.SnapshotTime
		}
		results[i] = r
	}
	return results, NearbyMeta{Count: len(results), RadiusM: radiusM}, nil
}

func GetCarparkDetail(ctx context.Context, pool *pgxpool.Pool, code, source string) (*CarparkDetail, error) {
	cp, err := repository.GetByCode(ctx, pool, code, source)
	if err != nil {
		return nil, err
	}
	if cp == nil {
		return nil, nil
	}

	avail, err := repository.GetLatestAvailability(ctx, pool, []string{code})
	if err != nil {
		return nil, err
	}

	feat, err := repository.GetFeatures(ctx, pool, code, source)
	if err != nil {
		return nil, err
	}

	detail := &CarparkDetail{
		CarparkCode:   cp.CarparkCode,
		CarparkName:   cp.CarparkName,
		DataSource:    cp.DataSource,
		CarparkType:   cp.CarparkType,
		ParkingSystem: cp.ParkingSystem,
		Lat:           cp.Lat,
		Lon:           cp.Lon,
		TotalLots:     cp.TotalLots,
		Availability:  []AvailabilityDetail{},
	}

	if feat != nil {
		detail.Features = &FeaturesDetail{
			ShortTermParking:  feat.ShortTermParking,
			FreeParking:       feat.FreeParking,
			NightParking:      feat.NightParking,
			CarParkDecks:      feat.CarParkDecks,
			GantryHeight:      feat.GantryHeight,
			CarParkBasement:   feat.CarParkBasement,
			IsCentralArea:     feat.IsCentralArea,
			IsPeakHourCarpark: feat.IsPeakHourCarpark,
		}
	}

	for _, a := range avail {
		if a.DataSource == source {
			detail.Availability = append(detail.Availability, AvailabilityDetail{
				VehicleType:   a.VehicleType,
				LotsAvailable: a.LotsAvailable,
				TotalLots:     a.TotalLots,
				SnapshotTime:  a.SnapshotTime,
			})
		}
	}
	return detail, nil
}

func GetAvailability(ctx context.Context, pool *pgxpool.Pool, code, source string) ([]AvailabilityDetail, error) {
	avail, err := repository.GetLatestAvailability(ctx, pool, []string{code})
	if err != nil {
		return nil, err
	}

	results := []AvailabilityDetail{}
	for _, a := range avail {
		if a.DataSource == source {
			results = append(results, AvailabilityDetail{
				VehicleType:   a.VehicleType,
				LotsAvailable: a.LotsAvailable,
				TotalLots:     a.TotalLots,
				SnapshotTime:  a.SnapshotTime,
			})
		}
	}
	return results, nil
}

func GetRates(ctx context.Context, pool *pgxpool.Pool, code, source string) (*RatesResult, error) {
	shortTerm, err := repository.GetShortTermRates(ctx, pool, code, source)
	if err != nil {
		return nil, err
	}

	season, err := repository.GetSeasonRates(ctx, pool, code, source)
	if err != nil {
		return nil, err
	}

	result := &RatesResult{
		ShortTerm: make([]ShortTermRateResult, len(shortTerm)),
		Season:    make([]SeasonRateResult, len(season)),
	}
	for i, r := range shortTerm {
		result.ShortTerm[i] = ShortTermRateResult{
			VehicleType:  r.VehicleType,
			DayType:      r.DayType,
			StartTime:    r.StartTime,
			EndTime:      r.EndTime,
			RatePer30Min: r.RatePer30Min,
			MinDuration:  r.MinDuration,
		}
	}
	for i, r := range season {
		result.Season[i] = SeasonRateResult{
			VehicleType: r.VehicleType,
			TicketType:  r.TicketType,
			ParkingHrs:  r.ParkingHrs,
			MonthlyRate: r.MonthlyRate,
		}
	}
	return result, nil
}
