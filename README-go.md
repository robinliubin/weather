# Weather (Go Version)

A Go-based weather application that provides access to the National Weather Service (NWS) API to retrieve weather alerts and forecasts. This is a Go port of the original Python application.

## Features

- Get active weather alerts for any US state
- Retrieve detailed weather forecasts for specific geographic coordinates
- Get 7-day weather forecasts for any city by name
- Cleanly formatted output for easy reading
- RESTful API endpoints for easy integration

## Installation

### Prerequisites

- Go 1.20 or higher
- Docker (optional, for containerized deployment)

### Running from Source

```bash
# Clone the repository
git clone https://github.com/robinliubin/weather.git
cd weather

# Build and run the application
go build -o weather-app .
./weather-app
```

### Using Docker

```bash
# Build and run with Docker
docker build -t weather-app .
docker run -p 8080:8080 weather-app

# Or using docker-compose
docker-compose up
```

## API Endpoints

The application exposes the following HTTP endpoints:

1. **Get Weather Alerts for a State**
   - Endpoint: `/alerts?state={state_code}`
   - Method: GET
   - Example: `/alerts?state=CA`

2. **Get Weather Forecast by Coordinates**
   - Endpoint: `/forecast?lat={latitude}&lon={longitude}`
   - Method: GET
   - Example: `/forecast?lat=37.7749&lon=-122.4194`

3. **Get Weather Forecast by City**
   - Endpoint: `/forecast/city?city={city_name}&state={state_code}`
   - Method: GET
   - Example: `/forecast/city?city=Seattle&state=WA`

## Example Queries

```bash
# Get active weather alerts for California
curl "http://localhost:8080/alerts?state=CA"

# Get weather forecast for San Francisco coordinates
curl "http://localhost:8080/forecast?lat=37.7749&lon=-122.4194"

# Get weather forecast for Seattle, WA
curl "http://localhost:8080/forecast/city?city=Seattle&state=WA"
```

## API Details

The application uses the National Weather Service (NWS) API:
- Base URL: https://api.weather.gov
- Endpoints used:
  - `/alerts/active/area/{state}` - For state alerts
  - `/points/{latitude},{longitude}` - For location information
  - `/forecast` - For detailed forecasts

For geocoding city names to coordinates, the application uses the OpenStreetMap Nominatim API.

## License

[Add your license information here]

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.