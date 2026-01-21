package openmateo

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Only extracting the wanted parts of the air quality API response
type AirQualityResult struct {
	GenerationTimeMs float64                `json:"generationtime_ms"`
	Hourly           *AirQualityHourly      `json:"hourly,omitempty"`
	HourlyUnits      *AirQualityHourlyUnits `json:"hourly_units,omitempty"`
}

// Hourly contains the time-series data for each hourly air quality parameter.
// Note: Pollen data fields are only available for European locations.
type AirQualityHourly struct {
	Time                []string  `json:"time"`
	PM10                []float64 `json:"pm10"`
	PM25                []float64 `json:"pm2_5"`
	CarbonMonoxide      []float64 `json:"carbon_monoxide"`
	NitrogenDioxide     []float64 `json:"nitrogen_dioxide"`
	SulphurDioxide      []float64 `json:"sulphur_dioxide"`
	Ozone               []float64 `json:"ozone"`
	AerosolOpticalDepth []float64 `json:"aerosol_optical_depth"`
	Dust                []float64 `json:"dust"`
	UVIndex             []float64 `json:"uv_index"`
	AlderPollen         []float64 `json:"alder_pollen"`
	BirchPollen         []float64 `json:"birch_pollen"`
	GrassPollen         []float64 `json:"grass_pollen"`
	MugwortPollen       []float64 `json:"mugwort_pollen"`
	OlivePollen         []float64 `json:"olive_pollen"`
	RagweedPollen       []float64 `json:"ragweed_pollen"`
}

// HourlyUnits contains the units for each hourly air quality parameter
type AirQualityHourlyUnits struct {
	Time                string `json:"time"`
	PM10                string `json:"pm10"`
	PM25                string `json:"pm2_5"`
	CarbonMonoxide      string `json:"carbon_monoxide"`
	NitrogenDioxide     string `json:"nitrogen_dioxide"`
	SulphurDioxide      string `json:"sulphur_dioxide"`
	Ozone               string `json:"ozone"`
	AerosolOpticalDepth string `json:"aerosol_optical_depth"`
	Dust                string `json:"dust"`
	UVIndex             string `json:"uv_index"`
	AlderPollen         string `json:"alder_pollen"`
	BirchPollen         string `json:"birch_pollen"`
	GrassPollen         string `json:"grass_pollen"`
	MugwortPollen       string `json:"mugwort_pollen"`
	OlivePollen         string `json:"olive_pollen"`
	RagweedPollen       string `json:"ragweed_pollen"`
}

// GetAirQuality fetches air quality data for a given latitude and longitude.
// It takes latitude, longitude, and a slice of hourly air quality parameters as input.
// If the hourlyAirQualityParameters slice is empty, it defaults to fetching PM10 and PM2.5 data.
// It returns an AirQualityResult pointer or an error if the request fails or the data cannot be unmarshaled.
func (aqc *AirQualityClient) GetAirQuality(
	latitude, longitude float64,
	hourlyAirQualityParameters []string,
) (*AirQualityResult, error) {
	if len(hourlyAirQualityParameters) == 0 {
		hourlyAirQualityParameters = []string{"pm10", "pm2_5"}
	}
	airQualityURL := fmt.Sprintf(
		"%sair-quality?latitude=%f&longitude=%f&hourly=%s",
		aqc.BaseURL,
		latitude,
		longitude,
		strings.Join(hourlyAirQualityParameters, ","),
	)

	data, err := aqc.doRequest(airQualityURL)
	if err != nil {
		return nil, fmt.Errorf("failed to get air quality data: %w", err)
	}

	var result AirQualityResult
	if err := json.Unmarshal(data, &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal air quality data: %w", err)
	}
	return &result, nil
}
