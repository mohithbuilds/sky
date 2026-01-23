package openmateo

import (
	"encoding/json"
	"fmt"
	"net/url"
	"strconv"
	"strings"
)

// TODO(sky): Add support for past_days, forecast_days, temperature_unit, and precipitation_unit parameters.
// The vision for this function is to be the base getter, with more specific wrapper functions built on top of it.
func (fc *ForecastClient) GetWeather(
	latitude, longitude float64,
	currentParameters []string,
	hourlyParameters []string,
	dailyParameters []string,
	temperatureUnit string,
	windSpeedUnit string,
	precipitationUnit string,
	pastDays int64,
	forecastDays int64,
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

	if temperatureUnit != "" {
		params.Add("temperature_unit", temperatureUnit)
	}

	if windSpeedUnit != "" {
		params.Add("wind_speed_unit", windSpeedUnit)
	}

	if precipitationUnit != "" {
		params.Add("precipitation_unit", precipitationUnit)
	}

	if pastDays >= 0 {
		params.Add("past_days", strconv.FormatInt(pastDays, 10))
	}

	if forecastDays >= 0 {
		params.Add("forecast_days", strconv.FormatInt(forecastDays, 10))
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
