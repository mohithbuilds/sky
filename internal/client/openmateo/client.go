package openmateo

import (
	"fmt"
	"io"
	"net/http"
	"time"
)

type baseClient struct {
	httpClient *http.Client
}

func (bc *baseClient) doRequest(url string) ([]byte, error) {
	resp, err := bc.httpClient.Get(url)
	if err != nil {
		return nil, fmt.Errorf("Failed to GET URL: %s: %w", url, err)
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API returned non-OK status: %s", resp.Status)
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("Unable to read the response body: %w", err)
	}

	return data, nil
}

// GEOCODING CLIENT
const geocodingBaseURL = "https://geocoding-api.open-meteo.com/v1/"

type GeocodingClient struct {
	*baseClient
	BaseURL string
}

func NewGeocodingClient() *GeocodingClient {
	return &GeocodingClient{
		baseClient: &baseClient{
			httpClient: &http.Client{Timeout: 3 * time.Second},
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

func NewForecastClient() *ForecastClient {
	return &ForecastClient{
		baseClient: &baseClient{
			httpClient: &http.Client{Timeout: 3 * time.Second},
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

func NewAirQualityClient() *AirQualityClient {
	return &AirQualityClient{
		baseClient: &baseClient{
			httpClient: &http.Client{Timeout: 3 * time.Second},
		},
		BaseURL: airQualityBaseURL,
	}
}
