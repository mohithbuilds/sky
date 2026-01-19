package openmateo

import (
	"fmt"
	"log"
	"net/http"
	"time"
)

type baseClient struct {
	httpClient *http.Client
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

// AIR-QUALITY CLIENT
const airQualityBaseURL = "https://air-quality-api.open-meteo.com/v1/"
