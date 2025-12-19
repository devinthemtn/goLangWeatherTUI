package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

// API constants
const (
	// OpenWeatherMap API (free tier, requires API key)
	openWeatherMapURL = "https://api.openweathermap.org/data/2.5/weather"
	geoLocationURL    = "https://api.openweathermap.org/geo/1.0/zip"

	// IP Geolocation service (free, no API key required)
	ipGeolocationURL = "https://ipapi.co/json/"

	// You need to set your OpenWeatherMap API key here
	// Get a free API key at: https://openweathermap.org/api
	// Can be set via environment variable OPENWEATHER_API_KEY or hardcoded below
)

// OpenWeatherMap API response structures
type OpenWeatherResponse struct {
	Weather []struct {
		Main        string `json:"main"`
		Description string `json:"description"`
	} `json:"weather"`
	Main struct {
		Temp      float64 `json:"temp"`
		FeelsLike float64 `json:"feels_like"`
		Humidity  int     `json:"humidity"`
		Pressure  float64 `json:"pressure"`
	} `json:"main"`
	Wind struct {
		Speed float64 `json:"speed"`
	} `json:"wind"`
	Visibility int    `json:"visibility"`
	Name       string `json:"name"`
	Sys        struct {
		Country string `json:"country"`
	} `json:"sys"`
}

// IP Geolocation response structure
type IPLocationResponse struct {
	City        string  `json:"city"`
	Region      string  `json:"region"`
	Country     string  `json:"country"`
	CountryName string  `json:"country_name"`
	Latitude    float64 `json:"latitude"`
	Longitude   float64 `json:"longitude"`
	Postal      string  `json:"postal"`
}

// ZIP code geolocation response
type ZipLocationResponse struct {
	Zip     string  `json:"zip"`
	Name    string  `json:"name"`
	Lat     float64 `json:"lat"`
	Lon     float64 `json:"lon"`
	Country string  `json:"country"`
}

// Message types for Bubble Tea
type WeatherMsg struct {
	WeatherData WeatherData
}

type ErrorMsg struct {
	Error string
}

// fetchWeatherByLocation gets weather for current location using IP geolocation
func (m Model) fetchWeatherByLocation() tea.Cmd {
	return func() tea.Msg {
		// Check if API key is set
		apiKey := getAPIKey()
		if apiKey == "your_api_key_here" || apiKey == "" {
			return ErrorMsg{Error: "API key not configured. Run with --debug flag for help: go run . --debug"}
		}

		// Get current location from IP
		location, err := getCurrentLocation()
		if err != nil {
			return ErrorMsg{Error: fmt.Sprintf("Failed to get location: %v", err)}
		}

		// Get weather data using coordinates
		weather, err := getWeatherByCoordinates(location.Latitude, location.Longitude)
		if err != nil {
			return ErrorMsg{Error: fmt.Sprintf("Failed to get weather: %v", err)}
		}

		// Use the location name from IP geolocation if available
		locationName := fmt.Sprintf("%s, %s", location.City, location.Region)
		if locationName == ", " {
			locationName = fmt.Sprintf("%s, %s", location.Country, location.CountryName)
		}

		weather.Location = locationName
		return WeatherMsg{WeatherData: *weather}
	}
}

// fetchWeatherByZip gets weather for a specific ZIP code
func (m Model) fetchWeatherByZip(zipCode string) tea.Cmd {
	return func() tea.Msg {
		// Check if API key is set
		apiKey := getAPIKey()
		if apiKey == "your_api_key_here" || apiKey == "" {
			return ErrorMsg{Error: "API key not configured. Run with --debug flag for help: go run . --debug"}
		}

		// First, get coordinates for the ZIP code
		coords, err := getCoordinatesFromZip(zipCode)
		if err != nil {
			return ErrorMsg{Error: fmt.Sprintf("Failed to find location for ZIP code: %v", err)}
		}

		// Get weather data using coordinates
		weather, err := getWeatherByCoordinates(coords.Lat, coords.Lon)
		if err != nil {
			return ErrorMsg{Error: fmt.Sprintf("Failed to get weather: %v", err)}
		}

		weather.Location = fmt.Sprintf("%s, %s %s", coords.Name, coords.Country, coords.Zip)
		return WeatherMsg{WeatherData: *weather}
	}
}

// getCurrentLocation gets the current location using IP geolocation
func getCurrentLocation() (*IPLocationResponse, error) {
	client := &http.Client{Timeout: 10 * time.Second}

	resp, err := client.Get(ipGeolocationURL)
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API request failed with status: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	var location IPLocationResponse
	if err := json.Unmarshal(body, &location); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	return &location, nil
}

