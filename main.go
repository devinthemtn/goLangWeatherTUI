package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	// Parse command line flags
	var debugMode = flag.Bool("debug", false, "Run in debug mode to check API configuration")
	var testAPI = flag.Bool("test-api", false, "Test API connectivity")
	flag.Parse()

	if *debugMode || *testAPI {
		fmt.Println("🌤️  Weather TUI - Debug Mode")
		fmt.Println("==============================\n")

		// Show environment debug info
		debugEnvironment()

		// Show API key debug info
		debugAPIKey()

		// Test API if requested
		if *testAPI {
			apiKey := getAPIKey()
			testAPIConnectivity(apiKey)
		}

		fmt.Println("\n" + strings.Repeat("=", 50))
		fmt.Println("Debug complete. To run the app normally:")
		fmt.Println("  go run . (without flags)")
		fmt.Println("  or ./weather-tui")
		return
	}

	// Normal application flow
	runNormalApp()
}

func runNormalApp() {
	// Create the TUI model
	m := initialModel()

	// Create the Bubble Tea program
	p := tea.NewProgram(m, tea.WithAltScreen())

	// Run the program
	if _, err := p.Run(); err != nil {
		log.Fatal(err)
		fmt.Printf("Error running program: %v", err)
		os.Exit(1)
	}
}
