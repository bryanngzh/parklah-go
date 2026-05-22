package hdb

import "github.com/bryanngzh/parklah-go/internal/models"

var centralAreaCarparks = map[string]bool{
	"ACB": true, "BBB": true, "BRB1": true, "CY": true, "DUXM": true,
	"HLM": true, "KAB": true, "KAM": true, "KAS": true, "PRM": true,
	"SLS": true, "SR1": true, "SR2": true, "TPM": true, "UCS": true, "WCB": true,
}

var peakHourCarparks = map[string]bool{
	"MP1": true, "MP2": true, "BRMS1": true, "BRMS2": true, "SE8": true,
	"SB2": true, "SB3": true, "SB4": true, "PM10": true, "HG14": true,
	"HG8": true, "HG5": true,
}

func IsCentralArea(code string) bool { return centralAreaCarparks[code] }
func IsPeakHour(code string) bool    { return peakHourCarparks[code] }

// DeriveShortTermRates generates the standardised HDB short-term rate rows for a carpark.
// Rates are sourced from the HDB website and derived from the carpark's classification.
func DeriveShortTermRates(carparkCode string) []models.ShortTermRate {
	isCentral := IsCentralArea(carparkCode)
	isPeak := IsPeakHour(carparkCode)

	var rates []models.ShortTermRate

	switch {
	case isCentral:
		rates = append(rates,
			hdbRate(carparkCode, "C", "weekday", "07:00", "17:00", 1.20),
			hdbRate(carparkCode, "C", "all", "00:00", "23:59", 0.60),
		)
	case isPeak:
		rates = append(rates,
			hdbRate(carparkCode, "C", "weekday", "07:30", "09:30", 0.65),
			hdbRate(carparkCode, "C", "weekday", "17:00", "20:00", 0.65),
			hdbRate(carparkCode, "C", "all", "00:00", "23:59", 0.60),
		)
	default:
		rates = append(rates, hdbRate(carparkCode, "C", "all", "00:00", "23:59", 0.60))
	}

	rates = append(rates,
		hdbRate(carparkCode, "M", "all", "00:00", "23:59", 0.65),
		hdbRate(carparkCode, "H", "all", "00:00", "23:59", 1.20),
	)

	return rates
}

func hdbRate(code, vt, dayType, start, end string, amount float64) models.ShortTermRate {
	return models.ShortTermRate{
		CarparkCode:  code,
		DataSource:   "hdb",
		VehicleType:  vt,
		DayType:      dayType,
		StartTime:    start,
		EndTime:      end,
		RatePer30Min: amount,
		MinDuration:  "",
	}
}
