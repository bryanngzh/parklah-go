package ura

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
)

func callURAAPI[T any](c *URAClient, endpoint string) (URAResponse[T], error) {
	if err := c.ensureValidToken(); err != nil {
		return URAResponse[T]{}, err
	}

	url := fmt.Sprintf("%s/invokeUraDS/v1?service=%s", c.BaseURL, endpoint)
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return URAResponse[T]{}, fmt.Errorf("[%s] creating request: %w", endpoint, err)
	}
	req.Header.Set("AccessKey", c.AccessKey)
	req.Header.Set("Token", c.Token)

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return URAResponse[T]{}, fmt.Errorf("[%s] performing request: %w", endpoint, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return URAResponse[T]{}, fmt.Errorf("[%s] non-200 status: %s", endpoint, resp.Status)
	}

	var payload URAResponse[T]
	if err := json.NewDecoder(resp.Body).Decode(&payload); err != nil {
		return URAResponse[T]{}, fmt.Errorf("[%s] decoding response: %w", endpoint, err)
	}

	if payload.Status != "Success" {
		return URAResponse[T]{}, fmt.Errorf("[%s] URA returned failure: %s", endpoint, payload.Message)
	}
	return payload, nil
}

// Gets new token - refreshes every 24 hours
func (c *URAClient) getToken() error {
	url := fmt.Sprintf("%s/insertNewToken/v1", c.BaseURL)

	// Set up request
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return fmt.Errorf("[get-token] creating request: %w", err)
	}
	req.Header.Set("AccessKey", c.AccessKey)

	// Get token
	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return fmt.Errorf("[get-token] requesting new token: %w", err)
	}
	// close network connection right before function returns
	defer resp.Body.Close()
	
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("[get-token] failed to get token, status: %s", resp.Status)
	}

	var payload URAResponse[string]
	if err := json.NewDecoder(resp.Body).Decode(&payload); err != nil {
		return fmt.Errorf("[get-token] decoding token response: %w", err)
	}

	if payload.Status != "Success" {
		return fmt.Errorf("[get-token] URA returned failure: %s", payload.Message)
	}

	c.Token = payload.Result
	c.LastFetched = time.Now()

	log.Println("[ura] New URA token fetched successfully ✅")
	return nil
}

// Ensures each API call has a valid token
func (c *URAClient) ensureValidToken() error {
	if c.Token == "" || time.Since(c.LastFetched).Hours() >= 24 {
		log.Println("[ura] Refreshing URA token (expired or missing)")
		return c.getToken()
	}
	return nil
}