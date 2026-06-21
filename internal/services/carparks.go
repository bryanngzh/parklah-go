package services

import (
	"context"
	"fmt"
	"math"
	"time"

	"github.com/bryanngzh/parklah-go/internal/repository"
	"github.com/jackc/pgx/v5/pgxpool"
)

var sgt = time.FixedZone("SGT", 8*60*60)

type NearbyCarparkResult struct {
	CarparkCode   string               `json:"carpark_code"`
	CarparkName   string               `json:"carpark_name"`
	DataSource    string               `json:"data_source"`
	Lat           float64              `json:"lat"`
	Lon           float64              `json:"lon"`
	DistanceM     float64              `json:"distance_m"`
	ParkingSystem *string              `json:"parking_system"`
	Availability  []AvailabilityDetail `json:"availability"`
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
	IsCurrent    bool    `json:"is_current"`
}

type SeasonRateResult struct {
	VehicleType string  `json:"vehicle_type"`
	TicketType  string  `json:"ticket_type"`
	ParkingHrs  string  `json:"parking_hrs"`
	MonthlyRate float64 `json:"monthly_rate"`
}

func GetNearby(ctx context.Context, pool *pgxpool.Pool, lat, lon, radiusM float64, limit int) ([]NearbyCarparkResult, NearbyMeta, error) {
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

	availMap := make(map[string][]AvailabilityDetail)
	for _, a := range avail {
		key := a.CarparkCode + ":" + a.DataSource
		availMap[key] = append(availMap[key], AvailabilityDetail{
			VehicleType:   a.VehicleType,
			LotsAvailable: a.LotsAvailable,
			TotalLots:     a.TotalLots,
			SnapshotTime:  a.SnapshotTime,
		})
	}

	results := make([]NearbyCarparkResult, len(nearby))
	for i, cp := range nearby {
		key := cp.CarparkCode + ":" + cp.DataSource
		avail := availMap[key]
		if avail == nil {
			avail = []AvailabilityDetail{}
		}
		results[i] = NearbyCarparkResult{
			CarparkCode:   cp.CarparkCode,
			CarparkName:   cp.CarparkName,
			DataSource:    cp.DataSource,
			Lat:           cp.Lat,
			Lon:           cp.Lon,
			DistanceM:     math.Round(cp.DistanceM),
			ParkingSystem: cp.ParkingSystem,
			Availability:  avail,
		}
	}
	return results, NearbyMeta{Count: len(results), RadiusM: radiusM}, nil
}

func GetBatch(ctx context.Context, pool *pgxpool.Pool, lat, lon float64, codes []string) ([]NearbyCarparkResult, error) {
	carparks, err := repository.GetByCodes(ctx, pool, codes, lat, lon)
	if err != nil {
		return nil, err
	}
	if len(carparks) == 0 {
		return []NearbyCarparkResult{}, nil
	}

	carparkCodes := make([]string, len(carparks))
	for i, cp := range carparks {
		carparkCodes[i] = cp.CarparkCode
	}

	avail, err := repository.GetLatestAvailability(ctx, pool, carparkCodes)
	if err != nil {
		return nil, err
	}

	availMap := make(map[string][]AvailabilityDetail)
	for _, a := range avail {
		key := a.CarparkCode + ":" + a.DataSource
		availMap[key] = append(availMap[key], AvailabilityDetail{
			VehicleType:   a.VehicleType,
			LotsAvailable: a.LotsAvailable,
			TotalLots:     a.TotalLots,
			SnapshotTime:  a.SnapshotTime,
		})
	}

	results := make([]NearbyCarparkResult, len(carparks))
	for i, cp := range carparks {
		key := cp.CarparkCode + ":" + cp.DataSource
		a := availMap[key]
		if a == nil {
			a = []AvailabilityDetail{}
		}
		results[i] = NearbyCarparkResult{
			CarparkCode:   cp.CarparkCode,
			CarparkName:   cp.CarparkName,
			DataSource:    cp.DataSource,
			Lat:           cp.Lat,
			Lon:           cp.Lon,
			DistanceM:     math.Round(cp.DistanceM),
			ParkingSystem: cp.ParkingSystem,
			Availability:  a,
		}
	}
	return results, nil
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

func GetRates(ctx context.Context, pool *pgxpool.Pool, code, source string, phDates map[string]bool) (*RatesResult, error) {
	shortTerm, err := repository.GetShortTermRates(ctx, pool, code, source)
	if err != nil {
		return nil, err
	}

	season, err := repository.GetSeasonRates(ctx, pool, code, source)
	if err != nil {
		return nil, err
	}

	now := time.Now().In(sgt)
	dayType := currentDayType(now, phDates)

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
			IsCurrent:    isCurrentRate(r.DayType, r.StartTime, r.EndTime, now, dayType),
		}
	}
	// If a specific day_type row is active for a vehicle type, suppress overlapping "all" rows
	// so that e.g. sunday_ph free parking takes priority over the base "all" rate.
	activeSpecific := make(map[string]bool)
	for _, r := range result.ShortTerm {
		if r.IsCurrent && r.DayType != "all" {
			activeSpecific[r.VehicleType] = true
		}
	}
	for i, r := range result.ShortTerm {
		if r.DayType == "all" && activeSpecific[r.VehicleType] {
			result.ShortTerm[i].IsCurrent = false
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

// currentDayType returns the URA-style day type for the given time,
// accounting for Singapore public holidays (treated as sunday_ph).
func currentDayType(t time.Time, phDates map[string]bool) string {
	if phDates[t.Format("2006-01-02")] {
		return "sunday_ph"
	}
	switch t.Weekday() {
	case time.Sunday:
		return "sunday_ph"
	case time.Saturday:
		return "saturday"
	default:
		return "weekday"
	}
}

// isCurrentRate returns true if the given rate row is active right now.
// Handles overnight ranges (e.g. 22:00–07:00) and day_type="all".
func isCurrentRate(dayType, startTime, endTime string, now time.Time, todayType string) bool {
	if dayType != "all" && dayType != todayType {
		return false
	}
	start, err1 := parseTimeOfDayMins(startTime)
	end, err2 := parseTimeOfDayMins(endTime)
	if err1 != nil || err2 != nil {
		return false
	}
	nowMins := now.Hour()*60 + now.Minute()
	if end <= start {
		// Overnight range: active if now >= start OR now < end
		return nowMins >= start || nowMins < end
	}
	return nowMins >= start && nowMins < end
}

// parseTimeOfDayMins parses "HH:MM" or "HH:MM:SS" into minutes since midnight.
func parseTimeOfDayMins(s string) (int, error) {
	if len(s) < 5 {
		return 0, fmt.Errorf("invalid time: %s", s)
	}
	var h, m int
	if _, err := fmt.Sscanf(s[:5], "%d:%d", &h, &m); err != nil {
		return 0, err
	}
	return h*60 + m, nil
}
