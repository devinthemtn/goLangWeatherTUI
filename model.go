package main

import (
	"fmt"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// State represents the current screen state
type State int

const (
	MenuState State = iota
	ZipInputState
	LoadingState
	WeatherDisplayState
	ErrorState
)

// Model represents the main TUI model
type Model struct {
	state        State
	textInput    textinput.Model
	weather      *WeatherData
	error        string
	loading      bool
	windowWidth  int
	windowHeight int
}

// WeatherData represents the weather information
type WeatherData struct {
	Location    string
	Temperature float64
	Description string
	Humidity    int
	WindSpeed   float64
	Pressure    float64
	FeelsLike   float64
	Visibility  float64
	UVIndex     float64
	Timestamp   time.Time
}

// initialModel creates and returns the initial model
func initialModel() Model {
	ti := textinput.New()
	ti.Placeholder = "Enter ZIP code (e.g., 90210)"
	ti.Focus()
	ti.CharLimit = 10
	ti.Width = 20

	return Model{
		state:     MenuState,
		textInput: ti,
		weather:   nil,
		error:     "",
		loading:   false,
	}
}

// Init initializes the model
func (m Model) Init() tea.Cmd {
	return textinput.Blink
}

// Update handles messages and updates the model
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.windowWidth = msg.Width
		m.windowHeight = msg.Height

	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			if m.state == WeatherDisplayState || m.state == ErrorState {
				// Return to menu instead of quitting
				m.state = MenuState
				m.error = ""
				m.weather = nil
				return m, nil
			}
			return m, tea.Quit

		case "esc":
			// Always return to menu on escape
			m.state = MenuState
			m.error = ""
			m.weather = nil
			m.textInput.SetValue("")
			return m, nil

		case "enter":
			switch m.state {
			case MenuState:
				return m, nil
			case ZipInputState:
				if strings.TrimSpace(m.textInput.Value()) != "" {
					m.state = LoadingState
					m.loading = true
					return m, m.fetchWeatherByZip(m.textInput.Value())
				}
			}

		case "1":
			if m.state == MenuState {
				m.state = LoadingState
				m.loading = true
				return m, m.fetchWeatherByLocation()
			}

		case "2":
			if m.state == MenuState {
				m.state = ZipInputState
				m.textInput.SetValue("")
				m.textInput.Focus()
				return m, textinput.Blink
			}

		case "r":
			if m.state == WeatherDisplayState {
				m.state = LoadingState
				m.loading = true
				// Refresh current weather data
				if m.weather != nil {
					return m, m.fetchWeatherByLocation()
				}
			}
		}

	case WeatherMsg:
		m.state = WeatherDisplayState
		m.weather = &msg.WeatherData
		m.loading = false
		m.error = ""

	case ErrorMsg:
		m.state = ErrorState
		m.error = msg.Error
		m.loading = false
	}

	// Update text input
	if m.state == ZipInputState {
		m.textInput, cmd = m.textInput.Update(msg)
	}

	return m, cmd
}

// View renders the current view
func (m Model) View() string {
	switch m.state {
	case MenuState:
		return m.menuView()
	case ZipInputState:
		return m.zipInputView()
	case LoadingState:
		return m.loadingView()
	case WeatherDisplayState:
		return m.weatherView()
	case ErrorState:
		return m.errorView()
	default:
		return "Unknown state"
	}
}

// Styles
var (
	titleStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FAFAFA")).
			Background(lipgloss.Color("#7D56F4")).
			Padding(0, 1).
			Bold(true)

	menuStyle = lipgloss.NewStyle().
			BorderStyle(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("#874BFD")).
			Padding(1, 2)

	weatherCardStyle = lipgloss.NewStyle().
				BorderStyle(lipgloss.RoundedBorder()).
				BorderForeground(lipgloss.Color("#04B575")).
				Padding(1, 2).
				MarginTop(1)

	errorStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FF5F87")).
			BorderStyle(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("#FF5F87")).
			Padding(1, 2)

	loadingStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FFA500")).
			Bold(true)
)

func (m Model) menuView() string {
	title := titleStyle.Render("🌤️  Weather TUI")

	menu := menuStyle.Render(
		"Welcome to Weather TUI!\n\n" +
			"Choose an option:\n\n" +
			"[1] Get weather for current location\n" +
			"[2] Enter ZIP code\n\n" +
			"[q] Quit",
	)

	return lipgloss.JoinVertical(
		lipgloss.Center,
		title,
		"",
		menu,
	)
}

func (m Model) zipInputView() string {
	title := titleStyle.Render("🌤️  Weather TUI")

	content := fmt.Sprintf(
		"Enter ZIP Code:\n\n%s\n\n[Enter] Get Weather\n[Esc] Back to Menu",
		m.textInput.View(),
	)

	input := menuStyle.Render(content)

	return lipgloss.JoinVertical(
		lipgloss.Center,
		title,
		"",
		input,
	)
}

func (m Model) loadingView() string {
	title := titleStyle.Render("🌤️  Weather TUI")

	loading := loadingStyle.Render("🔄 Loading weather data...")

	return lipgloss.JoinVertical(
		lipgloss.Center,
		title,
		"",
		loading,
	)
}

func (m Model) weatherView() string {
	if m.weather == nil {
		return "No weather data available"
	}

	title := titleStyle.Render("🌤️  Weather TUI")

	// Main weather info
	mainWeather := fmt.Sprintf(
		"📍 %s\n"+
			"🌡️  %.1f°F (Feels like %.1f°F)\n"+
			"☁️  %s\n"+
			"💧 Humidity: %d%%",
		m.weather.Location,
		m.weather.Temperature,
		m.weather.FeelsLike,
		m.weather.Description,
		m.weather.Humidity,
	)

	// Additional details
	details := fmt.Sprintf(
		"💨 Wind: %.1f mph\n"+
			"🔽 Pressure: %.1f hPa\n"+
			"👁️  Visibility: %.1f mi\n"+
			"☀️ UV Index: %.1f",
		m.weather.WindSpeed,
		m.weather.Pressure,
		m.weather.Visibility,
		m.weather.UVIndex,
	)

	timestamp := fmt.Sprintf("Last updated: %s", m.weather.Timestamp.Format("3:04 PM"))

	weatherCard := weatherCardStyle.Render(
		lipgloss.JoinVertical(
			lipgloss.Left,
			mainWeather,
			"",
			details,
			"",
			timestamp,
		),
	)

	controls := "Press [r] to refresh • [Esc] back to menu • [q] quit"

	return lipgloss.JoinVertical(
		lipgloss.Center,
		title,
		weatherCard,
		"",
		controls,
	)
}

func (m Model) errorView() string {
	title := titleStyle.Render("🌤️  Weather TUI")

	errorMsg := errorStyle.Render(
		fmt.Sprintf("❌ Error: %s\n\n[Esc] Back to Menu", m.error),
	)

	return lipgloss.JoinVertical(
		lipgloss.Center,
		title,
		"",
		errorMsg,
	)
}
