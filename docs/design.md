# Design Approach

## Data Provider

Currently want to focus on using Open-Mateo.
1. No API key is required
2. Provides a `search` API
  * Won't need to manage any local mappings or extra dependencies for longitude and latitude by name
3. Have many weather forecast models as options

## Repository Structure

This is my first actual project with Go, so will be iterating. But right now the basic idea after looking at [`golang-standards/project-layout`](https://github.com/golang-standards/project-layout) is:

```
├── cmd/
│   └── sky/
│       └── main.go
├── internal/
│   ├── client/
│   └── weather/
├── configs/
├── docs/
├── scripts/
├── .gitignore
├── go.mod
├── go.sum
└── README.md
```

### Explanation of the structure:

* **`/cmd/sky/main.go`**: This is the entry point for the application.
* **`/internal`**: This directory contains all the private application code. The Go compiler will prevent other projects from importing packages from the `internal` directory.
  * **`/internal/client`**: A place to put the code that interacts with the Open-Meteo API.
  * **`/internal/weather`**: For core application logic related to weather.
* **`/configs`**: This directory can hold configuration files.
* **`/docs`**: For your project's documentation, like the `design.md` you already have.
* **`/scripts`**: A place for shell scripts for building, installing, or other automation.
* **`go.mod` and `go.sum`**: These files are essential for managing the project's dependencies.
* **`README.md`**: A good `README.md` is crucial for any project. File to explain what the project does and how to use it.

## Weather Forecast Client Vision

The `GetWeather` function in the `openmateo` client is designed to be a base getter. More specific weather retrieval functions (wrappers) should be built on top of it, abstracting away the details of parameter selection for current, hourly, and daily forecasts, and potentially adding support for past days, forecast days, and specific temperature/precipitation units.