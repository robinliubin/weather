package weather

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

const (
	NWSAPIBase = "https://api.weather.gov"
	UserAgent  = "weather-app/1.0"
)

// APIClient handles API requests to the NWS
type APIClient struct {
	client *http.Client
}

// NewAPIClient creates a new API client with default settings
func NewAPIClient() *APIClient {
	return &APIClient{
		client: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// makeNWSRequest makes a request to the NWS API with proper error handling
func (c *APIClient) makeNWSRequest(url string) (map[string]interface{}, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("User-Agent", UserAgent)
	req.Header.Set("Accept", "application/geo+json")

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API request failed with status: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var result map[string]interface{}
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, err
	}

	return result, nil
}

// formatAlert formats an alert feature into a readable string
func formatAlert(feature map[string]interface{}) string {
	props, ok := feature["properties"].(map[string]interface{})
	if !ok {
		return "Error: Invalid alert format"
	}

	event, _ := props["event"].(string)
	if event == "" {
		event = "Unknown"
	}

	areaDesc, _ := props["areaDesc"].(string)
	if areaDesc == "" {
		areaDesc = "Unknown"
	}

	severity, _ := props["severity"].(string)
	if severity == "" {
		severity = "Unknown"
	}

	description, _ := props["description"].(string)
	if description == "" {
		description = "No description available"
	}

	instruction, _ := props["instruction"].(string)
	if instruction == "" {
		instruction = "No specific instructions provided"
	}

	return fmt.Sprintf(`
Event: %s
Area: %s
Severity: %s
Description: %s
Instructions: %s
`, event, areaDesc, severity, description, instruction)
}

// GetAlerts gets weather alerts for a US state
func GetAlerts(state string) (string, error) {
	client := NewAPIClient()
	url := fmt.Sprintf("%s/alerts/active/area/%s", NWSAPIBase, state)

	data, err := client.makeNWSRequest(url)
	if err != nil {
		return "Unable to fetch alerts.", err
	}

	features, ok := data["features"].([]interface{})
	if !ok {
		return "No alerts found or invalid response format.", nil
	}

	if len(features) == 0 {
		return "No active alerts for this state.", nil
	}

	var alerts []string
	for _, f := range features {
		if feature, ok := f.(map[string]interface{}); ok {
			alerts = append(alerts, formatAlert(feature))
		}
	}

	return strings.Join(alerts, "\n---\n"), nil
}

// GetForecast gets weather forecast for specific coordinates
func GetForecast(latitude, longitude float64) (string, error) {
	client := NewAPIClient()
	
	// First get the forecast grid endpoint
	pointsURL := fmt.Sprintf("%s/points/%f,%f", NWSAPIBase, latitude, longitude)
	pointsData, err := client.makeNWSRequest(pointsURL)
	if err != nil {
		return "Unable to fetch forecast data for this location.", err
	}

	// Extract forecast URL from points response
	properties, ok := pointsData["properties"].(map[string]interface{})
	if !ok {
		return "Invalid response format from weather API.", nil
	}

	forecastURL, ok := properties["forecast"].(string)
	if !ok {
		return "Unable to retrieve forecast URL.", nil
	}

	// Get the forecast data
	forecastData, err := client.makeNWSRequest(forecastURL)
	if err != nil {
		return "Unable to fetch detailed forecast.", err
	}

	forecastProps, ok := forecastData["properties"].(map[string]interface{})
	if !ok {
		return "Invalid forecast data format.", nil
	}

	periods, ok := forecastProps["periods"].([]interface{})
	if !ok {
		return "Invalid forecast periods data.", nil
	}

	var forecasts []string
	for _, p := range periods {
		period, ok := p.(map[string]interface{})
		if !ok {
			continue
		}

		name, _ := period["name"].(string)
		temperature, _ := period["temperature"].(float64)
		temperatureUnit, _ := period["temperatureUnit"].(string)
		windSpeed, _ := period["windSpeed"].(string)
		windDirection, _ := period["windDirection"].(string)
		detailedForecast, _ := period["detailedForecast"].(string)

		forecast := fmt.Sprintf(`
%s:
Temperature: %.0f°%s
Wind: %s %s
Forecast: %s
`, name, temperature, temperatureUnit, windSpeed, windDirection, detailedForecast)

		forecasts = append(forecasts, forecast)
	}

	return strings.Join(forecasts, "\n---\n"), nil
}

// GeocodingResult represents a result from the Nominatim geocoding service
type GeocodingResult struct {
	Lat string `json:"lat"`
	Lon string `json:"lon"`
}

// GeocodeCity converts a city name to latitude and longitude coordinates
func GeocodeCity(city, state string) (float64, float64, error) {
	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	query := city
	if state != "" {
		query = fmt.Sprintf("%s, %s", city, state)
	}

	// Use OSM Nominatim for geocoding
	baseURL := "https://nominatim.openstreetmap.org/search"
	req, err := http.NewRequest("GET", baseURL, nil)
	if err != nil {
		return 0, 0, err
	}

	q := req.URL.Query()
	q.Add("q", query)
	q.Add("format", "json")
	q.Add("limit", "1")
	req.URL.RawQuery = q.Encode()

	req.Header.Set("User-Agent", UserAgent)

	resp, err := client.Do(req)
	if err != nil {
		return 0, 0, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return 0, 0, fmt.Errorf("geocoding request failed with status: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return 0, 0, err
	}

	var results []GeocodingResult
	if err := json.Unmarshal(body, &results); err != nil {
		return 0, 0, err
	}

	if len(results) == 0 {
		return 0, 0, fmt.Errorf("no results found for location: %s", query)
	}

	lat, err := strconv.ParseFloat(results[0].Lat, 64)
	if err != nil {
		return 0, 0, err
	}

	lon, err := strconv.ParseFloat(results[0].Lon, 64)
	if err != nil {
		return 0, 0, err
	}

	return lat, lon, nil
}

// GetForecastByCity gets a weather forecast for a city by name
func GetForecastByCity(city, state string) (string, error) {
	// First geocode the city to get coordinates
	lat, lon, err := GeocodeCity(city, state)
	if err != nil {
		return fmt.Sprintf("Unable to find coordinates for %s%s.", 
			city, 
			func() string {
				if state != "" {
					return ", " + state
				}
				return ""
			}()), 
			err
	}

	// Use the existing forecast function with the obtained coordinates
	return GetForecast(lat, lon)
}

// HTTP Handlers

// HandleAlerts handles requests for weather alerts
func HandleAlerts(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	state := r.URL.Query().Get("state")
	if state == "" {
		http.Error(w, "State parameter is required", http.StatusBadRequest)
		return
	}

	alerts, err := GetAlerts(state)
	if err != nil {
		http.Error(w, "Error fetching alerts: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/plain")
	w.Write([]byte(alerts))
}

// HandleForecast handles requests for weather forecasts by coordinates
func HandleForecast(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	latStr := r.URL.Query().Get("lat")
	lonStr := r.URL.Query().Get("lon")
	
	if latStr == "" || lonStr == "" {
		http.Error(w, "Latitude and longitude parameters are required", http.StatusBadRequest)
		return
	}

	lat, err := strconv.ParseFloat(latStr, 64)
	if err != nil {
		http.Error(w, "Invalid latitude value", http.StatusBadRequest)
		return
	}

	lon, err := strconv.ParseFloat(lonStr, 64)
	if err != nil {
		http.Error(w, "Invalid longitude value", http.StatusBadRequest)
		return
	}

	forecast, err := GetForecast(lat, lon)
	if err != nil {
		http.Error(w, "Error fetching forecast: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/plain")
	w.Write([]byte(forecast))
}

// HandleForecastByCity handles requests for weather forecasts by city name
func HandleForecastByCity(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	city := r.URL.Query().Get("city")
	if city == "" {
		http.Error(w, "City parameter is required", http.StatusBadRequest)
		return
	}

	state := r.URL.Query().Get("state")

	forecast, err := GetForecastByCity(city, state)
	if err != nil {
		http.Error(w, "Error fetching forecast: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/plain")
	w.Write([]byte(forecast))
}