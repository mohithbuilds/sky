package weather

import (
	"testing"

	"sky/internal/client/openmateo"
)

func TestGetCurrentWeather_Success(t *testing.T) {
	mockClient := &mockForecastClient{
		GetWeatherFunc: func(latitude, longitude float64, currentParameters, hourlyParameters, dailyParameters []string, temperatureUnit, windSpeedUnit, precipitationUnit string, pastDays, forecastDays, pastHours, forecastHours int64) (*openmateo.ForecastResult, error) {
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
				CurrentUnits: &openmateo.ForecastCurrentUnits{
					Temperature2m: "°C",
					WindSpeed10m:  "km/h",
					Precipitation: "mm",
				},
			}, nil
		},
	}

	weatherClient := NewWeatherClient(mockClient)
	currentWeather, err := weatherClient.GetCurrentWeather(52.52, 13.41, "celsius", "kmh", "mm")
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
		t.Errorf(
			"Expected ApparentTemperature to be 8.0, got %f",
			currentWeather.ApparentTemperature,
		)
	}
	if currentWeather.Precipitation != 0.5 {
		t.Errorf("Expected Precipitation to be 0.5, got %f", currentWeather.Precipitation)
	}
	if currentWeather.WindSpeed != 5.0 {
		t.Errorf("Expected WindSpeed to be 5.0, got %f", currentWeather.WindSpeed)
	}
	if currentWeather.WeatherDescription != "Mainly clear, partly cloudy, and overcast" {
		t.Errorf(
			"Expected WeatherDescription to be 'Mainly clear, partly cloudy, and overcast', got '%s'",
			currentWeather.WeatherDescription,
		)
	}
	if currentWeather.IsDay != 1 {
		t.Errorf("Expected IsDay to be 1, got %d", currentWeather.IsDay)
	}
	if currentWeather.Units.Temperature != "°C" {
		t.Errorf("Expected Temperature unit to be '°C', got '%s'", currentWeather.Units.Temperature)
	}
}

func TestGetDailyForecast_Success(t *testing.T) {
	mockClient := &mockForecastClient{
		GetWeatherFunc: func(latitude, longitude float64, currentParameters, hourlyParameters, dailyParameters []string, temperatureUnit, windSpeedUnit, precipitationUnit string, pastDays, forecastDays, pastHours, forecastHours int64) (*openmateo.ForecastResult, error) {
			return &openmateo.ForecastResult{
				Timezone: "UTC",
				Daily: &openmateo.ForecastDaily{
					Time:             []string{"2023-01-01", "2023-01-02"},
					Temperature2mMax: []float64{12.0, 13.0},
					Temperature2mMin: []float64{2.0, 3.0},
					Sunrise: []string{
						"2023-01-01T07:00:00Z",
						"2023-01-02T07:01:00Z",
					},
					Sunset: []string{
						"2023-01-01T17:00:00Z",
						"2023-01-02T17:01:00Z",
					},
					PrecipitationSum:             []float64{0.1, 0.2},
					PrecipitationProbabilityMean: []float64{10.0, 20.0},
					WeatherCode:                  []int{3, 1},
					WindSpeed10mMax:              []float64{15.0, 16.0},
				},
				DailyUnits: &openmateo.ForecastDailyUnits{
					Temperature2mMax: "°C",
				},
			}, nil
		},
	}

	weatherClient := NewWeatherClient(mockClient)
	dailyForecast, err := weatherClient.GetDailyForecast(52.52, 13.41, 2, "celsius", "kmh", "mm")
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
	if dailyForecast[0].Units.Temperature != "°C" {
		t.Errorf(
			"Expected Temperature unit to be '°C', got '%s'",
			dailyForecast[0].Units.Temperature,
		)
	}
}

