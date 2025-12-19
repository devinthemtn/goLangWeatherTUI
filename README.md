# Weather TUI 🌤️

A beautiful Terminal User Interface (TUI) application written in Go that displays weather information for your current location or any ZIP code.

## Features

- 🌍 **Current Location Weather**: Automatically detects your location using IP geolocation
- 📮 **ZIP Code Search**: Get weather for any US ZIP code
- 🎨 **Beautiful TUI**: Clean, interactive terminal interface built with Bubble Tea
- 🔄 **Real-time Updates**: Refresh weather data with a single keypress
- 📊 **Comprehensive Data**: Temperature, humidity, wind speed, pressure, visibility, and more
- ⌨️ **Keyboard Navigation**: Intuitive keyboard controls

## Screenshots

```
┌─────────────────────────────┐
│  🌤️  Weather TUI            │
└─────────────────────────────┘

┌─────────────────────────────────────┐
│                                     │
│  Welcome to Weather TUI!            │
│                                     │
│  Choose an option:                  │
│                                     │
│  [1] Get weather for current location │
│  [2] Enter ZIP code                 │
│                                     │
│  [q] Quit                          │
│                                     │
└─────────────────────────────────────┘
```

## Prerequisites

- Go 1.21 or later
- OpenWeatherMap API key (free)

## Installation

1. **Clone the repository:**
   ```bash
   git clone <repository-url>
   cd goWeaatherTUI
   ```

2. **Get an OpenWeatherMap API key:**
   - Visit [OpenWeatherMap](https://openweathermap.org/api)
   - Sign up for a free account
   - Navigate to your API keys section
   - Copy your API key

3. **Configure the API key:**
   - Open `weather.go`
   - Find the line: `apiKey = "your_api_key_here"`
   - Replace `"your_api_key_here"` with your actual API key

4. **Install dependencies:**
   ```bash
   go mod tidy
   ```

5. **Build and run:**
   ```bash
   go build -o weather-tui
   ./weather-tui
   ```

   Or run directly:
   ```bash
   go run .
   ```

## Usage

### Main Menu
- Press `1` to get weather for your current location
- Press `2` to enter a ZIP code
- Press `q` to quit

### ZIP Code Input
- Enter a 5-digit US ZIP code
- Press `Enter` to fetch weather data
- Press `Esc` to return to the main menu

### Weather Display
- Press `r` to refresh the current weather data
- Press `Esc` to return to the main menu
- Press `q` to quit

## Keyboard Controls

| Key | Action |
|-----|--------|
| `1` | Get weather for current location (from main menu) |
| `2` | Enter ZIP code (from main menu) |
| `Enter` | Submit ZIP code / Confirm selection |
| `r` | Refresh weather data (from weather display) |
| `Esc` | Return to main menu |
| `q` | Quit application |
| `Ctrl+C` | Force quit |

## Weather Data Displayed

- 📍 **Location**: City, State/Country
- 🌡️ **Temperature**: Current temperature and "feels like"
- ☁️ **Conditions**: Weather description (e.g., "Clear sky", "Light rain")
- 💧 **Humidity**: Relative humidity percentage
- 💨 **Wind**: Wind speed in mph
- 🔽 **Pressure**: Atmospheric pressure in hPa
- 👁️ **Visibility**: Visibility distance in miles
- ⏰ **Last Updated**: Timestamp of data retrieval

## Dependencies

This project uses the following Go libraries:

- [Bubble Tea](https://github.com/charmbracelet/bubbletea) - TUI framework
- [Lip Gloss](https://github.com/charmbracelet/lipgloss) - Styling and layout
- [Bubbles](https://github.com/charmbracelet/bubbles) - TUI components

## API Services

- **OpenWeatherMap API**: Weather data and ZIP code geocoding
- **ipapi.co**: IP-based geolocation (no API key required)

## Configuration

### API Key Setup

The application requires an OpenWeatherMap API key. To set it up:

1. Open `weather.go`
2. Locate the `apiKey` constant
3. Replace the placeholder with your actual API key:

```go
const apiKey = "your_actual_api_key_here"
```

### Customization Options

You can customize various aspects of the application:

- **Units**: Change temperature units by modifying the `units` parameter in API calls
- **Styling**: Modify colors and styles in the `lipgloss.NewStyle()` definitions
- **Data Display**: Add or remove weather fields in the `WeatherData` struct

## Error Handling

The application handles various error scenarios:

- **No API Key**: Displays a helpful message to configure the API key
- **Network Issues**: Shows connection error messages
- **Invalid ZIP Codes**: Provides feedback for invalid or not found ZIP codes
- **API Limits**: Handles API rate limiting and quota errors

## Troubleshooting

### Common Issues

1. **"Please set your OpenWeatherMap API key"**
   - Solution: Follow the API key configuration steps above

2. **"Failed to get location"**
   - Solution: Check your internet connection
   - Alternative: Use ZIP code option instead

3. **"ZIP code not found"**
   - Solution: Ensure you're using a valid US ZIP code
   - Try a different ZIP code

4. **Build errors**
   - Solution: Run `go mod tidy` to resolve dependencies
   - Ensure you're using Go 1.21 or later

### Debug Mode

To see more detailed error information, you can modify the error messages in `weather.go` to include more context.

## Contributing

1. Fork the repository
2. Create a feature branch: `git checkout -b feature-name`
3. Make your changes
4. Run tests: `go test ./...`
5. Commit your changes: `git commit -am 'Add feature'`
6. Push to the branch: `git push origin feature-name`
7. Submit a pull request

## License

This project is licensed under the MIT License - see the LICENSE file for details.

## Acknowledgments

- [Charm](https://charm.sh/) for the amazing TUI libraries
- [OpenWeatherMap](https://openweathermap.org/) for the weather API
- [ipapi.co](https://ipapi.co/) for IP geolocation services

## Future Enhancements

- [ ] Support for international postal codes
- [ ] Weather forecast (5-day outlook)
- [ ] Weather alerts and warnings
- [ ] Favorite locations
- [ ] Configuration file support
- [ ] Weather history
- [ ] Multiple unit systems (metric/imperial)
- [ ] Weather maps and radar