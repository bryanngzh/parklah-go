package ura

import (
	"strconv"
	"strings"
	"time"

	"github.com/bryanngzh/parklah-go/internal/models"
	"github.com/bryanngzh/parklah-go/internal/util"
)

func TransformDetails(rows []CarparkDetailsResponse) ([]models.Carpark, []models.ShortTermRate) {
	carparkMap := make(map[string]models.Carpark)
	var rates []models.ShortTermRate

	for _, r := range rows {
		if _, ok := carparkMap[r.PpCode]; !ok {
			lat, lon := parseCoordinates(firstCoord(r.Geometries))
			carparkMap[r.PpCode] = models.Carpark{
				CarparkCode:   r.PpCode,
				CarparkName:   r.PpName,
				DataSource:    "ura",
				ParkingSystem: normalizeParkingSystem(r.ParkingSystem),
				Lat:           lat,
				Lon:           lon,
				TotalLots:     ptrInt(r.ParkCapacity),
			}
		}

		vt := normalizeVehicleType(r.VehCat)
		start := parseURATime(r.StartTime)
		end := parseURATime(r.EndTime)
		if start == "" {
			start = "00:00"
		}
		if end == "" {
			end = "23:59"
		}

		wRate := parseURARate(r.WeekdayRate)
		sRate := parseURARate(r.SatdayRate)
		phRate := parseURARate(r.SunPHRate)

		if wRate == sRate && sRate == phRate && r.WeekdayMin == r.SatdayMin && r.SatdayMin == r.SunPHMin {
			rates = append(rates, models.ShortTermRate{
				CarparkCode:  r.PpCode,
				DataSource:   "ura",
				VehicleType:  vt,
				DayType:      "all",
				StartTime:    start,
				EndTime:      end,
				RatePer30Min: wRate,
				MinDuration:  r.WeekdayMin,
			})
		} else {
			rates = append(rates,
				models.ShortTermRate{CarparkCode: r.PpCode, DataSource: "ura", VehicleType: vt, DayType: "weekday", StartTime: start, EndTime: end, RatePer30Min: wRate, MinDuration: r.WeekdayMin},
				models.ShortTermRate{CarparkCode: r.PpCode, DataSource: "ura", VehicleType: vt, DayType: "saturday", StartTime: start, EndTime: end, RatePer30Min: sRate, MinDuration: r.SatdayMin},
				models.ShortTermRate{CarparkCode: r.PpCode, DataSource: "ura", VehicleType: vt, DayType: "sunday_ph", StartTime: start, EndTime: end, RatePer30Min: phRate, MinDuration: r.SunPHMin},
			)
		}
	}

	carparks := make([]models.Carpark, 0, len(carparkMap))
	for _, cp := range carparkMap {
		carparks = append(carparks, cp)
	}
	return carparks, rates
}

func TransformSeasonDetails(rows []CarparkSeasonDetailsResponse) ([]models.Carpark, []models.SeasonRate) {
	carparkMap := make(map[string]models.Carpark)
	var rates []models.SeasonRate

	for _, r := range rows {
		if _, ok := carparkMap[r.PpCode]; !ok {
			lat, lon := parseCoordinates(firstCoord(r.Geometries))
			carparkMap[r.PpCode] = models.Carpark{
				CarparkCode: r.PpCode,
				CarparkName: r.PpName,
				DataSource:  "ura",
				Lat:         lat,
				Lon:         lon,
			}
		}
		rates = append(rates, models.SeasonRate{
			CarparkCode: r.PpCode,
			DataSource:  "ura",
			VehicleType: normalizeVehicleType(r.VehCat),
			TicketType:  r.TicketType,
			ParkingHrs:  r.ParkingHrs,
			MonthlyRate: parseURARate(r.MonthlyRate),
		})
	}

	carparks := make([]models.Carpark, 0, len(carparkMap))
	for _, cp := range carparkMap {
		carparks = append(carparks, cp)
	}
	return carparks, rates
}

func TransformAvailability(rows []CarparkAvailabilityResponse) []models.Availability {
	now := time.Now().UTC()
	avail := make([]models.Availability, 0, len(rows))
	for _, r := range rows {
		if r.LotType != "C" && r.LotType != "M" && r.LotType != "H" {
			continue
		}
		lots, err := strconv.Atoi(r.LotsAvailable)
		if err != nil {
			continue
		}
		avail = append(avail, models.Availability{
			CarparkCode:   r.CarparkNo,
			DataSource:    "ura",
			VehicleType:   r.LotType,
			LotsAvailable: lots,
			TotalLots:     nil,
			SnapshotTime:  now,
		})
	}
	return avail
}

func firstCoord(geometries []struct {
	Coordinates string `json:"coordinates"`
}) string {
	if len(geometries) == 0 {
		return ""
	}
	return geometries[0].Coordinates
}

func parseCoordinates(s string) (*float64, *float64) {
	if s == "" {
		return nil, nil
	}
	parts := strings.SplitN(s, ",", 2)
	if len(parts) != 2 {
		return nil, nil
	}
	x, err1 := strconv.ParseFloat(strings.TrimSpace(parts[0]), 64)
	y, err2 := strconv.ParseFloat(strings.TrimSpace(parts[1]), 64)
	if err1 != nil || err2 != nil {
		return nil, nil
	}
	lat, lon := util.SVY21ToWGS84(x, y)
	return &lat, &lon
}

func parseURATime(s string) string {
	if s == "" {
		return ""
	}
	t, err := time.Parse("03.04 PM", s)
	if err != nil {
		return ""
	}
	return t.Format("15:04")
}

func parseURARate(s string) float64 {
	s = strings.TrimSpace(strings.TrimPrefix(strings.TrimSpace(s), "$"))
	v, _ := strconv.ParseFloat(s, 64)
	return v
}

func normalizeParkingSystem(s string) *string {
	var v string
	switch s {
	case "B", "P":
		v = "electronic"
	case "C":
		v = "coupon"
	default:
		return nil
	}
	return &v
}

func normalizeVehicleType(s string) string {
	switch s {
	case "Motorcycle":
		return "M"
	case "Heavy Vehicle":
		return "H"
	default:
		return "C"
	}
}

func ptrInt(v int) *int {
	if v == 0 {
		return nil
	}
	return &v
}
