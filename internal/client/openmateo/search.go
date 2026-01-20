package openmateo

import (
	"encoding/json"
	"fmt"
	"net/url"
)

type Results struct {
	Locations []Location `json:"results"`
}

type Location struct {
	ID          int     `json:"id"`
	Name        string  `json:"name"`
	Latitude    float64 `json:"latitude"`
	Longitude   float64 `json:"longitude"`
	Elevation   float64 `json:"elevation"`
	Timezone    string  `json:"timezone"`
	Population  int     `json:"population"`
	CountryCode string  `json:"country_code"`
	Country     string  `json:"country"`
}

func (gc *GeocodingClient) Search(locationName string) (*Location, error) {
	var searchURL string = fmt.Sprintf("%s/search?name=%s&count=1", gc.BaseURL, url.QueryEscape(locationName))

	data, err := gc.doRequest(searchURL)
	if err != nil {
		return nil, fmt.Errorf("Search request for %s failed: %w", locationName, err)
	}

	var result Results
	err = json.Unmarshal(data, &result)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal search response for %s: %w", locationName, err)
	}

	if len(result.Locations) == 0 {
		return nil, fmt.Errorf("no location found for %s", locationName)
	}

	return &result.Locations[0], nil
}
