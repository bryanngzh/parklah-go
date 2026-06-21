package util

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

// FetchSGPublicHolidays fetches Singapore public holidays for a given year
// from the Nager.Date public API (no auth required).
// Returns a map of "YYYY-MM-DD" → true.
func FetchSGPublicHolidays(ctx context.Context, year int) (map[string]bool, error) {
	url := fmt.Sprintf("https://date.nager.at/api/v3/PublicHolidays/%d/SG", year)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var holidays []struct {
		Date string `json:"date"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&holidays); err != nil {
		return nil, err
	}

	result := make(map[string]bool, len(holidays))
	for _, h := range holidays {
		result[h.Date] = true
	}
	return result, nil
}

func IsSGPublicHoliday(t time.Time, phDates map[string]bool) bool {
	return phDates[t.Format("2006-01-02")]
}
