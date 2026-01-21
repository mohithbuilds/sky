# Open-Mateo

## Search

The API endpoint we are trying to hit: `https://geocoding-api.open-meteo.com/v1/search`

The goal of this for sky is to use it to retrieve the latitude and longitude [along with some extra data] for a location that we will be given by the end user.
- The extra data I'm referring to is elevation, timezone, population, country id, and country code

```go
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
```

JSON return object:
```json
{
    "results": [
        {
            "id": 2950159,
            "name": "Berlin",
            "latitude": 52.52437,
            "longitude": 13.41053,
            "elevation": 74.0,
            "feature_code": "PPLC",
            "country_code": "DE",
            "admin1_id": 2950157,
            "admin2_id": 0,
            "admin3_id": 6547383,
            "admin4_id": 6547539,
            "timezone": "Europe/Berlin",
            "population": 3426354,
            "postcodes": [
                "10967",
                "13347"
            ],
            "country_id": 2921044,
            "country": "Deutschland",
            "admin1": "Berlin",
            "admin2": "",
            "admin3": "Berlin, Stadt",
            "admin4": "Berlin"
        },

        {
        ...
        }
    ]
}
```

## Forecast

The API endpoint used for weather forecasts is: `https://api.open-meteo.com/v1/forecast`

This API provides current, hourly, and daily weather data for specified geographical coordinates.
The `sky` client interacts with this API to retrieve various weather parameters like temperature, precipitation, wind speed, weather codes, and more. The data is unmarshaled into `ForecastResult`, `ForecastCurrent`, `ForecastHourly`, and `ForecastDaily` structs, along with their respective unit structs (`ForecastCurrentUnits`, `ForecastHourlyUnits`, `ForecastDailyUnits`).

For more details on available parameters and response structure, refer to the official Open-Meteo Weather Forecast API Documentation: [https://open-meteo.com/en/docs](https://open-meteo.com/en/docs)

## Air Quality

The API endpoint for air quality data is: `https://air-quality-api.open-meteo.com/v1/air-quality`

This API provides hourly air quality data for specified geographical coordinates.
The `sky` client interacts with this API to retrieve parameters such as PM10, PM2.5, Ozone, Carbon Monoxide, and various pollen levels. The data is unmarshaled into `AirQualityResult` and `AirQualityHourly` structs, along with their unit struct (`AirQualityHourlyUnits`).

For more details on available parameters and response structure, refer to the official Open-Meteo Air Quality API Documentation: [https://www.open-meteo.com/en/docs/air-quality-api](https://www.open-meteo.com/en/docs/air-quality-api)
