package weather

import (
	"sky/internal/client/openmateo"
	"testing"
)

// mockForecastClient is a mock implementation of the openmateo.ForecastClient
// for testing purposes.
type mockForecastClient struct {
	GetWeatherFunc func(
		latitude, longitude float64,
		currentParameters []string,
		hourlyParameters []string,
		dailyParameters []string,
		temperatureUnit string,
		windSpeedUnit string,
		precipitationUnit string,
		pastDays int64,
		forecastDays int64,
	) (*openmateo.ForecastResult, error)
}

func (m *mockForecastClient) GetWeather(
	latitude, longitude float64,
	currentParameters []string,
	hourlyParameters []string,
	dailyParameters []string,
	temperatureUnit string,
	windSpeedUnit string,
	precipitationUnit string,
	pastDays int64,
	forecastDays int64,
) (*openmateo.ForecastResult, error) {
	return m.GetWeatherFunc(latitude, longitude, currentParameters, hourlyParameters, dailyParameters, temperatureUnit, windSpeedUnit, precipitationUnit, pastDays, forecastDays)
}

func TestGetCurrentWeather_Success(t *testing.T) {
	mockClient := &mockForecastClient{
		GetWeatherFunc: func(latitude, longitude float64, currentParameters, hourlyParameters, dailyParameters []string, temperatureUnit, windSpeedUnit, precipitationUnit string, pastDays, forecastDays int64) (*openmateo.ForecastResult, error) {
			return &openmateo.ForecastResult{
				Timezone: "UTC",
				Current: &openmateo.ForecastCurrent{
					Time:                "2023-01-01T12:00:00Z",
					Temperature2m:       10.0,
					RelativeHumidity2m:  80.0,
					ApparentTemperature: 8.0,
					Precipitation:       0.5,
					WindSpeed10m:        5.0,
					WeatherCode:         3,
					IsDay:               1,
				},
			}, nil
		},
	}

	weatherClient := NewWeatherClient(mockClient, "celsius", "kmh", "mm")
	currentWeather, err := weatherClient.GetCurrentWeather(52.52, 13.41)
	if err != nil {
		t.Fatalf("GetCurrentWeather failed: %v", err)
	}

	if currentWeather.Temperature != 10.0 {
		t.Errorf("Expected Temperature to be 10.0, got %f", currentWeather.Temperature)
	}
	if currentWeather.Humidity != 80.0 {
		t.Errorf("Expected Humidity to be 80.0, got %f", currentWeather.Humidity)
	}
	if currentWeather.ApparentTemperature != 8.0 {
		t.Errorf("Expected ApparentTemperature to be 8.0, got %f", currentWeather.ApparentTemperature)
	}
	if currentWeather.Precipitation != 0.5 {
		t.Errorf("Expected Precipitation to be 0.5, got %f", currentWeather.Precipitation)
	}
	if currentWeather.WindSpeed != 5.0 {
		t.Errorf("Expected WindSpeed to be 5.0, got %f", currentWeather.WindSpeed)
	}
	if currentWeather.WeatherDescription != "Mainly clear, partly cloudy, and overcast" {
		t.Errorf("Expected WeatherDescription to be 'Mainly clear, partly cloudy, and overcast', got '%s'", currentWeather.WeatherDescription)
	}
	if currentWeather.IsDay != 1 {
		t.Errorf("Expected IsDay to be 1, got %d", currentWeather.IsDay)
	}
}

