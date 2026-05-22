package hdb

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
)

type HDBClient struct {
	APIKey     string
	HTTPClient *http.Client
	BaseURL    string
}

// Initialize new HDB client
func NewClient(apiKey string) *HDBClient {
	return &HDBClient{
		APIKey:     apiKey,
		BaseURL:    "https://api.data.gov.sg",
		HTTPClient: &http.Client{Timeout: 10 * time.Second},
	}
}

// Fetch static carpark information from data.gov.sg
func (c *HDBClient) FetchCarparkInfo() ([]CarparkInfoResponse, error) {
	datasetID := "d_23f946fa557947f93a8043bbef41dd09"
	url := fmt.Sprintf("https://data.gov.sg/api/action/datastore_search?resource_id=%s&limit=5000", datasetID)

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("[hdb-carpark-info] creating request: %w", err)
	}

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("[hdb-carpark-info] performing request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("[hdb-carpark-info] non-200 status: %s", resp.Status)
	}

	var apiResp CarparkInfoAPIResponse
	if err := json.NewDecoder(resp.Body).Decode(&apiResp); err != nil {
		return nil, fmt.Errorf("[hdb-carpark-info] decoding response: %w", err)
	}

	if !apiResp.Success {
		return nil, fmt.Errorf("[hdb-carpark-info] API returned failure")
	}

	log.Printf("[hdb] Fetched %d carpark records\n", len(apiResp.Result.Records))
	return apiResp.Result.Records, nil
}

// Fetch real-time carpark availability from data.gov.sg
func (c *HDBClient) FetchCarparkAvailability() (CarparkAvailabilityResponse, error) {
	url := fmt.Sprintf("%s/v1/transport/carpark-availability", c.BaseURL)

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return CarparkAvailabilityResponse{}, fmt.Errorf("[hdb-availability] creating request: %w", err)
	}
	req.Header.Set("X-Api-Key", c.APIKey)

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return CarparkAvailabilityResponse{}, fmt.Errorf("[hdb-availability] performing request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return CarparkAvailabilityResponse{}, fmt.Errorf("[hdb-availability] non-200 status: %s", resp.Status)
	}

	var availability CarparkAvailabilityResponse
	if err := json.NewDecoder(resp.Body).Decode(&availability); err != nil {
		return CarparkAvailabilityResponse{}, fmt.Errorf("[hdb-availability] decoding response: %w", err)
	}

	totalCarparks := len(availability.Items[0].CarparkData)
	
	log.Printf("[hdb] Fetched availability for %d carparks\n", totalCarparks)
	return availability, nil
}
