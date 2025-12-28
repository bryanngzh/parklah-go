package ura

import (
	"fmt"
	"log"
	"net/http"
	"time"
)

type URAClient struct {
	BaseURL string
	AccessKey string
	Token string
	LastFetched time.Time
	HTTPClient *http.Client
}

// Initialize New Client
func NewClient(accessKey string) (*URAClient, error) {
	client := &URAClient{
		BaseURL: "https://eservice.ura.gov.sg/uraDataService",
		AccessKey: accessKey,
		HTTPClient: &http.Client{Timeout: 10 * time.Second}, // creates http client and returns the pointer
	}

	if err := client.getToken(); err != nil {
		return nil, fmt.Errorf("[new-client] failed to get initial URA token: %w", err)
	}

	log.Println("[ura] Client initialized successfully ✅")
	return client, nil
}

// URA API Calls
// FetchCarparkAvailability returns the latest available lots;
func (c *URAClient) FetchCarparkAvailability() ([]CarparkAvailabilityResponse, error) {
	payload, err := callURAAPI[[]CarparkAvailabilityResponse](c, "Car_Park_Availability")
	if err != nil {
		return nil, err
	}
	return payload.Result, nil
}

// FetchCarparkDetails returns static carpark details
func (c *URAClient) FetchCarparkDetails() ([]CarparkDetailsResponse, error) {
	payload, err := callURAAPI[[]CarparkDetailsResponse](c, "Car_Park_Details")
	if err != nil {
		return nil, err
	}
	return payload.Result, nil
}

// FetchCarparkSeasonDetails returns season parking details
func (c *URAClient) FetchCarparkSeasonDetails() ([]CarparkSeasonDetailsResponse, error) {
	payload, err := callURAAPI[[]CarparkSeasonDetailsResponse](c, "Season_Car_Park_Details")
	if err != nil {
		return nil, err
	}
	return payload.Result, nil
}