func TestGetDailyForecast_Success(t *testing.T) {
	mockClient := &mockForecastClient{
		GetWeatherFunc: func(latitude, longitude float64, currentParameters, hourlyParameters, dailyParameters []string, temperatureUnit, windSpeedUnit, precipitationUnit string, pastDays, forecastDays int64) (*openmateo.ForecastResult, error) {
			return &openmateo.ForecastResult{
				Timezone: "UTC",
				Daily: &openmateo.ForecastDaily{
					Time:             []string{"2023-01-01", "2023-01-02"},
					Temperature2mMax: []float64{12.0, 13.0},
					Temperature2mMin: []float64{2.0, 3.0},
					Sunrise:          []string{"2023-01-01T07:00:00Z", "2023-01-02T07:01:00Z"},
					Sunset:           []string{"2023-01-01T17:00:00Z", "2023-01-02T17:01:00Z"},
					PrecipitationSum: []float64{0.1, 0.2},
					PrecipitationProbabilityMean: []float64{10.0, 20.0},
					WeatherCode:      []int{3, 1},
					WindSpeed10mMax:    []float64{15.0, 16.0},
				},
			}, nil
		},
	}

	weatherClient := NewWeatherClient(mockClient, "celsius", "kmh", "mm")
	dailyForecast, err := weatherClient.GetDailyForecast(52.52, 13.41, 2)
	if err != nil {
		t.Fatalf("GetDailyForecast failed: %v", err)
	}

	if len(dailyForecast) != 2 {
		t.Fatalf("Expected 2 daily forecasts, got %d", len(dailyForecast))
	}

	if dailyForecast[0].MaxTemperature != 12.0 {
		t.Errorf("Expected MaxTemperature to be 12.0, got %f", dailyForecast[0].MaxTemperature)
	}
	if dailyForecast[1].MaxTemperature != 13.0 {
		t.Errorf("Expected MaxTemperature to be 13.0, got %f", dailyForecast[1].MaxTemperature)
	}
}

func TestGetDailyForecast_NumDaysZero(t *testing.T) {
	mockClient := &mockForecastClient{
		GetWeatherFunc: func(latitude, longitude float64, currentParameters, hourlyParameters, dailyParameters []string, temperatureUnit, windSpeedUnit, precipitationUnit string, pastDays, forecastDays int64) (*openmateo.ForecastResult, error) {
			if forecastDays != 1 {
				t.Errorf("Expected forecastDays to be 1, got %d", forecastDays)
			}
			return &openmateo.ForecastResult{
				Timezone: "UTC",
				Daily: &openmateo.ForecastDaily{
					Time:             []string{"2023-01-01"},
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
	dailyForecast, err := weatherClient.GetDailyForecast(52.52, 13.41, 0)
	if err != nil {
		t.Fatalf("GetDailyForecast with numDays=0 failed: %v", err)
	}

	if len(dailyForecast) != 1 {
		t.Errorf("Expected 1 daily forecast, got %d", len(dailyForecast))
	}
	if dailyForecast[0].MaxTemperature != 12.0 {
		t.Errorf("Expected MaxTemperature to be 12.0, got %f", dailyForecast[0].MaxTemperature)
	}
}

func TestGetDailyForecast_NumDaysOutOfRange(t *testing.T) {
	mockClient := &mockForecastClient{
		GetWeatherFunc: func(latitude, longitude float64, currentParameters, hourlyParameters, dailyParameters []string, temperatureUnit, windSpeedUnit, precipitationUnit string, pastDays, forecastDays int64) (*openmateo.ForecastResult, error) {
			if forecastDays != 1 {
				t.Errorf("Expected forecastDays to be 1, got %d", forecastDays)
			}
			return &openmateo.ForecastResult{
				Timezone: "UTC",
				Daily: &openmateo.ForecastDaily{
					Time:             []string{"2023-01-01"},
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

	dailyForecast, err := weatherClient.GetDailyForecast(52.52, 13.41, 17)
	if err != nil {
		t.Fatalf("GetDailyForecast with numDays=17 failed: %v", err)
	}
	if len(dailyForecast) != 1 {
		t.Errorf("Expected 1 daily forecast for numDays=17, got %d", len(dailyForecast))
	}
	if dailyForecast[0].MaxTemperature != 12.0 {
		t.Errorf("Expected MaxTemperature to be 12.0 for numDays=17, got %f", dailyForecast[0].MaxTemperature)
	}

	dailyForecast, err = weatherClient.GetDailyForecast(52.52, 13.41, -1)
	if err != nil {
		t.Fatalf("GetDailyForecast with numDays=-1 failed: %v", err)
	}
	if len(dailyForecast) != 1 {
		t.Errorf("Expected 1 daily forecast for numDays=-1, got %d", len(dailyForecast))
	}
	if dailyForecast[0].MaxTemperature != 12.0 {
		t.Errorf("Expected MaxTemperature to be 12.0 for numDays=-1, got %f", dailyForecast[0].MaxTemperature)
	}
}
