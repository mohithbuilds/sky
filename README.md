# Sky

A command-line weather tool written in Go.

## Data Provider

This project uses the [Open-Meteo API](https://open-meteo.com/) for weather forecast and geocoding data. It was chosen because it does not require an API key and provides a simple way to get weather data.

## Project Structure

The project follows the standard Go project layout:

```
├── cmd/sky/main.go     # Main application entry point
├── internal/           # Private application logic
│   ├── client/         # Client for interacting with external APIs
│   │   └── openmeteo/  # Open-Meteo API client
│   └── weather/        # Core weather application logic
├── go.mod              # Go module definition
└── README.md
```

## Getting Started

### Prerequisites

*   [Go](https://golang.org/doc/install)

### Building

```sh
go build -o sky cmd/sky/main.go
```

### Running

```sh
./sky
```

## Roadmap

*   [ ] Implement the Open-Meteo client in `internal/client/openmeteo`.
    *   [ ] `search.go`: Functionality to search for locations.
    *   [ ] `forecast.go`: Functionality to get the weather forecast for a location.
*   [ ] Integrate the Open-Meteo client with the main application in `cmd/sky/main.go`.
*   [ ] Implement the core weather logic in `internal/weather`.
*   [ ] Parse and display the weather information to the user.
