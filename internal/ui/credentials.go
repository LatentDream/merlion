// internal/ui/credentials.go
package ui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"merlion/internal/auth"
)

type credentialsModel struct {
	emailInput    textinput.Model
	passwordInput textinput.Model
	err           error
	done          bool
	credentials   *auth.Credentials
	width         int
	height        int
}

func NewCredentialsUI() *credentialsModel {
	emailInput := textinput.New()
	emailInput.Placeholder = "Enter your email"
	emailInput.Focus()
	emailInput.CharLimit = 64
	emailInput.Width = 32

	passwordInput := textinput.New()
	passwordInput.Placeholder = "Enter your password"
	passwordInput.CharLimit = 64
	passwordInput.Width = 32
	passwordInput.EchoMode = textinput.EchoPassword
	passwordInput.EchoCharacter = '•'

	return &credentialsModel{
		emailInput:    emailInput,
		passwordInput: passwordInput,
	}
}

func (m credentialsModel) Init() tea.Cmd {
	return textinput.Blink
}

func (m credentialsModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		return m, nil

	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "esc":
			return m, tea.Quit

		case "tab", "shift+tab":
			// Cycle between inputs
			if m.emailInput.Focused() {
				m.emailInput.Blur()
				m.passwordInput.Focus()
			} else {
				m.passwordInput.Blur()
				m.emailInput.Focus()
			}
			return m, textinput.Blink

		case "enter":
			if m.passwordInput.Focused() {
				m.done = true
				m.credentials = &auth.Credentials{
					Email:    strings.TrimSpace(m.emailInput.Value()),
					Password: m.passwordInput.Value(),
				}
				return m, tea.Quit
			}
			// Move to password field when pressing enter in email field
			if m.emailInput.Focused() {
				m.emailInput.Blur()
				m.passwordInput.Focus()
				return m, textinput.Blink
			}
		}
	}

	// Handle input updates
	if m.emailInput.Focused() {
		m.emailInput, cmd = m.emailInput.Update(msg)
		return m, cmd
	}
	m.passwordInput, cmd = m.passwordInput.Update(msg)
	return m, cmd
}

var (
	titleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#FAFAFA")).
			Background(lipgloss.Color("#7D56F4")).
			Padding(0, 1)

	inputStyle = lipgloss.NewStyle().
			Padding(0, 1)

	containerStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			Padding(1, 2).
			BorderForeground(lipgloss.Color("#7D56F4"))
)

func (m credentialsModel) View() string {
	if m.width == 0 {
		return "Loading..."
	}

	// Create the form content
	var b strings.Builder
	b.WriteString(titleStyle.Render("Welcome to Merlion!") + "\n\n")
	b.WriteString("Please enter your credentials\n\n")
	b.WriteString(inputStyle.Render("Email: "+m.emailInput.View()) + "\n")
	b.WriteString(inputStyle.Render("Password: "+m.passwordInput.View()) + "\n\n")
	b.WriteString("(Press tab to switch fields, enter to submit)")

	if m.err != nil {
		b.WriteString(fmt.Sprintf("\nError: %v\n", m.err))
	}

	// Center the container in the terminal
	container := containerStyle.Render(b.String())

	return lipgloss.Place(
		m.width,
		m.height,
		lipgloss.Center,
		lipgloss.Center,
		container,
	)
}

func GetCredentials() (*auth.Credentials, error) {
	p := tea.NewProgram(
		NewCredentialsUI(),
		tea.WithAltScreen(),       // Use alternate screen
		tea.WithMouseCellMotion(), // Enable mouse support
	)

	m, err := p.Run()
	if err != nil {
		return nil, fmt.Errorf("running credentials UI: %w", err)
	}

	if m, ok := m.(credentialsModel); ok && m.done {
		return m.credentials, nil
	}

	return nil, fmt.Errorf("credentials input cancelled")
}
