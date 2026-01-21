package openmateo

import (
	"encoding/json"
	"fmt"
	"net/url"
	"strings"
)

// I also want to think about getting past days, forecast days, temperature unit, and precipitation unit
// Vision for this is to be the base getter and then have a bunch of wrappers that wrap this
func (fc *ForecastClient) GetWeather(
	latitude, longitude float64,
	currentParameters []string,
	hourlyParameters []string,
	dailyParameters []string,
) (*ForecastResult, error) {
	params := url.Values{}
	params.Add("latitude", fmt.Sprintf("%f", latitude))
	params.Add("longitude", fmt.Sprintf("%f", longitude))
	params.Add("temperature_unit", "celsius") // Default temperature unit
	params.Add("timezone", "auto")            // Default timezone

	if len(currentParameters) > 0 {
		params.Add("current", strings.Join(currentParameters, ","))
	}

	if len(hourlyParameters) > 0 {
		params.Add("hourly", strings.Join(hourlyParameters, ","))
	}

	if len(dailyParameters) > 0 {
		params.Add("daily", strings.Join(dailyParameters, ","))
	}

	fullURL := fc.BaseURL + "forecast?" + params.Encode()

	data, err := fc.doRequest(fullURL)
	if err != nil {
		return nil, fmt.Errorf("failed to get weather data: %w", err)
	}

	var result ForecastResult
	if err := json.Unmarshal(data, &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal weather data: %w", err)
	}

	return &result, nil
}
