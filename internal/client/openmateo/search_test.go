package openmateo

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestSearch(t *testing.T) {
	// Start a local HTTP server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Check that the query parameters are correct
		if r.URL.Query().Get("name") != "Berlin" {
			t.Errorf("Expected name to be 'Berlin', got '%s'", r.URL.Query().Get("name"))
		}
		if r.URL.Query().Get("count") != "1" {
			t.Errorf("Expected count to be '1', got '%s'", r.URL.Query().Get("count"))
		}

		// Send response to be tested
		w.Header().Set("Content-Type", "application/json")
		_, _ = fmt.Fprintln(w, `{
			"results": [
				{
					"id": 2950159,
					"name": "Berlin",
					"latitude": 52.52437,
					"longitude": 13.41053,
					"country_code": "DE"
				}
			]
		}`)
	}))
	// Close the server when test finishes
	defer server.Close()

	// Create a GeocodingClient with the test server's URL
	client := NewGeocodingClient(server.Client())
	client.BaseURL = server.URL // Set BaseURL to the test server's URL

	// Call the Search function
	location, err := client.Search("Berlin")
	if err != nil {
		t.Fatalf("Search failed: %v", err)
	}

	// Check the results
	if location.Name != "Berlin" {
		t.Errorf("Expected location name 'Berlin', got '%s'", location.Name)
	}
	if location.Latitude != 52.52437 {
		t.Errorf("Expected latitude 52.52437, got '%f'", location.Latitude)
	}
	if location.ID != 2950159 {
		t.Errorf("Expected ID 2950159, got '%d'", location.ID)
	}
	if location.CountryCode != "DE" {
		t.Errorf("Expected CountryCode DE, got '%s'", location.CountryCode)
	}
}

func TestSearch_NoLocationFound(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = fmt.Fprintln(w, `{"results": []}`) // Empty results array
	}))
	defer server.Close()

	client := NewGeocodingClient(server.Client())
	client.BaseURL = server.URL

	_, err := client.Search("NonExistentCity")
	if err == nil {
		t.Fatal("Expected an error for no location found, but got nil")
	}

	expectedErrMsg := "no location found"
	if !strings.Contains(err.Error(), expectedErrMsg) {
		t.Errorf("Expected error message to contain '%s', got '%s'", expectedErrMsg, err.Error())
	}
}

func TestSearch_APIError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(
			http.StatusInternalServerError,
		) // Simulate a 500 Internal Server Error
		_, _ = fmt.Fprintln(
			w,
			`{"error": true, "reason": "Internal Server Error"}`,
		) // Provide a response body
	}))
	defer server.Close()

	client := NewGeocodingClient(server.Client())
	client.BaseURL = server.URL

	_, err := client.Search("AnyCity")
	if err == nil {
		t.Fatal("Expected an error for API failure, but got nil")
	}

	expectedErrMsg := "API error (500 Internal Server Error): Internal Server Error"
	if !strings.Contains(err.Error(), expectedErrMsg) {
		t.Errorf("Expected error message to contain '%s', got '%s'", expectedErrMsg, err.Error())
	}
}

func TestSearch_MalformedJSON(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = fmt.Fprintln(w, `this is not valid json`) // Malformed JSON response
	}))
	defer server.Close()

	client := NewGeocodingClient(server.Client())
	client.BaseURL = server.URL

	_, err := client.Search("AnyCity")
	if err == nil {
		t.Fatal("Expected an error for malformed JSON, but got nil")
	}

	expectedErrMsg := "failed to unmarshal search response"
	if !strings.Contains(err.Error(), expectedErrMsg) {
		t.Errorf("Expected error message to contain '%s', got '%s'", expectedErrMsg, err.Error())
	}
}