func TestGetDailyForecast_NumDaysZero(t *testing.T) {
	mockClient := &mockForecastClient{
		GetWeatherFunc: func(latitude, longitude float64, currentParameters, hourlyParameters, dailyParameters []string, temperatureUnit, windSpeedUnit, precipitationUnit string, pastDays, forecastDays, pastHours, forecastHours int64) (*openmateo.ForecastResult, error) {
			if forecastDays != 1 {
				t.Errorf("Expected forecastDays to be 1, got %d", forecastDays)
			}
			return &openmateo.ForecastResult{
				Timezone: "UTC",
				Daily: &openmateo.ForecastDaily{
					Time:                         []string{"2023-01-01"},
					Temperature2mMax:             []float64{12.0},
					Temperature2mMin:             []float64{2.0},
					Sunrise:                      []string{"2023-01-01T07:00:00Z"},
					Sunset:                       []string{"2023-01-01T17:00:00Z"},
					PrecipitationSum:             []float64{0.1},
					PrecipitationProbabilityMean: []float64{10.0},
					WeatherCode:                  []int{3},
					WindSpeed10mMax:              []float64{15.0},
				},
				DailyUnits: &openmateo.ForecastDailyUnits{
					Temperature2mMax: "°C",
				},
			}, nil
		},
	}
	weatherClient := NewWeatherClient(mockClient)
	dailyForecast, err := weatherClient.GetDailyForecast(52.52, 13.41, 0, "celsius", "kmh", "mm")
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
		GetWeatherFunc: func(latitude, longitude float64, currentParameters, hourlyParameters, dailyParameters []string, temperatureUnit, windSpeedUnit, precipitationUnit string, pastDays, forecastDays, pastHours, forecastHours int64) (*openmateo.ForecastResult, error) {
			if forecastDays != 1 {
				t.Errorf("Expected forecastDays to be 1, got %d", forecastDays)
			}
			return &openmateo.ForecastResult{
				Timezone: "UTC",
				Daily: &openmateo.ForecastDaily{
					Time:                         []string{"2023-01-01"},
					Temperature2mMax:             []float64{12.0},
					Temperature2mMin:             []float64{2.0},
					Sunrise:                      []string{"2023-01-01T07:00:00Z"},
					Sunset:                       []string{"2023-01-01T17:00:00Z"},
					PrecipitationSum:             []float64{0.1},
					PrecipitationProbabilityMean: []float64{10.0},
					WeatherCode:                  []int{3},
					WindSpeed10mMax:              []float64{15.0},
				},
				DailyUnits: &openmateo.ForecastDailyUnits{
					Temperature2mMax: "°C",
				},
			}, nil
		},
	}
	weatherClient := NewWeatherClient(mockClient)

	dailyForecast, err := weatherClient.GetDailyForecast(52.52, 13.41, 17, "celsius", "kmh", "mm")
	if err != nil {
		t.Fatalf("GetDailyForecast with numDays=17 failed: %v", err)
	}
	if len(dailyForecast) != 1 {
		t.Errorf("Expected 1 daily forecast for numDays=17, got %d", len(dailyForecast))
	}
	if dailyForecast[0].MaxTemperature != 12.0 {
		t.Errorf(
			"Expected MaxTemperature to be 12.0 for numDays=17, got %f",
			dailyForecast[0].MaxTemperature,
		)
	}

	dailyForecast, err = weatherClient.GetDailyForecast(52.52, 13.41, -1, "celsius", "kmh", "mm")
	if err != nil {
		t.Fatalf("GetDailyForecast with numDays=-1 failed: %v", err)
	}
	if len(dailyForecast) != 1 {
		t.Errorf("Expected 1 daily forecast for numDays=-1, got %d", len(dailyForecast))
	}
	if dailyForecast[0].MaxTemperature != 12.0 {
		t.Errorf(
			"Expected MaxTemperature to be 12.0 for numDays=-1, got %f",
			dailyForecast[0].MaxTemperature,
		)
	}
}

func TestGetHourlyForecast_Success(t *testing.T) {
	mockClient := &mockForecastClient{
		GetWeatherFunc: func(latitude, longitude float64, currentParameters, hourlyParameters, dailyParameters []string, temperatureUnit, windSpeedUnit, precipitationUnit string, pastDays, forecastDays, pastHours, forecastHours int64) (*openmateo.ForecastResult, error) {
			return &openmateo.ForecastResult{
				Timezone: "UTC",
				Hourly: &openmateo.ForecastHourly{
					Time:                []string{"2023-01-01T12:00:00Z", "2023-01-01T13:00:00Z"},
					Temperature2m:       []float64{10.0, 11.0},
					RelativeHumidity2m:  []float64{80.0, 81.0},
					ApparentTemperature: []float64{8.0, 9.0},
					CloudCover:          []float64{50.0, 55.0},
					WindSpeed10m:        []float64{5.0, 6.0},
					Precipitation:       []float64{0.5, 0.6},
					Snowfall:            []float64{0.0, 0.0},
					PrecipitationProbability: []float64{10.0, 15.0},
					WeatherCode:         []int{3, 1},
					IsDay:               []int{1, 1},
				},
				HourlyUnits: &openmateo.ForecastHourlyUnits{
					Temperature2m: "°C",
					WindSpeed10m:  "km/h",
					Precipitation: "mm",
				},
			}, nil
		},
	}

	weatherClient := NewWeatherClient(mockClient)
	hourlyForecast, err := weatherClient.GetHourlyForecast(52.52, 13.41, 2, "celsius", "kmh", "mm")
	if err != nil {
		t.Fatalf("GetHourlyForecast failed: %v", err)
	}

	if len(hourlyForecast) != 2 {
		t.Fatalf("Expected 2 hourly forecasts, got %d", len(hourlyForecast))
	}

	if hourlyForecast[0].Temperature != 10.0 {
		t.Errorf("Expected Temperature to be 10.0, got %f", hourlyForecast[0].Temperature)
	}
	if hourlyForecast[1].Temperature != 11.0 {
		t.Errorf("Expected Temperature to be 11.0, got %f", hourlyForecast[1].Temperature)
	}
	if hourlyForecast[0].Units.Temperature != "°C" {
		t.Errorf(
			"Expected Temperature unit to be '°C', got '%s'",
			hourlyForecast[0].Units.Temperature,
		)
	}
}
