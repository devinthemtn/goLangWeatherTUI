package main

import (
	"os"
	"testing"
	"time"
)

func TestGetAPIKey(t *testing.T) {
	tests := []struct {
		name        string
		envValue    string
		expected    string
		description string
	}{
		{
			name:        "API key from environment",
			envValue:    "test_api_key_123",
			expected:    "test_api_key_123",
			description: "Should return API key from environment variable",
		},
		{
			name:        "No environment variable",
			envValue:    "",
			expected:    "your_api_key_here",
			description: "Should return default value when no env var is set",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Isolate from any real secrets.enc.env / .env in the repo
			chdirTemp(t)

			// Store original value
			originalValue := os.Getenv("OPENWEATHER_API_KEY")
			defer os.Setenv("OPENWEATHER_API_KEY", originalValue)

			// Set test environment
			if tt.envValue != "" {
				os.Setenv("OPENWEATHER_API_KEY", tt.envValue)
			} else {
				os.Unsetenv("OPENWEATHER_API_KEY")
			}

			// Test the function
			result := getAPIKey()
			if result != tt.expected {
				t.Errorf("getAPIKey() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestIsAPIKeyConfigured(t *testing.T) {
	tests := []struct {
		name     string
		envValue string
		expected bool
	}{
		{
			name:     "Valid API key",
			envValue: "sk-1234567890abcdef",
			expected: true,
		},
		{
			name:     "Default placeholder value",
			envValue: "your_api_key_here",
			expected: false,
		},
		{
			name:     "Empty API key",
			envValue: "",
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Isolate from any real secrets.enc.env / .env in the repo
			chdirTemp(t)

			// Store original value
			originalValue := os.Getenv("OPENWEATHER_API_KEY")
			defer os.Setenv("OPENWEATHER_API_KEY", originalValue)

			// Set test environment
			if tt.envValue != "" {
				os.Setenv("OPENWEATHER_API_KEY", tt.envValue)
			} else {
				os.Unsetenv("OPENWEATHER_API_KEY")
			}

			// Test the function
			result := isAPIKeyConfigured()
			if result != tt.expected {
				t.Errorf("isAPIKeyConfigured() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestWeatherDataStructure(t *testing.T) {
	// Test that WeatherData can be created and populated
	weather := WeatherData{
		Location:    "New York, NY",
		Temperature: 72.5,
		Description: "Clear sky",
		Humidity:    45,
		WindSpeed:   5.2,
		Pressure:    1013.25,
		FeelsLike:   75.0,
		Visibility:  10.0,
		UVIndex:     6.0,
		Timestamp:   time.Now(),
	}

	// Verify all fields are set correctly
	if weather.Location != "New York, NY" {
		t.Errorf("Expected location 'New York, NY', got %s", weather.Location)
	}
	if weather.Temperature != 72.5 {
		t.Errorf("Expected temperature 72.5, got %f", weather.Temperature)
	}
	if weather.Description != "Clear sky" {
		t.Errorf("Expected description 'Clear sky', got %s", weather.Description)
	}
	if weather.Humidity != 45 {
		t.Errorf("Expected humidity 45, got %d", weather.Humidity)
	}
	if weather.WindSpeed != 5.2 {
		t.Errorf("Expected wind speed 5.2, got %f", weather.WindSpeed)
	}
	if weather.Pressure != 1013.25 {
		t.Errorf("Expected pressure 1013.25, got %f", weather.Pressure)
	}
	if weather.FeelsLike != 75.0 {
		t.Errorf("Expected feels like 75.0, got %f", weather.FeelsLike)
	}
	if weather.Visibility != 10.0 {
		t.Errorf("Expected visibility 10.0, got %f", weather.Visibility)
	}
	if weather.UVIndex != 6.0 {
		t.Errorf("Expected UV index 6.0, got %f", weather.UVIndex)
	}
}

func TestModelStateTransitions(t *testing.T) {
	// Test initial model creation
	model := initialModel()

	// Verify initial state
	if model.state != MenuState {
		t.Errorf("Expected initial state to be MenuState, got %v", model.state)
	}
	if model.weather != nil {
		t.Errorf("Expected initial weather to be nil, got %v", model.weather)
	}
	if model.error != "" {
		t.Errorf("Expected initial error to be empty, got %s", model.error)
	}
	if model.loading != false {
		t.Errorf("Expected initial loading to be false, got %v", model.loading)
	}

	// Test text input configuration
	if model.textInput.Placeholder != "Enter ZIP code (e.g., 90210)" {
		t.Errorf("Expected placeholder text not set correctly")
	}
	if model.textInput.CharLimit != 10 {
		t.Errorf("Expected char limit 10, got %d", model.textInput.CharLimit)
	}
}

func TestWeatherMessageTypes(t *testing.T) {
	// Test WeatherMsg creation
	weatherData := WeatherData{
		Location:    "Test City",
		Temperature: 68.0,
		Description: "Sunny",
		Humidity:    50,
		WindSpeed:   3.5,
		Pressure:    1015.0,
		FeelsLike:   70.0,
		Visibility:  8.5,
		UVIndex:     4.0,
		Timestamp:   time.Now(),
	}

	weatherMsg := WeatherMsg{WeatherData: weatherData}
	if weatherMsg.WeatherData.Location != "Test City" {
		t.Errorf("WeatherMsg not created correctly")
	}

	// Test ErrorMsg creation
	errorMsg := ErrorMsg{Error: "Test error message"}
	if errorMsg.Error != "Test error message" {
		t.Errorf("ErrorMsg not created correctly")
	}
}

func TestOpenWeatherResponseStructure(t *testing.T) {
	// Test that the OpenWeatherResponse struct can be properly initialized
	response := OpenWeatherResponse{
		Weather: []struct {
			Main        string `json:"main"`
			Description string `json:"description"`
		}{
			{Main: "Clear", Description: "clear sky"},
		},
		Main: struct {
			Temp      float64 `json:"temp"`
			FeelsLike float64 `json:"feels_like"`
			Humidity  int     `json:"humidity"`
			Pressure  float64 `json:"pressure"`
		}{
			Temp:      295.15,
			FeelsLike: 297.15,
			Humidity:  60,
			Pressure:  1013.25,
		},
		Wind: struct {
			Speed float64 `json:"speed"`
		}{
			Speed: 3.5,
		},
		Visibility: 10000,
		Name:       "London",
		Sys: struct {
			Country string `json:"country"`
		}{
			Country: "GB",
		},
	}

	// Verify the structure is properly populated
	if len(response.Weather) != 1 {
		t.Errorf("Expected 1 weather entry, got %d", len(response.Weather))
	}
	if response.Weather[0].Main != "Clear" {
		t.Errorf("Expected weather main 'Clear', got %s", response.Weather[0].Main)
	}
	if response.Main.Temp != 295.15 {
		t.Errorf("Expected temperature 295.15, got %f", response.Main.Temp)
	}
	if response.Name != "London" {
		t.Errorf("Expected name 'London', got %s", response.Name)
	}
}

func TestIPLocationResponseStructure(t *testing.T) {
	// Test IP location response structure
	location := IPLocationResponse{
		City:        "San Francisco",
		Region:      "California",
		Country:     "US",
		CountryName: "United States",
		Latitude:    37.7749,
		Longitude:   -122.4194,
		Postal:      "94102",
	}

	if location.City != "San Francisco" {
		t.Errorf("Expected city 'San Francisco', got %s", location.City)
	}
	if location.Latitude != 37.7749 {
		t.Errorf("Expected latitude 37.7749, got %f", location.Latitude)
	}
	if location.Longitude != -122.4194 {
		t.Errorf("Expected longitude -122.4194, got %f", location.Longitude)
	}
}

func TestZipLocationResponseStructure(t *testing.T) {
	// Test ZIP location response structure
	zipLocation := ZipLocationResponse{
		Zip:     "90210",
		Name:    "Beverly Hills",
		Lat:     34.0901,
		Lon:     -118.4065,
		Country: "US",
	}

	if zipLocation.Zip != "90210" {
		t.Errorf("Expected ZIP '90210', got %s", zipLocation.Zip)
	}
	if zipLocation.Name != "Beverly Hills" {
		t.Errorf("Expected name 'Beverly Hills', got %s", zipLocation.Name)
	}
	if zipLocation.Lat != 34.0901 {
		t.Errorf("Expected latitude 34.0901, got %f", zipLocation.Lat)
	}
}

// Benchmark tests for performance
func BenchmarkGetAPIKey(b *testing.B) {
	os.Setenv("OPENWEATHER_API_KEY", "test_key")
	defer os.Unsetenv("OPENWEATHER_API_KEY")

	for i := 0; i < b.N; i++ {
		getAPIKey()
	}
}

func BenchmarkIsAPIKeyConfigured(b *testing.B) {
	os.Setenv("OPENWEATHER_API_KEY", "test_key")
	defer os.Unsetenv("OPENWEATHER_API_KEY")

	for i := 0; i < b.N; i++ {
		isAPIKeyConfigured()
	}
}

func BenchmarkInitialModel(b *testing.B) {
	for i := 0; i < b.N; i++ {
		initialModel()
	}
}

// chdirTemp switches into a fresh temp dir for the duration of the test so
// local secret files (secrets.enc.env, .env) don't leak into results.
func chdirTemp(t *testing.T) {
	t.Helper()
	orig, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	if err := os.Chdir(t.TempDir()); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = os.Chdir(orig) })
}

func TestParseDotenvValue(t *testing.T) {
	content := "# comment\n\nOTHER=x\nOPENWEATHER_API_KEY=\"quoted123\"\n"
	if got := parseDotenvValue(content, "OPENWEATHER_API_KEY"); got != "quoted123" {
		t.Errorf("got %q, want %q", got, "quoted123")
	}
	if got := parseDotenvValue("OPENWEATHER_API_KEY='single'", "OPENWEATHER_API_KEY"); got != "single" {
		t.Errorf("got %q, want %q", got, "single")
	}
	if got := parseDotenvValue("NOPE=1", "OPENWEATHER_API_KEY"); got != "" {
		t.Errorf("got %q, want empty", got)
	}
}

func TestReadAPIKeyFromSopsMissingFile(t *testing.T) {
	chdirTemp(t)
	if got := readAPIKeyFromSops(); got != "" {
		t.Errorf("expected empty key when %s is absent, got %q", sopsSecretsFile, got)
	}
}
