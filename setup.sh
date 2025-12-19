#!/bin/bash

# Weather TUI Setup Script
# This script helps you configure the OpenWeatherMap API key

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

# Function to validate API key format
validate_api_key() {
    local key=$1
    # OpenWeatherMap API keys are typically 32 characters long
    if [ ${#key} -lt 20 ]; then
        return 1
    fi
    # Check for common placeholder values
    if [[ "$key" == "your_api_key_here" ]] || [[ "$key" == "YOUR_API_KEY" ]] || [[ "$key" == "xxx" ]]; then
        return 1
    fi
    return 0
}

# Main setup function
main() {
    clear
    print_header "🌤️  Weather TUI - Setup Script"

    print_color $CYAN "This script will help you configure your OpenWeatherMap API key."
    echo

    # Check if API key is already configured
    if [ -n "$OPENWEATHER_API_KEY" ]; then
        print_color $GREEN "✅ API key is already set in environment variable"
        print_color $YELLOW "Current key: ${OPENWEATHER_API_KEY:0:8}...${OPENWEATHER_API_KEY: -4}"
        echo
        read -p "Do you want to update it? (y/N): " update_key
        if [[ ! "$update_key" =~ ^[Yy]$ ]]; then
            print_color $GREEN "Setup complete! Your API key is already configured."
            exit 0
        fi
    elif [ -f ".env" ] && grep -q "OPENWEATHER_API_KEY=" .env; then
        print_color $GREEN "✅ Found .env file with API key configuration"
        echo
        read -p "Do you want to update the API key in .env? (y/N): " update_key
        if [[ ! "$update_key" =~ ^[Yy]$ ]]; then
            print_color $GREEN "Setup complete! Check your .env file for the API key."
            exit 0
        fi
    fi

    print_header "🔑 API Key Setup"

    print_color $YELLOW "To use Weather TUI, you need a free OpenWeatherMap API key."
    echo
    print_color $CYAN "Steps to get your API key:"
    echo "  1. Visit: https://openweathermap.org/api"
    echo "  2. Click 'Sign Up' (or 'Sign In' if you have an account)"
    echo "  3. Create a free account"
    echo "  4. Go to 'API keys' section in your account"
    echo "  5. Copy your default API key"
    echo

    # Option to open browser
    if command -v xdg-open >/dev/null 2>&1 || command -v open >/dev/null 2>&1; then
        read -p "Open OpenWeatherMap API page in browser? (y/N): " open_browser
        if [[ "$open_browser" =~ ^[Yy]$ ]]; then
            if command -v xdg-open >/dev/null 2>&1; then
                xdg-open "https://openweathermap.org/api" 2>/dev/null || true
            elif command -v open >/dev/null 2>&1; then
                open "https://openweathermap.org/api" 2>/dev/null || true
            fi
            print_color $GREEN "✅ Browser opened. Come back here when you have your API key."
            echo
        fi
    fi

    # Get API key from user
    echo "Please enter your OpenWeatherMap API key:"
    read -p "API Key: " api_key

    # Trim whitespace
    api_key=$(echo "$api_key" | tr -d '[:space:]')

    # Validate API key
    if ! validate_api_key "$api_key"; then
        print_color $RED "❌ Invalid API key format."
        print_color $YELLOW "API keys should be at least 20 characters long and not contain placeholder text."
        echo
        read -p "Continue anyway? (y/N): " continue_anyway
        if [[ ! "$continue_anyway" =~ ^[Yy]$ ]]; then
            print_color $RED "Setup cancelled."
            exit 1
        fi
    fi

    print_header "💾 Configuration Method"

    echo "How would you like to store your API key?"
    echo "1) Environment variable (recommended for permanent setup)"
    echo "2) .env file (recommended for development)"
    echo "3) Both (most flexible)"
    echo
    read -p "Choose option (1-3): " config_method

    case $config_method in
        1)
            setup_environment_variable "$api_key"
            ;;
        2)
            setup_env_file "$api_key"
            ;;
        3)
            setup_env_file "$api_key"
            setup_environment_variable "$api_key"
            ;;
        *)
            print_color $RED "Invalid option. Setting up .env file as default."
            setup_env_file "$api_key"
            ;;
    esac

    print_header "✅ Setup Complete!"

    print_color $GREEN "Your API key has been configured successfully!"
    echo

    # Test the configuration
    print_color $YELLOW "Testing configuration..."
    if ./goWeatherTUI --test-api 2>/dev/null; then
        print_color $GREEN "✅ API test successful!"
    else
        print_color $YELLOW "⚠️  API test failed, but this might be due to network issues."
        print_color $CYAN "Try running the app to see if it works: ./goWeatherTUI"
    fi

    echo
    print_color $CYAN "You can now run the Weather TUI application:"
    print_color $GREEN "  ./goWeatherTUI"
    print_color $CYAN "  or: go run ."
    echo
    print_color $CYAN "For debugging: ./goWeatherTUI --debug"
}

