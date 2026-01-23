package weather

import (
	"errors"
	"sky/internal/client/openmateo"
	"testing"
)

func TestGetCurrentWeather_NilCurrent(t *testing.T) {
	mockClient := &mockForecastClient{
		GetWeatherFunc: func(latitude, longitude float64, currentParameters, hourlyParameters, dailyParameters []string, temperatureUnit, windSpeedUnit, precipitationUnit string, pastDays, forecastDays int64) (*openmateo.ForecastResult, error) {
			return &openmateo.ForecastResult{
				Current: nil,
			}, nil
		},
	}

	weatherClient := NewWeatherClient(mockClient, "celsius", "kmh", "mm")
	_, err := weatherClient.GetCurrentWeather(52.52, 13.41)
	if err == nil {
		t.Fatal("Expected an error for nil current weather data, but got nil")
	}
}

func TestGetCurrentWeather_TimeParseError(t *testing.T) {
	mockClient := &mockForecastClient{
		GetWeatherFunc: func(latitude, longitude float64, currentParameters, hourlyParameters, dailyParameters []string, temperatureUnit, windSpeedUnit, precipitationUnit string, pastDays, forecastDays int64) (*openmateo.ForecastResult, error) {
			return &openmateo.ForecastResult{
				Current: &openmateo.ForecastCurrent{
					Time: "invalid-time",
				},
			}, nil
		},
	}

	weatherClient := NewWeatherClient(mockClient, "celsius", "kmh", "mm")
	_, err := weatherClient.GetCurrentWeather(52.52, 13.41)
	if err == nil {
		t.Fatal("Expected an error for invalid time format, but got nil")
	}
}

func TestGetDailyForecast_NilDaily(t *testing.T) {
	mockClient := &mockForecastClient{
		GetWeatherFunc: func(latitude, longitude float64, currentParameters, hourlyParameters, dailyParameters []string, temperatureUnit, windSpeedUnit, precipitationUnit string, pastDays, forecastDays int64) (*openmateo.ForecastResult, error) {
			return &openmateo.ForecastResult{
				Daily: nil,
			}, nil
		},
	}

	weatherClient := NewWeatherClient(mockClient, "celsius", "kmh", "mm")
	_, err := weatherClient.GetDailyForecast(52.52, 13.41, 1)
	if err == nil {
		t.Fatal("Expected an error for nil daily weather data, but got nil")
	}
}

func TestGetDailyForecast_TimeParseError(t *testing.T) {
	mockClient := &mockForecastClient{
		GetWeatherFunc: func(latitude, longitude float64, currentParameters, hourlyParameters, dailyParameters []string, temperatureUnit, windSpeedUnit, precipitationUnit string, pastDays, forecastDays int64) (*openmateo.ForecastResult, error) {
			return &openmateo.ForecastResult{
				Daily: &openmateo.ForecastDaily{
					Time:             []string{"invalid-time"},
					Temperature2mMax: []float64{12.0},
					Temperature2mMin: []float64{2.0},
					Sunrise:          []string{"2023-01-01T07:00:00Z"},
					Sunset:           []string{"2023-01-01T17:00:00Z"},
					PrecipitationSum: []float64{0.1},
					PrecipitationProbabilityMean: []float64{10.0},
					WeatherCode:      []int{3},
					WindSpeed10mMax:    []float64{15.0},
				},
			}, nil
		},
	}

	weatherClient := NewWeatherClient(mockClient, "celsius", "kmh", "mm")
	_, err := weatherClient.GetDailyForecast(52.52, 13.41, 1)
	if err == nil {
		t.Fatal("Expected an error for invalid daily time format, but got nil")
	}
}

func TestGetDailyForecast_IncompleteData(t *testing.T) {
	mockClient := &mockForecastClient{
		GetWeatherFunc: func(latitude, longitude float64, currentParameters, hourlyParameters, dailyParameters []string, temperatureUnit, windSpeedUnit, precipitationUnit string, pastDays, forecastDays int64) (*openmateo.ForecastResult, error) {
			return &openmateo.ForecastResult{
				Daily: &openmateo.ForecastDaily{
					Time:             []string{"2023-01-01"},
					Temperature2mMax: nil, // Incomplete data
				},
			}, nil
		},
	}
	weatherClient := NewWeatherClient(mockClient, "celsius", "kmh", "mm")
	_, err := weatherClient.GetDailyForecast(52.52, 13.41, 1)
	if err == nil {
		t.Fatal("Expected an error for incomplete daily data, but got nil")
	}
}

func TestGetDailyForecast_APIError(t *testing.T) {
	mockClient := &mockForecastClient{
		GetWeatherFunc: func(latitude, longitude float64, currentParameters, hourlyParameters, dailyParameters []string, temperatureUnit, windSpeedUnit, precipitationUnit string, pastDays, forecastDays int64) (*openmateo.ForecastResult, error) {
			return nil, errors.New("API error")
		},
	}
	weatherClient := NewWeatherClient(mockClient, "celsius", "kmh", "mm")
	_, err := weatherClient.GetDailyForecast(52.52, 13.41, 1)
	if err == nil {
		t.Fatal("Expected an error for API error, but got nil")
	}
}

func TestGetDailyForecast_InconsistentLengths(t *testing.T) {
	mockClient := &mockForecastClient{
		GetWeatherFunc: func(latitude, longitude float64, currentParameters, hourlyParameters, dailyParameters []string, temperatureUnit, windSpeedUnit, precipitationUnit string, pastDays, forecastDays int64) (*openmateo.ForecastResult, error) {
			return &openmateo.ForecastResult{
				Daily: &openmateo.ForecastDaily{
					Time:             []string{"2023-01-01"},
					Temperature2mMax: []float64{12.0, 13.0}, // Inconsistent length
					Temperature2mMin: []float64{2.0},
					Sunrise:          []string{"2023-01-01T07:00:00Z"},
					Sunset:           []string{"2023-01-01T17:00:00Z"},
					PrecipitationSum: []float64{0.1},
					PrecipitationProbabilityMean: []float64{10.0},
					WeatherCode:      []int{3},
					WindSpeed10mMax:    []float64{15.0},
				},
			}, nil
		},
	}
	weatherClient := NewWeatherClient(mockClient, "celsius", "kmh", "mm")
	_, err := weatherClient.GetDailyForecast(52.52, 13.41, 1)
	if err == nil {
		t.Fatal("Expected an error for inconsistent daily data lengths, but got nil")
	}
}
