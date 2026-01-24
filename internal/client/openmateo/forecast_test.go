package openmateo

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strconv"
	"strings"
	"testing"
)

func TestGetWeather_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		lat, _ := strconv.ParseFloat(r.URL.Query().Get("latitude"), 64)
		if lat != 52.52 {
			t.Errorf("Expected latitude to be '52.52', got '%f'", lat)
		}
		lon, _ := strconv.ParseFloat(r.URL.Query().Get("longitude"), 64)
		if lon != 13.41 {
			t.Errorf("Expected longitude to be '13.41', got '%f'", lon)
		}
		if r.URL.Query().Get("current") != "temperature_2m,weather_code" {
			t.Errorf(
				"Expected current to be 'temperature_2m,weather_code', got '%s'",
				r.URL.Query().Get("current"),
			)
		}
		if r.URL.Query().Get("hourly") != "temperature_2m,relative_humidity_2m" {
			t.Errorf(
				"Expected hourly to be 'temperature_2m,relative_humidity_2m', got '%s'",
				r.URL.Query().Get("hourly"),
			)
		}
		if r.URL.Query().Get("daily") != "weather_code,temperature_2m_max,temperature_2m_min" {
			t.Errorf(
				"Expected daily to be 'weather_code,temperature_2m_max,temperature_2m_min', got '%s'",
				r.URL.Query().Get("daily"),
			)
		}
		if r.URL.Query().Get("past_hours") != "2" {
			t.Errorf("Expected past_hours to be '2', got '%s'", r.URL.Query().Get("past_hours"))
		}
		if r.URL.Query().Get("forecast_hours") != "4" {
			t.Errorf("Expected forecast_hours to be '4', got '%s'", r.URL.Query().Get("forecast_hours"))
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = fmt.Fprintln(w, `{
			"latitude": 52.52,
			"longitude": 13.41,
			"generationtime_ms": 0.5,
			"timezone": "Europe/Berlin",
			"elevation": 74.0,
			"current": {
				"time": "2023-01-01T12:00",
				"temperature_2m": 10.0,
				"weather_code": 3
			},
			"hourly": {
				"time": ["2023-01-01T00:00"],
				"temperature_2m": [3.0],
				"relative_humidity_2m": [80.0]
			},
			"daily": {
				"time": ["2023-01-01"],
				"weather_code": [3],
				"temperature_2m_max": [12.0],
				"temperature_2m_min": [2.0]
			}
		}`)
	}))
	defer server.Close()

	client := NewForecastClient(server.Client())
	client.BaseURL = server.URL + "/"

	result, err := client.GetWeather(
		52.52,
		13.41,
		[]string{"temperature_2m", "weather_code"},
		[]string{"temperature_2m", "relative_humidity_2m"},
		[]string{"weather_code", "temperature_2m_max", "temperature_2m_min"},
		"celsius",
		"kmh",
		"mm",
		0,
		0,
		2,
		4,
	)
	if err != nil {
		t.Fatalf("GetWeather failed: %v", err)
	}

	expectedResult := &ForecastResult{
		GenerationTimeMs: 0.5,
		Timezone:         "Europe/Berlin", // Add expected timezone
		Elevation:        74.0,            // Add expected elevation
		Current: &ForecastCurrent{
			Time:          "2023-01-01T12:00",
			Temperature2m: 10.0,
			WeatherCode:   3,
		},
		Hourly: &ForecastHourly{
			Time:               []string{"2023-01-01T00:00"},
			Temperature2m:      []float64{3.0},
			RelativeHumidity2m: []float64{80.0},
		},
		Daily: &ForecastDaily{
			Time:             []string{"2023-01-01"},
			WeatherCode:      []int{3},
			Temperature2mMax: []float64{12.0},
			Temperature2mMin: []float64{2.0},
		},
	}

	// We don't care about comparing all fields, just the ones we requested
	if !reflect.DeepEqual(result.GenerationTimeMs, expectedResult.GenerationTimeMs) {
		t.Errorf(
			"Expected GenerationTimeMs '%v', got '%v'",
			expectedResult.GenerationTimeMs,
			result.GenerationTimeMs,
		)
	}
	if !reflect.DeepEqual(result.Timezone, expectedResult.Timezone) {
		t.Errorf(
			"Expected Timezone '%v', got '%v'",
			expectedResult.Timezone,
			result.Timezone,
		)
	}
	if !reflect.DeepEqual(result.Elevation, expectedResult.Elevation) {
		t.Errorf(
			"Expected Elevation '%v', got '%v'",
			expectedResult.Elevation,
			result.Elevation,
		)
	}
	if !reflect.DeepEqual(result.Current, expectedResult.Current) {
		t.Errorf("Expected Current '%v', got '%v'", expectedResult.Current, result.Current)
	}
	if !reflect.DeepEqual(result.Hourly, expectedResult.Hourly) {
		t.Errorf("Expected Hourly '%v', got '%v'", expectedResult.Hourly, result.Hourly)
	}
	if !reflect.DeepEqual(result.Daily, expectedResult.Daily) {
		t.Errorf("Expected Daily '%v', got '%v'", expectedResult.Daily, result.Daily)
	}
}

