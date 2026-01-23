package weather

import (
	"fmt"
	"time"

	"sky/internal/client/openmateo"
)

const openMeteoLayout = "2006-01-02T15:04"

// Units holds the unit strings for the weather data.
type Units struct {
	Temperature   string
	WindSpeed     string
	Precipitation string
}

// CurrentWeather represents the simplified current weather information
// that your application cares about.
type CurrentWeather struct {
	Temperature         float64
	Humidity            float64
	ApparentTemperature float64
	Precipitation       float64
	WindSpeed           float64
	WeatherDescription  string
	ObservationTime     time.Time
	IsDay               int
	Units               Units
}

// DailyForecast represents the simplified daily forecast information.
type DailyForecast struct {
	Date               time.Time
	MaxTemperature     float64
	MinTemperature     float64
	WeatherDescription string
	Sunrise            time.Time
	Sunset             time.Time
	PrecipitationSum   float64
	PrecipitationProb  float64 // Mean daily precipitation probability
	WindGusts          float64 // Max daily 10m wind speed
	Units              Units
}

// ForecastClient is an interface for a client that can fetch weather data.
type ForecastClient interface {
	GetWeather(
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

// WeatherClient is your application's client for weather-related operations.
// It composes a ForecastClient.
type WeatherClient struct {
	openmateoClient ForecastClient
}

// NewWeatherClient creates a new instance of the WeatherClient.
func NewWeatherClient(fc ForecastClient) *WeatherClient {
	return &WeatherClient{
		openmateoClient: fc,
	}
}

// GetCurrentWeather fetches the current weather conditions for a given location.
// This function abstracts away the specific parameters needed by the openmateo API.
func (w *WeatherClient) GetCurrentWeather(
	latitude, longitude float64,
	tempUnit, windUnit, precipUnit string,
) (*CurrentWeather, error) {
	// Define the specific current parameters we want from the Open-Meteo API
	currentParams := []string{
		"temperature_2m",
		"relative_humidity_2m",
		"weather_code",
		"is_day",
		"apparent_temperature",
		"precipitation",
		"wind_speed_10m",
	}

	// Call the low-level openmateo client's GetWeather function
	// We only care about current data, so other slices are empty.
	forecast, err := w.openmateoClient.GetWeather(
		latitude,
		longitude,
		currentParams,
		[]string{}, // No hourly data
		[]string{}, // No daily data
		tempUnit,
		windUnit,
		precipUnit,
		0,
		0,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to get raw weather data: %w", err)
	}

	if forecast.Current == nil {
		return nil, fmt.Errorf(
			"no current weather data returned for %.2f, %.2f",
			latitude,
			longitude,
		)
	}

	// Process the raw forecast data into your simplified CurrentWeather struct
	// You might have a helper function to map weather codes to descriptions
	// GetCurrentWeather always calls time.LoadLocation(forecast.Timezone). When the API response (or mocks) has an empty Timezone,
	// LoadLocation("") errors and prevents parsing the timestamp, causing GetCurrentWeather to fail even if the time is RFC3339.
	// Consider defaulting to time.UTC when forecast.Timezone is empty (or when LoadLocation fails) so callers aren’t forced to have a timezone in the response.
	var location *time.Location
	if forecast.Timezone == "" {
		// Default to UTC when no timezone is provided
		location = time.UTC
	} else {
		var locErr error
		location, locErr = time.LoadLocation(forecast.Timezone)
		if locErr != nil {
			// Fall back to UTC if the provided timezone cannot be loaded
			location = time.UTC
		}
	}

	obsTime, err := time.ParseInLocation(openMeteoLayout, forecast.Current.Time, location)
	if err != nil {
		// Fallback for full ISO8601 with offset
		obsTime, err = time.Parse(time.RFC3339, forecast.Current.Time)
		if err != nil {
			return nil, fmt.Errorf("failed to parse observation time: %w", err)
		}
	}

	weatherDesc := mapWeatherCodeToDescription(forecast.Current.WeatherCode)
	current := &CurrentWeather{
		Temperature:         forecast.Current.Temperature2m,
		Humidity:            forecast.Current.RelativeHumidity2m,
		ApparentTemperature: forecast.Current.ApparentTemperature,
		Precipitation:       forecast.Current.Precipitation,
		WindSpeed:           forecast.Current.WindSpeed10m,
		WeatherDescription:  weatherDesc,
		ObservationTime:     obsTime,
		IsDay:               forecast.Current.IsDay,
		Units: Units{
			Temperature:   forecast.CurrentUnits.Temperature2m,
			Precipitation: forecast.CurrentUnits.Precipitation,
			WindSpeed:     forecast.CurrentUnits.WindSpeed10m,
		},
	}

	return current, nil
}

// GetDailyForecast fetches the daily forecast for a given location for a specified number of days.
func (w *WeatherClient) GetDailyForecast(
	latitude, longitude float64,
	numDays int64,
	tempUnit, windUnit, precipUnit string,
) ([]DailyForecast, error) {
	if numDays < 1 || numDays > 16 { // Open-Meteo typically supports up to 16 days
		numDays = 1 // Default to 1 day
	}

	// Define the specific daily parameters we want from the Open-Meteo API
	dailyParams := []string{
		"temperature_2m_max",
		"temperature_2m_min",
		"weather_code",
		"sunrise",
		"sunset",
		"precipitation_sum",
		"precipitation_probability_mean",
		"wind_speed_10m_max",
	}

	forecast, err := w.openmateoClient.GetWeather(
		latitude,
		longitude,
		[]string{}, // No current data
		[]string{}, // No hourly data
		dailyParams,
		tempUnit,
		windUnit,
		precipUnit,
		0,       // No past days
		numDays, // Request numDays of forecast
	)
	if err != nil {
		return nil, fmt.Errorf("failed to get raw daily forecast data: %w", err)
	}

	// Add comprehensive nil checks for all slices we will access.
	if forecast.Daily == nil || forecast.Daily.Time == nil ||
		forecast.Daily.Temperature2mMax == nil || forecast.Daily.Temperature2mMin == nil ||
		forecast.Daily.WeatherCode == nil || forecast.Daily.Sunrise == nil ||
		forecast.Daily.Sunset == nil || forecast.Daily.PrecipitationSum == nil ||
		forecast.Daily.PrecipitationProbabilityMean == nil || forecast.Daily.WindSpeed10mMax == nil {
		return nil, fmt.Errorf("daily forecast data is incomplete or missing from API response")
	}

	// Check that all daily slices have the same length
	numDaysReturned := len(forecast.Daily.Time)
	if len(forecast.Daily.Temperature2mMax) != numDaysReturned ||
		len(forecast.Daily.Temperature2mMin) != numDaysReturned ||
		len(forecast.Daily.WeatherCode) != numDaysReturned ||
		len(forecast.Daily.Sunrise) != numDaysReturned ||
		len(forecast.Daily.Sunset) != numDaysReturned ||
		len(forecast.Daily.PrecipitationSum) != numDaysReturned ||
		len(forecast.Daily.PrecipitationProbabilityMean) != numDaysReturned ||
		len(forecast.Daily.WindSpeed10mMax) != numDaysReturned {
		return nil, fmt.Errorf("API returned daily forecast data with inconsistent lengths")
	}

	if len(forecast.Daily.Time) == 0 {
		return nil, fmt.Errorf(
			"no daily forecast data returned for %.2f, %.2f",
			latitude,
			longitude,
		)
	}

	location, err := time.LoadLocation(forecast.Timezone)
	if err != nil {
		return nil, fmt.Errorf("failed to load location from timezone: %w", err)
	}

	dailyForecasts := make([]DailyForecast, len(forecast.Daily.Time))
	units := Units{
		Temperature:   forecast.DailyUnits.Temperature2mMax,
		Precipitation: forecast.DailyUnits.PrecipitationSum,
		WindSpeed:     forecast.DailyUnits.WindSpeed10mMax,
	}

	for i := range forecast.Daily.Time {
		forecastDate, err := time.ParseInLocation("2006-01-02", forecast.Daily.Time[i], location)
		if err != nil {
			return nil, fmt.Errorf("failed to parse forecast date: %w", err)
		}

		sunriseTime, err := time.ParseInLocation(
			openMeteoLayout,
			forecast.Daily.Sunrise[i],
			location,
		)
		if err != nil {
			sunriseTime, err = time.Parse(time.RFC3339, forecast.Daily.Sunrise[i])
			if err != nil {
				return nil, fmt.Errorf("failed to parse sunrise time: %w", err)
			}
		}
		sunsetTime, err := time.ParseInLocation(openMeteoLayout, forecast.Daily.Sunset[i], location)
		if err != nil {
			sunsetTime, err = time.Parse(time.RFC3339, forecast.Daily.Sunset[i])
			if err != nil {
				return nil, fmt.Errorf("failed to parse sunset time: %w", err)
			}
		}

		dailyForecasts[i] = DailyForecast{
			Date:               forecastDate,
			MaxTemperature:     forecast.Daily.Temperature2mMax[i],
			MinTemperature:     forecast.Daily.Temperature2mMin[i],
			WeatherDescription: mapWeatherCodeToDescription(forecast.Daily.WeatherCode[i]),
			Sunrise:            sunriseTime,
			Sunset:             sunsetTime,
			PrecipitationSum:   forecast.Daily.PrecipitationSum[i],
			PrecipitationProb:  forecast.Daily.PrecipitationProbabilityMean[i],
			WindGusts:          forecast.Daily.WindSpeed10mMax[i],
			Units:              units,
		}
	}

	return dailyForecasts, nil
}

// mapWeatherCodeToDescription is a helper function to convert Open-Meteo
// weather codes into more human-readable descriptions.
// This would typically be more extensive.
func mapWeatherCodeToDescription(code int) string {
	switch code {
	case 0:
		return "Clear sky"
	case 1, 2, 3:
		return "Mainly clear, partly cloudy, and overcast"
	case 45, 48:
		return "Fog and depositing rime fog"
	case 51, 53, 55:
		return "Drizzle: Light, moderate, and dense intensity"
	case 56, 57:
		return "Freezing Drizzle: Light and dense intensity"
	case 61, 63, 65:
		return "Rain: Slight, moderate and heavy intensity"
	case 66, 67:
		return "Freezing Rain: Light and heavy intensity"
	case 71, 73, 75:
		return "Snow fall: Slight, moderate, and heavy intensity"
	case 77:
		return "Snow grains"
	case 80, 81, 82:
		return "Rain showers: Slight, moderate, and violent"
	case 85, 86:
		return "Snow showers: Slight and heavy"
	case 95:
		return "Thunderstorm: Slight or moderate"
	case 96, 99:
		return "Thunderstorm with slight and heavy hail"
	default:
		return fmt.Sprintf("Unknown weather code: %d", code)
	}
}