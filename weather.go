package main

import (
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

	// IP geolocation services (free, no API key required), tried in order
	// since free tiers rate-limit aggressively.
	ipwhoURL  = "https://ipwho.is/"
	ipinfoURL = "https://ipinfo.io/json"
	ipapiURL  = "https://ipapi.co/json/"

	// OpenWeatherMap API key: get one free at https://openweathermap.org/api
	// and store it encrypted with sops+age (see README / `make secrets`).
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
			if zip := getDefaultZip(); zip != "" {
				return m.fetchWeatherByZip(zip)()
			}
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

// getCurrentLocation gets the current location using IP geolocation,
// trying several free providers until one succeeds.
func getCurrentLocation() (*IPLocationResponse, error) {
	providers := []struct {
		name  string
		fetch func() (*IPLocationResponse, error)
	}{
		{"ipwho.is", fetchIPWho},
		{"ipinfo.io", fetchIPInfo},
		{"ipapi.co", fetchIPAPI},
	}

	var errs []string
	for _, p := range providers {
		loc, err := p.fetch()
		if err == nil && (loc.Latitude != 0 || loc.Longitude != 0) {
			return loc, nil
		}
		if err == nil {
			err = fmt.Errorf("no coordinates returned")
		}
		errs = append(errs, fmt.Sprintf("%s: %v", p.name, err))
	}
	return nil, fmt.Errorf("all geolocation providers failed (%s)", strings.Join(errs, "; "))
}

// getJSON GETs url and decodes the JSON body into v.
func getJSON(url string, v any) error {
	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Get(url)
	if err != nil {
		return fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("status %d", resp.StatusCode)
	}
	if err := json.NewDecoder(resp.Body).Decode(v); err != nil {
		return fmt.Errorf("failed to parse response: %w", err)
	}
	return nil
}

func fetchIPWho() (*IPLocationResponse, error) {
	var r struct {
		Success     bool    `json:"success"`
		Message     string  `json:"message"`
		City        string  `json:"city"`
		Region      string  `json:"region"`
		CountryCode string  `json:"country_code"`
		Country     string  `json:"country"`
		Latitude    float64 `json:"latitude"`
		Longitude   float64 `json:"longitude"`
		Postal      string  `json:"postal"`
	}
	if err := getJSON(ipwhoURL, &r); err != nil {
		return nil, err
	}
	if !r.Success {
		return nil, fmt.Errorf("%s", r.Message)
	}
	return &IPLocationResponse{
		City: r.City, Region: r.Region, Country: r.CountryCode, CountryName: r.Country,
		Latitude: r.Latitude, Longitude: r.Longitude, Postal: r.Postal,
	}, nil
}

func fetchIPInfo() (*IPLocationResponse, error) {
	var r struct {
		City    string `json:"city"`
		Region  string `json:"region"`
		Country string `json:"country"`
		Loc     string `json:"loc"` // "lat,lon"
		Postal  string `json:"postal"`
	}
	if err := getJSON(ipinfoURL, &r); err != nil {
		return nil, err
	}
	var lat, lon float64
	if _, err := fmt.Sscanf(r.Loc, "%f,%f", &lat, &lon); err != nil {
		return nil, fmt.Errorf("bad loc %q", r.Loc)
	}
	return &IPLocationResponse{
		City: r.City, Region: r.Region, Country: r.Country,
		Latitude: lat, Longitude: lon, Postal: r.Postal,
	}, nil
}

func fetchIPAPI() (*IPLocationResponse, error) {
	var loc IPLocationResponse
	if err := getJSON(ipapiURL, &loc); err != nil {
		return nil, err
	}
	return &loc, nil
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

// Helper function to get API key. Lookup order:
//  1. OPENWEATHER_API_KEY environment variable
//  2. sops+age encrypted secrets.enc.env (recommended)
//  3. plaintext .env file (legacy, gitignored)
func getAPIKey() string {
	if key := os.Getenv("OPENWEATHER_API_KEY"); key != "" {
		return key
	}

	if key := readAPIKeyFromSops(); key != "" {
		return key
	}

	if key := readAPIKeyFromEnvFile(); key != "" {
		return key
	}

	return "your_api_key_here"
}

// Helper function to read API key from a plaintext .env file (legacy).
func readAPIKeyFromEnvFile() string {
	data, err := os.ReadFile(".env")
	if err != nil {
		return ""
	}
	return parseDotenvValue(string(data), "OPENWEATHER_API_KEY")
}

// Helper function to check if API key is configured
func isAPIKeyConfigured() bool {
	apiKey := getAPIKey()
	return apiKey != "" && apiKey != "your_api_key_here"
}