func TestGetWeather_APIError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = fmt.Fprintln(w, `{"error": true, "reason": "Internal Server Error"}`)
	}))
	defer server.Close()

	client := NewForecastClient(server.Client())
	client.BaseURL = server.URL + "/"

	_, err := client.GetWeather(52.52, 13.41, []string{}, []string{}, []string{}, "celsius", "kmh", "mm", 0, 0, 0, 0)
	if err == nil {
		t.Fatalf("expected GetWeather to return an error for API 500 response, got nil")
	}

	expectedErrMsg := "failed to get weather data"
	if !strings.Contains(err.Error(), expectedErrMsg) {
		t.Errorf("Expected error message to contain '%s', got '%s'", expectedErrMsg, err.Error())
	}
}

func TestGetWeather_MalformedJSON(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = fmt.Fprintln(w, `this is not valid json`)
	}))
	defer server.Close()

	client := NewForecastClient(server.Client())
	client.BaseURL = server.URL + "/"

	_, err := client.GetWeather(52.52, 13.41, []string{}, []string{}, []string{}, "celsius", "kmh", "mm", 0, 0, 0, 0)
	if err == nil {
		t.Fatal("Expected an error for malformed JSON, but got nil")
	}

	expectedErrMsg := "failed to unmarshal weather data"
	if !strings.Contains(err.Error(), expectedErrMsg) {
		t.Errorf("Expected error message to contain '%s', got '%s'", expectedErrMsg, err.Error())
	}
}

func TestGetWeather_NoParameters(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Query().Get("current") != "" {
			t.Errorf("Expected current to be empty, got '%s'", r.URL.Query().Get("current"))
		}
		if r.URL.Query().Get("hourly") != "" {
			t.Errorf("Expected hourly to be empty, got '%s'", r.URL.Query().Get("hourly"))
		}
		if r.URL.Query().Get("daily") != "" {
			t.Errorf("Expected daily to be empty, got '%s'", r.URL.Query().Get("daily"))
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = fmt.Fprintln(w, `{
			"latitude": 52.52,
			"longitude": 13.41,
			"generationtime_ms": 0.5,
			"timezone": "Europe/Berlin",
			"elevation": 74.0
		}`)
	}))
	defer server.Close()

	client := NewForecastClient(server.Client())
	client.BaseURL = server.URL + "/"

	result, err := client.GetWeather(52.52, 13.41, []string{}, []string{}, []string{}, "celsius", "kmh", "mm", 0, 0, 0, 0)
	if err != nil {
		t.Fatalf("GetWeather with no parameters failed: %v", err)
	}

	if result.GenerationTimeMs == 0 {
		t.Errorf("Expected GenerationTimeMs to be non-zero, got %f", result.GenerationTimeMs)
	}
	if result.Timezone == "" {
		t.Errorf("Expected Timezone to be non-empty, got '%s'", result.Timezone)
	}
	if result.Elevation == 0 {
		t.Errorf("Expected Elevation to be non-zero, got %f", result.Elevation)
	}
	if result.Current != nil {
		t.Errorf("Expected Current to be nil, got '%+v'", result.Current)
	}
	if result.Hourly != nil {
		t.Errorf("Expected Hourly to be nil, got '%+v'", result.Hourly)
	}
	if result.Daily != nil {
		t.Errorf("Expected Daily to be nil, got '%+v'", result.Daily)
	}
}
