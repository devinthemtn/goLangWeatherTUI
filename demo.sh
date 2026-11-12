#!/bin/bash

# Weather TUI Demo Script
# This script demonstrates the Weather TUI application

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
PURPLE='\033[0;35m'
CYAN='\033[0;36m'
NC='\033[0m' # No Color

# Function to print colored output
print_color() {
    color=$1
    message=$2
    echo -e "${color}${message}${NC}"
}

# Function to print section header
print_header() {
    echo
    print_color $PURPLE "==============================================="
    print_color $PURPLE "$1"
    print_color $PURPLE "==============================================="
    echo
}

# Function to check if command exists
command_exists() {
    command -v "$1" >/dev/null 2>&1
}

# Main demo function
main() {
    clear

    print_header "🌤️  Weather TUI Application Demo"

    print_color $CYAN "Welcome to the Weather TUI demo!"
    print_color $CYAN "This script will guide you through setting up and running the application."
    echo

    # Check prerequisites
    print_header "📋 Checking Prerequisites"

    if command_exists go; then
        GO_VERSION=$(go version | cut -d' ' -f3)
        print_color $GREEN "✓ Go is installed: $GO_VERSION"
    else
        print_color $RED "✗ Go is not installed"
        print_color $YELLOW "Please install Go from: https://golang.org/"
        exit 1
    fi

    if [ -f "go.mod" ]; then
        print_color $GREEN "✓ Go module found"
    else
        print_color $RED "✗ Go module not found"
        print_color $YELLOW "Please run this script from the goWeatherTUI directory"
        exit 1
    fi

    # Check API key
    print_header "🔑 API Key Configuration"

    if [ -n "$OPENWEATHER_API_KEY" ]; then
        print_color $GREEN "✓ Environment variable OPENWEATHER_API_KEY is set"
        API_KEY_STATUS="configured"
    elif [ -f ".env" ]; then
        if grep -q "OPENWEATHER_API_KEY=.*[^your_api_key_here]" .env 2>/dev/null; then
            print_color $GREEN "✓ API key found in .env file"
            API_KEY_STATUS="configured"
        else
            print_color $YELLOW "⚠ .env file exists but API key may not be configured"
            API_KEY_STATUS="needs_setup"
        fi
    else
        print_color $YELLOW "⚠ No API key configuration found"
        API_KEY_STATUS="needs_setup"
    fi

    if [ "$API_KEY_STATUS" = "needs_setup" ]; then
        echo
        print_color $YELLOW "To get weather data, you'll need an OpenWeatherMap API key:"
        print_color $CYAN "1. Visit: https://openweathermap.org/api"
        print_color $CYAN "2. Sign up for a free account"
        print_color $CYAN "3. Get your API key"
        print_color $CYAN "4. Set it as: export OPENWEATHER_API_KEY=your_key_here"
        print_color $CYAN "   OR create a .env file with: OPENWEATHER_API_KEY=your_key_here"
        echo
        read -p "Press Enter to continue with demo (limited functionality without API key)..."
    fi

    # Build the application
    print_header "🔨 Building the Application"

    print_color $YELLOW "Installing dependencies..."
    if go mod tidy; then
        print_color $GREEN "✓ Dependencies installed"
    else
        print_color $RED "✗ Failed to install dependencies"
        exit 1
    fi

    print_color $YELLOW "Building application..."
    if go build -o weather-tui; then
        print_color $GREEN "✓ Application built successfully"
    else
        print_color $RED "✗ Build failed"
        exit 1
    fi

    # Show application features
    print_header "✨ Application Features"

    print_color $CYAN "The Weather TUI application provides:"
    echo "  🌍 Current location weather (via IP geolocation)"
    echo "  📮 ZIP code weather lookup"
    echo "  🎨 Beautiful terminal interface"
    echo "  🔄 Real-time weather updates"
    echo "  ⌨️  Intuitive keyboard controls"
    echo

    print_color $CYAN "Keyboard Controls:"
    echo "  [1] - Get weather for current location"
    echo "  [2] - Enter ZIP code"
    echo "  [r] - Refresh weather data"
    echo "  [Esc] - Return to menu"
    echo "  [q] - Quit application"
    echo

    # Demo scenarios
    print_header "🎬 Demo Scenarios"

    echo "Choose a demo scenario:"
    echo "1) Run the full interactive application"
    echo "2) Show application help"
    echo "3) Display project structure"
    echo "4) Run development mode info"
    echo "5) Exit demo"
    echo

    read -p "Enter your choice (1-5): " choice

    case $choice in
        1)
            print_header "🚀 Starting Weather TUI Application"
            print_color $GREEN "Launching the application..."
            print_color $YELLOW "Use the keyboard controls mentioned above to navigate."
            print_color $YELLOW "Press Ctrl+C to exit when done."
            echo
            sleep 2
            ./weather-tui
            ;;
        2)
            print_header "📚 Application Help"
            if command_exists task; then
                task --list
            else
                print_color $CYAN "Available commands:"
                echo "  go run .              # Run the application"
                echo "  go build              # Build the application"
                echo "  go test ./...         # Run tests"
                echo "  go mod tidy           # Update dependencies"
            fi
            ;;
        3)
            print_header "📁 Project Structure"
            print_color $CYAN "Project files:"
            if command_exists tree; then
                tree -I 'bin|.git'
            else
                find . -type f -name "*.go" -o -name "*.md" -o -name "Taskfile.yml" -o -name ".env*" | grep -v "/\." | sort
            fi
            ;;
        4)
            print_header "🔧 Development Information"
            print_color $CYAN "Development commands:"
            echo
            if command_exists task; then
                print_color $GREEN "Using Taskfile (https://taskfile.dev):"
                echo "  task setup    # Setup development environment"
                echo "  task run      # Run the application"
                echo "  task build    # Build the application"
                echo "  task dev      # Development mode (auto-reload)"
                echo "  task test     # Run tests"
                echo "  task clean    # Clean build artifacts"
            else
                print_color $GREEN "Direct Go commands:"
                echo "  go run .              # Run directly"
                echo "  go build -o weather-tui  # Build binary"
                echo "  go test ./...         # Run tests"
            fi

            echo
            print_color $CYAN "Environment setup:"
            echo "  export OPENWEATHER_API_KEY=your_key"
            echo "  # OR create .env file with API key"
            ;;
        5)
            print_color $GREEN "Thanks for trying the Weather TUI demo!"
            exit 0
            ;;
        *)
            print_color $RED "Invalid choice. Exiting."
            exit 1
            ;;
    esac

    echo
    print_color $GREEN "Demo completed!"
    print_color $CYAN "To run the application again: ./weather-tui"
    print_color $CYAN "To get help: task --list (if task is installed)"
}

# Cleanup function
cleanup() {
    print_color $YELLOW "Demo interrupted. Cleaning up..."
    exit 0
}

# Set up signal handlers
trap cleanup SIGINT SIGTERM

# Run main function
main "$@"
