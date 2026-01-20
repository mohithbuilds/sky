package openmateo

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

type APIError struct {
	Error  bool   `json:"error"`
	Reason string `json:"reason"`
}

type baseClient struct {
	httpClient *http.Client
}

func (bc *baseClient) doRequest(url string) ([]byte, error) {
	resp, err := bc.httpClient.Get(url)
	if err != nil {
		return nil, fmt.Errorf("Failed to GET URL: %s: %w", url, err)
	}

	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("Unable to read the response body: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		var apiErr APIError

		if json.Unmarshal(data, &apiErr) == nil && apiErr.Error {
			return nil, fmt.Errorf("API error (%s): %s", resp.Status, apiErr.Reason)
		}

		return nil, fmt.Errorf(
			"API returned non-OK status: %s, body: %s",
			resp.Status,
			string(data),
		)
	}

	return data, nil
}

// GEOCODING CLIENT
const geocodingBaseURL = "https://geocoding-api.open-meteo.com/v1/"

type GeocodingClient struct {
	*baseClient
	BaseURL string
}

func NewGeocodingClient(httpClient *http.Client) *GeocodingClient {
	if httpClient == nil {
		httpClient = &http.Client{Timeout: 3 * time.Second}
	}
	return &GeocodingClient{
		baseClient: &baseClient{
			httpClient: httpClient,
		},
		BaseURL: geocodingBaseURL,
	}
}

// FORECAST CLIENT
const forecastBaseURL = "https://api.open-meteo.com/v1/"

type ForecastClient struct {
	*baseClient
	BaseURL string
}

func NewForecastClient(httpClient *http.Client) *ForecastClient {
	if httpClient == nil {
		httpClient = &http.Client{Timeout: 3 * time.Second}
	}
	return &ForecastClient{
		baseClient: &baseClient{
			httpClient: httpClient,
		},
		BaseURL: forecastBaseURL,
	}
}

// AIR-QUALITY CLIENT
const airQualityBaseURL = "https://air-quality-api.open-meteo.com/v1/"

type AirQualityClient struct {
	*baseClient
	BaseURL string
}

func NewAirQualityClient(httpClient *http.Client) *AirQualityClient {
	if httpClient == nil {
		httpClient = &http.Client{Timeout: 3 * time.Second}
	}
	return &AirQualityClient{
		baseClient: &baseClient{
			httpClient: httpClient,
		},
		BaseURL: airQualityBaseURL,
	}
}