setup_environment_variable() {
    local api_key=$1

    print_color $YELLOW "Setting up environment variable..."

    # Detect shell
    shell_name=$(basename "$SHELL")

    case $shell_name in
        bash)
            shell_file="$HOME/.bashrc"
            ;;
        zsh)
            shell_file="$HOME/.zshrc"
            ;;
        fish)
            shell_file="$HOME/.config/fish/config.fish"
            ;;
        *)
            shell_file="$HOME/.profile"
            ;;
    esac

    # Add export line
    export_line="export OPENWEATHER_API_KEY=\"$api_key\""

    if [ -f "$shell_file" ]; then
        # Check if already exists
        if grep -q "OPENWEATHER_API_KEY" "$shell_file"; then
            print_color $YELLOW "Updating existing API key in $shell_file"
            # Use sed to replace the line
            if [[ "$OSTYPE" == "darwin"* ]]; then
                # macOS
                sed -i '' "s/.*OPENWEATHER_API_KEY.*/$export_line/" "$shell_file"
            else
                # Linux
                sed -i "s/.*OPENWEATHER_API_KEY.*/$export_line/" "$shell_file"
            fi
        else
            print_color $YELLOW "Adding API key to $shell_file"
            echo "" >> "$shell_file"
            echo "# OpenWeatherMap API Key for Weather TUI" >> "$shell_file"
            echo "$export_line" >> "$shell_file"
        fi
    else
        print_color $YELLOW "Creating $shell_file"
        echo "# OpenWeatherMap API Key for Weather TUI" > "$shell_file"
        echo "$export_line" >> "$shell_file"
    fi

    # Set for current session
    export OPENWEATHER_API_KEY="$api_key"

    print_color $GREEN "✅ Environment variable configured in $shell_file"
    print_color $YELLOW "Note: Restart your terminal or run 'source $shell_file' to use in new sessions"
}

setup_env_file() {
    local api_key=$1

    print_color $YELLOW "Setting up .env file..."

    # Create or update .env file
    if [ -f ".env" ]; then
        if grep -q "OPENWEATHER_API_KEY" .env; then
            print_color $YELLOW "Updating existing API key in .env"
            # Use sed to replace the line
            if [[ "$OSTYPE" == "darwin"* ]]; then
                # macOS
                sed -i '' "s/.*OPENWEATHER_API_KEY.*/OPENWEATHER_API_KEY=$api_key/" .env
            else
                # Linux
                sed -i "s/.*OPENWEATHER_API_KEY.*/OPENWEATHER_API_KEY=$api_key/" .env
            fi
        else
            print_color $YELLOW "Adding API key to existing .env"
            echo "" >> .env
            echo "# OpenWeatherMap API Key" >> .env
            echo "OPENWEATHER_API_KEY=$api_key" >> .env
        fi
    else
        print_color $YELLOW "Creating new .env file"
        cat > .env << EOF
# Weather TUI Configuration
# OpenWeatherMap API Key - Get yours at: https://openweathermap.org/api
OPENWEATHER_API_KEY=$api_key

# Optional: Default settings
# DEFAULT_UNITS=imperial
# DEFAULT_COUNTRY=US
EOF
    fi

    print_color $GREEN "✅ .env file configured"
}

# Check prerequisites
check_prerequisites() {
    if ! command -v go >/dev/null 2>&1; then
        print_color $RED "❌ Go is not installed"
        print_color $YELLOW "Please install Go from: https://golang.org/"
        exit 1
    fi

    if [ ! -f "goWeatherTUI" ] && [ ! -f "weather-tui" ]; then
        print_color $YELLOW "⚠️  Application not built yet. Building now..."
        if go build -o goWeatherTUI; then
            print_color $GREEN "✅ Application built successfully"
        else
            print_color $RED "❌ Failed to build application"
            exit 1
        fi
    fi
}

# Cleanup function
cleanup() {
    print_color $YELLOW "Setup interrupted."
    exit 0
}

# Set up signal handlers
trap cleanup SIGINT SIGTERM

# Check prerequisites first
check_prerequisites

# Run main function
main "$@"
