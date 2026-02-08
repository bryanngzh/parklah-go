package main

import (
	"fmt"
	"sync"

	"github.com/bryanngzh/parklah-go/internal/ura"
)

// staticCarparkData groups static carpark data (details and season details) fetched daily.
type staticCarparkData struct {
    Details       []ura.CarparkDetailsResponse
    SeasonDetails []ura.CarparkSeasonDetailsResponse
}

// fetchStaticCarparkData fetches details and season details in parallel (daily).
func fetchStaticCarparkData(client *ura.URAClient) (staticCarparkData, error) {
    var (
        wg       sync.WaitGroup
        mu       sync.Mutex
        firstErr error
        result   staticCarparkData
    )

    setErr := func(err error) {
        if err == nil {
            return
        }
        mu.Lock()
        if firstErr == nil {
            firstErr = err
        }
        mu.Unlock()
    }

    run := func(label string, fn func() error) {
        defer wg.Done()
        if err := fn(); err != nil {
            setErr(fmt.Errorf("%s: %w", label, err))
        }
    }

    wg.Add(2)

    go run("details", func() error {
        res, err := client.FetchCarparkDetails()
        if err != nil {
            return err
        }
        mu.Lock()
        result.Details = res
        mu.Unlock()
        return nil
    })

    go run("season-details", func() error {
        res, err := client.FetchCarparkSeasonDetails()
        if err != nil {
            return err
        }
        mu.Lock()
        result.SeasonDetails = res
        mu.Unlock()
        return nil
    })

    wg.Wait()

    if firstErr != nil {
        return staticCarparkData{}, firstErr
    }

    return result, nil
}

// fetchCarparkAvailability fetches availability data (every 5 mins).
func fetchCarparkAvailability(client *ura.URAClient) ([]ura.CarparkAvailabilityResponse, error) {
    return client.FetchCarparkAvailability()
}

// logStaticCarparkDataSummary prints counts and a sample record for details and season details.
func logStaticCarparkDataSummary(data staticCarparkData) {
    fmt.Printf("\n[Carpark Details] Fetched %d records\n", len(data.Details))
    printSample(data.Details)

    fmt.Printf("\n[Carpark Season Details] Fetched %d records\n", len(data.SeasonDetails))
    printSample(data.SeasonDetails)
}

// logAvailabilitySummary prints count and a sample for availability data.
func logAvailabilitySummary(availability []ura.CarparkAvailabilityResponse) {
    fmt.Printf("\n[Carpark Availability] Fetched %d records\n", len(availability))
    printSample(availability)
}

func printSample[T any](items []T) {
    if len(items) == 0 {
        return
    }
    fmt.Printf("Sample: %+v\n", items[0])
}
