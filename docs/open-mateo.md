# Open-Mateo

## Search

The API endpoint we are trying to hit: `https://geocoding-api.open-meteo.com/v1/search`

The goal of this for sky is to use it to retrieve the latitude and longitude [along with some extra data] for a location that we will be given by the end user.
- The extra data I'm referring to is elevation, timezone, population, country id, and country code
```
type Location struct {
	ID          int     `json:"id"`
	Name        string  `json:"name"`
	Latitude    float64 `json:"latitude"`
	Longitude   float64 `json:"longitude"`
	Elevation   float64 `json:"elevation"`
	Timezone    string  `json:"timezone"`
	Population  int     `json:"population"`
	CountryCode int     `json:"country_code"`
	Country     string  `json:"country"`
}
```

JSON return object:
```
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

## Air Quality
