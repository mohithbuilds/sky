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

func TestGetAirQuality_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		lat, _ := strconv.ParseFloat(r.URL.Query().Get("latitude"), 64)
		if lat != 52.52 {
			t.Errorf("Expected latitude to be '52.52', got '%f'", lat)
		}
		lon, _ := strconv.ParseFloat(r.URL.Query().Get("longitude"), 64)
		if lon != 13.41 {
			t.Errorf("Expected longitude to be '13.41', got '%f'", lon)
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = fmt.Fprintln(w, `{
			"latitude": 52.52,
			"longitude": 13.41,
			"generationtime_ms": 0.2,
			"hourly_units": {
				"time": "iso8601",
				"pm2_5": "μg/m³"
			},
			"hourly": {
				"time": ["2023-01-01T00:00"],
				"pm2_5": [8.0]
			}
		}`)
	}))
	defer server.Close()

	client := NewAirQualityClient(server.Client())
	client.BaseURL = server.URL + "/"

	result, err := client.GetAirQuality(52.52, 13.41, []string{"pm2_5"})
	if err != nil {
		t.Fatalf("GetAirQuality failed: %v", err)
	}

	expectedResult := &AirQualityResult{
		GenerationTimeMs: 0.2,
		Hourly: Hourly{
			Time: []string{"2023-01-01T00:00"},
			PM25: []float64{8.0},
		},
		HourlyUnits: HourlyUnits{
			Time: "iso8601",
			PM25: "μg/m³",
		},
	}

	if !reflect.DeepEqual(result, expectedResult) {
		t.Errorf("Expected result '%v', got '%v'", expectedResult, result)
	}
}

func TestGetAirQuality_APIError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = fmt.Fprintln(w, `{"error": true, "reason": "Internal Server Error"}`)
	}))
	defer server.Close()

	client := NewAirQualityClient(server.Client())
	client.BaseURL = server.URL + "/"

	_, err := client.GetAirQuality(52.52, 13.41, []string{"pm2_5"})
	if err == nil {
		t.Fatal("Expected an error for API failure, but got nil")
	}

	expectedErrMsg := "failed to get air quality data"
	if !strings.Contains(err.Error(), expectedErrMsg) {
		t.Errorf("Expected error message to contain '%s', got '%s'", expectedErrMsg, err.Error())
	}
}

func TestGetAirQuality_MalformedJSON(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = fmt.Fprintln(w, `this is not valid json`)
	}))
	defer server.Close()

	client := NewAirQualityClient(server.Client())
	client.BaseURL = server.URL + "/"

	_, err := client.GetAirQuality(52.52, 13.41, []string{"pm2_5"})
	if err == nil {
		t.Fatal("Expected an error for malformed JSON, but got nil")
	}

	expectedErrMsg := "failed to unmarshal air quality data"
	if !strings.Contains(err.Error(), expectedErrMsg) {
		t.Errorf("Expected error message to contain '%s', got '%s'", expectedErrMsg, err.Error())
	}
}
