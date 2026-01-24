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

// HourlyForecast represents the simplified hourly forecast information.
type HourlyForecast struct {
	DateTime            time.Time
	Temperature         float64
	Humidity            float64
	ApparentTemperature float64
	Cloudy              float64
	WindSpeed           float64
	Precipitation       float64
	SnowFall            float64
	PrecipitationProb   float64
	WeatherDescription  string
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
		pastHours int64,
		forecastHours int64,
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

	obsTime, err := parseTime(forecast.Current.Time, forecast.Timezone)
	if err != nil {
		return nil, err
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

// GetHourlyForecast fetches the hourly forecast for a given location.
func (w *WeatherClient) GetHourlyForecast(
	latitude, longitude float64,
	numHours int64,
	tempUnit, windUnit, precipUnit string,
) ([]HourlyForecast, error) {
	if numHours < 1 {
		numHours = 1
	}

	hourlyParams := []string{
		"temperature_2m",
		"relative_humidity_2m",
		"apparent_temperature",
		"cloud_cover",
		"wind_speed_10m",
		"precipitation",
		"snowfall",
		"precipitation_probability",
		"weather_code",
		"is_day",
	}

	forecast, err := w.openmateoClient.GetWeather(
		latitude,
		longitude,
		[]string{},
		hourlyParams,
		[]string{},
		tempUnit,
		windUnit,
		precipUnit,
		0,
		0,
		0,
		numHours,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to get raw hourly forecast data: %w", err)
	}

	// Add comprehensive nil checks for all slices we will access.
	if forecast.Hourly == nil || forecast.Hourly.Time == nil ||
		forecast.Hourly.Temperature2m == nil || forecast.Hourly.RelativeHumidity2m == nil ||
		forecast.Hourly.ApparentTemperature == nil || forecast.Hourly.CloudCover == nil ||
		forecast.Hourly.WindSpeed10m == nil || forecast.Hourly.Precipitation == nil ||
		forecast.Hourly.Snowfall == nil || forecast.Hourly.PrecipitationProbability == nil ||
		forecast.Hourly.WeatherCode == nil || forecast.Hourly.IsDay == nil {
		return nil, fmt.Errorf("hourly forecast data is incomplete or missing from API response")
	}

	// Check that all hourly slices have the same length
	numHoursReturned := len(forecast.Hourly.Time)
	if len(forecast.Hourly.Temperature2m) != numHoursReturned ||
		len(forecast.Hourly.RelativeHumidity2m) != numHoursReturned ||
		len(forecast.Hourly.ApparentTemperature) != numHoursReturned ||
		len(forecast.Hourly.CloudCover) != numHoursReturned ||
		len(forecast.Hourly.WindSpeed10m) != numHoursReturned ||
		len(forecast.Hourly.Precipitation) != numHoursReturned ||
		len(forecast.Hourly.Snowfall) != numHoursReturned ||
		len(forecast.Hourly.PrecipitationProbability) != numHoursReturned ||
		len(forecast.Hourly.WeatherCode) != numHoursReturned ||
		len(forecast.Hourly.IsDay) != numHoursReturned {
		return nil, fmt.Errorf("API returned hourly forecast data with inconsistent lengths")
	}

	if len(forecast.Hourly.Time) == 0 {
		return nil, fmt.Errorf(
			"no hourly forecast data returned for %.2f, %.2f",
			latitude,
			longitude,
		)
	}

	hourlyForecasts := make([]HourlyForecast, len(forecast.Hourly.Time))
	units := Units{
		Temperature:   forecast.HourlyUnits.Temperature2m,
		Precipitation: forecast.HourlyUnits.Precipitation,
		WindSpeed:     forecast.HourlyUnits.WindSpeed10m,
	}

	for i := range forecast.Hourly.Time {
		forecastTime, err := parseTime(forecast.Hourly.Time[i], forecast.Timezone)
		if err != nil {
			return nil, err
		}

		hourlyForecasts[i] = HourlyForecast{
			DateTime:            forecastTime,
			Temperature:         forecast.Hourly.Temperature2m[i],
			Humidity:            forecast.Hourly.RelativeHumidity2m[i],
			ApparentTemperature: forecast.Hourly.ApparentTemperature[i],
			Cloudy:              forecast.Hourly.CloudCover[i],
			WindSpeed:           forecast.Hourly.WindSpeed10m[i],
			Precipitation:       forecast.Hourly.Precipitation[i],
			SnowFall:            forecast.Hourly.Snowfall[i],
			PrecipitationProb:   forecast.Hourly.PrecipitationProbability[i],
			WeatherDescription:  mapWeatherCodeToDescription(forecast.Hourly.WeatherCode[i]),
			IsDay:               forecast.Hourly.IsDay[i],
			Units:               units,
		}
	}

	return hourlyForecasts, nil
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
		"snow_depth",
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
		0,
		0,
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

	dailyForecasts := make([]DailyForecast, len(forecast.Daily.Time))
	units := Units{
		Temperature:   forecast.DailyUnits.Temperature2mMax,
		Precipitation: forecast.DailyUnits.PrecipitationSum,
		WindSpeed:     forecast.DailyUnits.WindSpeed10mMax,
	}

	for i := range forecast.Daily.Time {
		forecastDate, err := parseTime(forecast.Daily.Time[i], forecast.Timezone)
		if err != nil {
			return nil, err
		}

		sunriseTime, err := parseTime(forecast.Daily.Sunrise[i], forecast.Timezone)
		if err != nil {
			return nil, err
		}
		sunsetTime, err := parseTime(forecast.Daily.Sunset[i], forecast.Timezone)
		if err != nil {
			return nil, err
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

func parseTime(timeStr, timezoneStr string) (time.Time, error) {
	location, err := loadTimezone(timezoneStr)
	if err != nil {
		return time.Time{}, err
	}

	parsedTime, err := time.ParseInLocation(openMeteoLayout, timeStr, location)
	if err != nil {
		parsedTime, err = time.Parse(time.RFC3339, timeStr)
		if err != nil {
			parsedTime, err = time.ParseInLocation("2006-01-02", timeStr, location)
			if err != nil {
				return time.Time{}, fmt.Errorf("failed to parse time: %w", err)
			}
		}
	}
	return parsedTime, nil
}

func loadTimezone(timezoneStr string) (*time.Location, error) {
	if timezoneStr == "" {
		return time.UTC, nil
	}
	location, err := time.LoadLocation(timezoneStr)
	if err != nil {
		return nil, fmt.Errorf("failed to load location from timezone: %w", err)
	}
	return location, nil
}

