package hdb

import (
	"strconv"
	"strings"
	"time"

	"github.com/bryanngzh/parklah-go/internal/models"
	"github.com/bryanngzh/parklah-go/internal/util"
)

func TransformCarparkInfo(rows []CarparkInfoResponse) ([]models.Carpark, []models.Features) {
	carparks := make([]models.Carpark, 0, len(rows))
	features := make([]models.Features, 0, len(rows))

	for _, r := range rows {
		lat, lon := parseHDBCoordinates(r.XCoord, r.YCoord)
		ct := normalizeCarParkType(r.CarParkType)
		ps := normalizeHDBParkingSystem(r.TypeOfParkingSystem)

		carparks = append(carparks, models.Carpark{
			CarparkCode:   r.CarParkNo,
			CarparkName:   r.Address,
			DataSource:    "hdb",
			CarparkType:   ct,
			ParkingSystem: ps,
			Lat:           lat,
			Lon:           lon,
		})

		decks, _ := strconv.Atoi(strings.TrimSpace(r.CarParkDecks))
		gantry, _ := strconv.ParseFloat(strings.TrimSpace(r.GantryHeight), 64)

		features = append(features, models.Features{
			CarparkCode:       r.CarParkNo,
			DataSource:        "hdb",
			ShortTermParking:  r.ShortTermParking,
			FreeParking:       r.FreeParking,
			NightParking:      strings.EqualFold(r.NightParking, "YES"),
			CarParkDecks:      decks,
			GantryHeight:      gantry,
			CarParkBasement:   strings.EqualFold(r.CarParkBasement, "YES"),
			IsCentralArea:     IsCentralArea(r.CarParkNo),
			IsPeakHourCarpark: IsPeakHour(r.CarParkNo),
		})
	}
	return carparks, features
}

func TransformHDBAvailability(resp CarparkAvailabilityResponse) []models.Availability {
	if len(resp.Items) == 0 {
		return nil
	}
	item := resp.Items[0]

	snapshotTime, err := time.Parse(time.RFC3339, item.Timestamp)
	if err != nil {
		snapshotTime = time.Now().UTC()
	}

	var avail []models.Availability
	for _, cp := range item.CarparkData {
		for _, info := range cp.CarparkInfo {
			vt := info.LotType
			if vt != "C" && vt != "M" && vt != "H" {
				continue
			}
			lots, err := strconv.Atoi(info.LotsAvailable)
			if err != nil {
				continue
			}
			total, err := strconv.Atoi(info.TotalLots)
			var totalPtr *int
			if err == nil && total > 0 {
				totalPtr = &total
			}
			avail = append(avail, models.Availability{
				CarparkCode:   cp.CarparkNumber,
				DataSource:    "hdb",
				VehicleType:   vt,
				LotsAvailable: lots,
				TotalLots:     totalPtr,
				SnapshotTime:  snapshotTime,
			})
		}
	}
	return avail
}

func parseHDBCoordinates(xStr, yStr string) (*float64, *float64) {
	x, err1 := strconv.ParseFloat(strings.TrimSpace(xStr), 64)
	y, err2 := strconv.ParseFloat(strings.TrimSpace(yStr), 64)
	if err1 != nil || err2 != nil || x == 0 || y == 0 {
		return nil, nil
	}
	lat, lon := util.SVY21ToWGS84(x, y)
	return &lat, &lon
}

func normalizeHDBParkingSystem(s string) *string {
	var v string
	switch strings.ToUpper(strings.TrimSpace(s)) {
	case "ELECTRONIC PARKING":
		v = "electronic"
	case "COUPON PARKING":
		v = "coupon"
	default:
		return nil
	}
	return &v
}

func normalizeCarParkType(s string) *string {
	s = strings.TrimSpace(s)
	if s == "" {
		return nil
	}
	upper := strings.ToUpper(s)
	return &upper
}