// getCoordinatesFromZip gets coordinates for a ZIP code
func getCoordinatesFromZip(zipCode string) (*ZipLocationResponse, error) {
	client := &http.Client{Timeout: 10 * time.Second}
	apiKey := getAPIKey()

	// Build URL with ZIP code and API key
	u, err := url.Parse(geoLocationURL)
	if err != nil {
		return nil, err
	}

	q := u.Query()
	q.Set("zip", zipCode+",US") // Assuming US ZIP codes, can be extended
	q.Set("appid", apiKey)
	u.RawQuery = q.Encode()

	resp, err := client.Get(u.String())
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		if resp.StatusCode == 401 {
			return nil, fmt.Errorf("invalid API key (401). Check your OpenWeatherMap API key")
		} else if resp.StatusCode == 429 {
			return nil, fmt.Errorf("API rate limit exceeded (429). Try again in a few minutes")
		}
		return nil, fmt.Errorf("ZIP code lookup failed (status %d): %s", resp.StatusCode, string(body))
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	var coords ZipLocationResponse
	if err := json.Unmarshal(body, &coords); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	return &coords, nil
}

// getWeatherByCoordinates gets weather data for specific coordinates
func getWeatherByCoordinates(lat, lon float64) (*WeatherData, error) {
	client := &http.Client{Timeout: 10 * time.Second}
	apiKey := getAPIKey()

	// Build URL with coordinates and API key
	u, err := url.Parse(openWeatherMapURL)
	if err != nil {
		return nil, err
	}

	q := u.Query()
	q.Set("lat", fmt.Sprintf("%.6f", lat))
	q.Set("lon", fmt.Sprintf("%.6f", lon))
	q.Set("appid", apiKey)
	q.Set("units", "imperial") // Use Fahrenheit
	u.RawQuery = q.Encode()

	resp, err := client.Get(u.String())
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		if resp.StatusCode == 401 {
			return nil, fmt.Errorf("invalid API key (401). Verify your OpenWeatherMap API key is correct")
		} else if resp.StatusCode == 429 {
			return nil, fmt.Errorf("API rate limit exceeded (429). Free tier allows 60 calls/minute")
		} else if resp.StatusCode == 404 {
			return nil, fmt.Errorf("location not found (404). Check coordinates or ZIP code")
		}
		return nil, fmt.Errorf("weather API failed (status %d): %s", resp.StatusCode, string(body))
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	var owmResp OpenWeatherResponse
	if err := json.Unmarshal(body, &owmResp); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	// Convert to our WeatherData structure
	weather := &WeatherData{
		Location:    fmt.Sprintf("%s, %s", owmResp.Name, owmResp.Sys.Country),
		Temperature: owmResp.Main.Temp,
		FeelsLike:   owmResp.Main.FeelsLike,
		Humidity:    owmResp.Main.Humidity,
		WindSpeed:   owmResp.Wind.Speed,
		Pressure:    owmResp.Main.Pressure,
		Visibility:  float64(owmResp.Visibility) / 1609.34, // Convert meters to miles
		UVIndex:     0,                                     // OpenWeatherMap free tier doesn't include UV index
		Timestamp:   time.Now(),
	}

	// Set weather description
	if len(owmResp.Weather) > 0 {
		weather.Description = owmResp.Weather[0].Description
		// Capitalize first letter
		if len(weather.Description) > 0 {
			weather.Description = string(weather.Description[0]-32) + weather.Description[1:]
		}
	}

	return weather, nil
}

// Helper function to get API key from environment or fallback
func getAPIKey() string {
	// Try to get from environment variable first
	if key := os.Getenv("OPENWEATHER_API_KEY"); key != "" {
		return key
	}

	// Try to read from .env file
	if key := readAPIKeyFromEnvFile(); key != "" {
		return key
	}

	// Fallback to hardcoded value (not recommended for production)
	return "your_api_key_here"
}

// Helper function to read API key from .env file
func readAPIKeyFromEnvFile() string {
	file, err := os.Open(".env")
	if err != nil {
		return ""
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		// Skip empty lines and comments
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		// Look for OPENWEATHER_API_KEY=value
		if strings.HasPrefix(line, "OPENWEATHER_API_KEY=") {
			parts := strings.SplitN(line, "=", 2)
			if len(parts) == 2 {
				value := strings.TrimSpace(parts[1])
				// Remove quotes if present
				if len(value) >= 2 &&
					((strings.HasPrefix(value, "\"") && strings.HasSuffix(value, "\"")) ||
						(strings.HasPrefix(value, "'") && strings.HasSuffix(value, "'"))) {
					value = value[1 : len(value)-1]
				}
				return value
			}
		}
	}

	return ""
}

// Helper function to check if API key is configured
func isAPIKeyConfigured() bool {
	apiKey := getAPIKey()
	return apiKey != "" && apiKey != "your_api_key_here"
}
