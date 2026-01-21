package openmateo

// ForecastResult extracts the relevant parts of the forecast API response.
// It groups the optional forecast sections (current, hourly, daily) that can
// be selectively requested via API parameters.
type ForecastResult struct {
	GenerationTimeMs float64               `json:"generationtime_ms"`
	Timezone         string                `json:"timezone"`
	Elevation        float64               `json:"elevation"`
	Current          *ForecastCurrent      `json:"current,omitempty"`
	CurrentUnits     *ForecastCurrentUnits `json:"current_units,omitempty"`
	Hourly           *ForecastHourly       `json:"hourly,omitempty"`
	HourlyUnits      *ForecastHourlyUnits  `json:"hourly_units,omitempty"`
	Daily            *ForecastDaily        `json:"daily,omitempty"`
	DailyUnits       *ForecastDailyUnits   `json:"daily_units,omitempty"`
}

// ForecastCurrent holds the current weather conditions.
type ForecastCurrent struct {
	Time                string  `json:"time"`
	Temperature2m       float64 `json:"temperature_2m"`
	RelativeHumidity2m  float64 `json:"relative_humidity_2m"`
	Precipitation       float64 `json:"precipitation"`
	Snowfall            float64 `json:"snowfall"`
	WeatherCode         int     `json:"weather_code"`
	WindSpeed10m        float64 `json:"wind_speed_10m"`
	WindDirection10m    float64 `json:"wind_direction_10m"`
	IsDay               int     `json:"is_day"`
	ApparentTemperature float64 `json:"apparent_temperature"`
}

// ForecastCurrentUnits holds the units for the current weather conditions.
type ForecastCurrentUnits struct {
	Time                string `json:"time"`
	Temperature2m       string `json:"temperature_2m"`
	RelativeHumidity2m  string `json:"relative_humidity_2m"`
	Precipitation       string `json:"precipitation"`
	Snowfall            string `json:"snowfall"`
	WeatherCode         string `json:"weather_code"`
	WindSpeed10m        string `json:"wind_speed_10m"`
	WindDirection10m    string `json:"wind_direction_10m"`
	IsDay               string `json:"is_day"`
	ApparentTemperature string `json:"apparent_temperature"`
}

// ForecastHourly holds the time-series data for the hourly forecast.
type ForecastHourly struct {
	Time                     []string  `json:"time"`
	Temperature2m            []float64 `json:"temperature_2m"`
	RelativeHumidity2m       []float64 `json:"relative_humidity_2m"`
	Precipitation            []float64 `json:"precipitation"`
	WeatherCode              []int     `json:"weather_code"`
	WindSpeed10m             []float64 `json:"wind_speed_10m"`
	ApparentTemperature      []float64 `json:"apparent_temperature"`
	CloudCover               []float64 `json:"cloud_cover"`
	WindDirection10m         []float64 `json:"wind_direction_10m"`
	Snowfall                 []float64 `json:"snowfall"`
	PrecipitationProbability []float64 `json:"precipitation_probability"`
	SnowDepth                []float64 `json:"snow_depth"`
	IsDay                    []int     `json:"is_day"`
}

// ForecastHourlyUnits holds the units for the hourly forecast data.
type ForecastHourlyUnits struct {
	Time                     string `json:"time"`
	Temperature2m            string `json:"temperature_2m"`
	RelativeHumidity2m       string `json:"relative_humidity_2m"`
	Precipitation            string `json:"precipitation"`
	WeatherCode              string `json:"weather_code"`
	WindSpeed10m             string `json:"wind_speed_10m"`
	ApparentTemperature      string `json:"apparent_temperature"`
	CloudCover               string `json:"cloud_cover"`
	WindDirection10m         string `json:"wind_direction_10m"`
	Snowfall                 string `json:"snowfall"`
	PrecipitationProbability string `json:"precipitation_probability"`
	SnowDepth                string `json:"snow_depth"`
	IsDay                    string `json:"is_day"`
}

// ForecastDaily holds the time-series data for the daily forecast.
type ForecastDaily struct {
	Time                         []string  `json:"time"`
	Temperature2mMax             []float64 `json:"temperature_2m_max"`
	Temperature2mMin             []float64 `json:"temperature_2m_min"`
	Sunrise                      []string  `json:"sunrise"`
	Sunset                       []string  `json:"sunset"`
	DaylightDuration             []float64 `json:"daylight_duration"`
	PrecipitationSum             []float64 `json:"precipitation_sum"`
	SnowfallSum                  []float64 `json:"snowfall_sum"`
	PrecipitationProbabilityMean []float64 `json:"precipitation_probability_mean"`
	WeatherCode                  []int     `json:"weather_code"`
	WindSpeed10mMax              []float64 `json:"wind_speed_10m_max"`
	ApparentTemperatureMax       []float64 `json:"apparent_temperature_max"`
	ApparentTemperatureMin       []float64 `json:"apparent_temperature_min"`
}

// ForecastDailyUnits holds the units for the daily forecast data.
type ForecastDailyUnits struct {
	Time                         string `json:"time"`
	Temperature2mMax             string `json:"temperature_2m_max"`
	Temperature2mMin             string `json:"temperature_2m_min"`
	Sunrise                      string `json:"sunrise"`
	Sunset                       string `json:"sunset"`
	DaylightDuration             string `json:"daylight_duration"`
	PrecipitationSum             string `json:"precipitation_sum"`
	SnowfallSum                  string `json:"snowfall_sum"`
	PrecipitationProbabilityMean string `json:"precipitation_probability_mean"`
	WeatherCode                  string `json:"weather_code"`
	WindSpeed10mMax              string `json:"wind_speed_10m_max"`
	ApparentTemperatureMax       string `json:"apparent_temperature_max"`
	ApparentTemperatureMin       string `json:"apparent_temperature_min"`
}
