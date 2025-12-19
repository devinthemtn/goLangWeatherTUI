package main

import (
	"fmt"
	"os"
	"strings"
)

// Debug function to check API key configuration
func debugAPIKey() {
	fmt.Println("🔍 API Key Debug Information:")
	fmt.Println("============================")

	// Check environment variable
	envKey := os.Getenv("OPENWEATHER_API_KEY")
	fmt.Printf("Environment Variable OPENWEATHER_API_KEY: ")
	if envKey == "" {
		fmt.Println("❌ NOT SET")
	} else {
		// Show first and last 4 characters for security
		if len(envKey) > 8 {
			masked := envKey[:4] + strings.Repeat("*", len(envKey)-8) + envKey[len(envKey)-4:]
			fmt.Printf("✅ SET (%s)\n", masked)
		} else {
			fmt.Printf("✅ SET (length: %d)\n", len(envKey))
		}
	}

	// Check what getAPIKey() returns
	apiKey := getAPIKey()
	fmt.Printf("getAPIKey() returns: ")
	if apiKey == "" {
		fmt.Println("❌ EMPTY")
	} else if apiKey == "your_api_key_here" {
		fmt.Println("❌ DEFAULT PLACEHOLDER")
	} else {
		if len(apiKey) > 8 {
			masked := apiKey[:4] + strings.Repeat("*", len(apiKey)-8) + apiKey[len(apiKey)-4:]
			fmt.Printf("✅ VALID (%s)\n", masked)
		} else {
			fmt.Printf("✅ SET (length: %d)\n", len(apiKey))
		}
	}

	// Check if configured according to our function
	fmt.Printf("isAPIKeyConfigured(): ")
	if isAPIKeyConfigured() {
		fmt.Println("✅ TRUE")
	} else {
		fmt.Println("❌ FALSE")
	}

	// Check all environment variables starting with OPENWEATHER
	fmt.Println("\nAll OPENWEATHER environment variables:")
	found := false
	for _, env := range os.Environ() {
		if strings.HasPrefix(env, "OPENWEATHER") {
			parts := strings.SplitN(env, "=", 2)
			if len(parts) == 2 {
				key := parts[0]
				value := parts[1]
				if len(value) > 8 {
					masked := value[:4] + strings.Repeat("*", len(value)-8) + value[len(value)-4:]
					fmt.Printf("  %s=%s\n", key, masked)
				} else {
					fmt.Printf("  %s=(length: %d)\n", key, len(value))
				}
				found = true
			}
		}
	}
	if !found {
		fmt.Println("  ❌ No OPENWEATHER variables found")
	}

	fmt.Println("\n💡 Troubleshooting Tips:")
	fmt.Println("1. Set environment variable: export OPENWEATHER_API_KEY=your_key")
	fmt.Println("2. Or create .env file with: OPENWEATHER_API_KEY=your_key")
	fmt.Println("3. Get free API key at: https://openweathermap.org/api")
	fmt.Println("4. Make sure there are no extra spaces in your key")
	fmt.Println("5. Check if your terminal session has the variable set")
}

// Test API connectivity with debug info
func testAPIConnectivity(apiKey string) {
	fmt.Println("\n🌐 Testing API Connectivity:")
	fmt.Println("============================")

	if apiKey == "" || apiKey == "your_api_key_here" {
		fmt.Println("❌ Cannot test - no valid API key provided")
		return
	}

	fmt.Println("Testing with OpenWeatherMap API...")

	// Test with a simple coordinates request (London)
	lat, lon := 51.5074, -0.1278
	weather, err := getWeatherByCoordinates(lat, lon)

	if err != nil {
		fmt.Printf("❌ API Test Failed: %v\n", err)

		// Provide specific error guidance
		errStr := err.Error()
		if strings.Contains(errStr, "401") {
			fmt.Println("💡 Error 401: Invalid API key")
			fmt.Println("   - Check if your API key is correct")
			fmt.Println("   - Make sure your API key is activated (can take a few minutes)")
		} else if strings.Contains(errStr, "429") {
			fmt.Println("💡 Error 429: API rate limit exceeded")
			fmt.Println("   - Wait a few minutes and try again")
		} else if strings.Contains(errStr, "timeout") {
			fmt.Println("💡 Network timeout")
			fmt.Println("   - Check your internet connection")
		}
	} else {
		fmt.Printf("✅ API Test Successful!\n")
		fmt.Printf("   Location: %s\n", weather.Location)
		fmt.Printf("   Temperature: %.1f°F\n", weather.Temperature)
		fmt.Printf("   Description: %s\n", weather.Description)
	}
}

// Show current working directory and file info
func debugEnvironment() {
	fmt.Println("\n📁 Environment Debug:")
	fmt.Println("====================")

	// Current working directory
	if cwd, err := os.Getwd(); err == nil {
		fmt.Printf("Current Directory: %s\n", cwd)
	}

	// Check for .env file
	if _, err := os.Stat(".env"); err == nil {
		fmt.Println(".env file: ✅ EXISTS")
		// Try to read API key from .env file
		envFileKey := readAPIKeyFromEnvFile()
		if envFileKey != "" && envFileKey != "your_api_key_here" {
			if len(envFileKey) > 8 {
				masked := envFileKey[:4] + strings.Repeat("*", len(envFileKey)-8) + envFileKey[len(envFileKey)-4:]
				fmt.Printf("  - API key found: %s\n", masked)
			} else {
				fmt.Printf("  - API key found (length: %d)\n", len(envFileKey))
			}
		} else if envFileKey == "your_api_key_here" {
			fmt.Println("  - Contains placeholder API key")
		} else {
			fmt.Println("  - No valid API key found in .env")
		}
	} else {
		fmt.Println(".env file: ❌ NOT FOUND")
	}

	// Check shell
	shell := os.Getenv("SHELL")
	if shell != "" {
		fmt.Printf("Shell: %s\n", shell)
	}
}
