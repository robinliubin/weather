# Weather

A Python-based weather application that provides access to the National Weather Service (NWS) API to retrieve weather alerts and forecasts.

## Features

- Get active weather alerts for any US state
- Retrieve detailed weather forecasts for specific geographic coordinates
- Get 7-day weather forecasts for any city by name
- Cleanly formatted output for easy reading
- Built as an MCP (Anthropic's Model Control Protocol) tool for use with Claude

## Installation

Requires Python 3.10 or higher.

```bash
# Clone the repository
git clone https://github.com/robinliubin/weather.git
cd weather

# Install dependencies
pip install -e .
```

## Usage

### As an MCP Tool with Claude

This application is designed to be used as a tool with Claude via the MCP protocol. Once installed, Claude can access weather data through three main functions:

1. `get_alerts(state)` - Get active weather alerts for a US state
2. `get_forecast(latitude, longitude)` - Get a detailed weather forecast for specific coordinates
3. `get_forecast_by_city(city, state)` - Get a 7-day weather forecast for a city by name

### Example Queries for Claude

```
What weather alerts are active in CA?
What's the weather forecast for latitude 37.7749 and longitude -122.4194 (San Francisco)?
What's the 7-day weather forecast for Seattle, WA?
```

### Running Locally

You can also run the application directly:

```bash
python main.py
```

## API Details

The application uses the National Weather Service (NWS) API:
- Base URL: https://api.weather.gov
- Endpoints used:
  - `/alerts/active/area/{state}` - For state alerts
  - `/points/{latitude},{longitude}` - For location information
  - `/forecast` - For detailed forecasts

## Dependencies

- httpx: For asynchronous HTTP requests
- mcp[cli]: Anthropic's Model Control Protocol implementation

## License

[Add your license information here]

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.