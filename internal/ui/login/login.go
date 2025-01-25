package login

import (
	"fmt"
	"strings"

	"merlion/internal/api"
	"merlion/internal/auth"
	"merlion/internal/styles"
	"merlion/internal/ui/navigation"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/log"
)

type Model struct {
	emailInput         textinput.Model
	passwordInput      textinput.Model
	err                error
	credentials        *auth.Credentials
	width              int
	height             int
	validating         bool
	credentialsManager *auth.CredentialsManager
	styles             *styles.Styles
	themeManager       *styles.ThemeManager
}

func NewModel(credentialsManager *auth.CredentialsManager, themeManager *styles.ThemeManager) navigation.View {
	appStyles := themeManager.Styles()

	emailInput := textinput.New()
	emailInput.Placeholder = "Enter your email"
	emailInput.Focus()
	emailInput.CharLimit = 64
	emailInput.Width = 32

	// Apply theme to input
	emailInput.PromptStyle = appStyles.Input
	emailInput.TextStyle = appStyles.Input

	passwordInput := textinput.New()
	passwordInput.Placeholder = "Enter your password"
	passwordInput.CharLimit = 64
	passwordInput.Width = 32
	passwordInput.EchoMode = textinput.EchoPassword
	passwordInput.EchoCharacter = '•'
	// Apply theme to input
	passwordInput.PromptStyle = appStyles.Input
	passwordInput.TextStyle = appStyles.Input

	return Model{
		emailInput:         emailInput,
		passwordInput:      passwordInput,
		styles:             appStyles,
		themeManager:       themeManager,
		credentialsManager: credentialsManager,
	}
}

func (m Model) SetClient(client *api.Client) {
	// empty
}

func (m Model) Init() tea.Cmd {
	return textinput.Blink
}

func (m Model) Update(msg tea.Msg) (navigation.View, tea.Cmd) {
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
			if m.validating {
				return m, nil // Ignore tab while validating
			}
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
			if m.validating {
				return m, nil // Ignore enter while validating
			}

			if m.passwordInput.Focused() {
				// Validate credentials
				m.validating = true
				creds := auth.Credentials{
					Email:    strings.TrimSpace(m.emailInput.Value()),
					Password: m.passwordInput.Value(),
				}

				client, err := api.NewClient(nil)
				if err != nil {
					m.err = fmt.Errorf("could not initialize client: %w", err)
					m.validating = false
					return m, nil
				}

				if err := client.ValidateCredentials(creds); err != nil {
					m.err = err
					m.validating = false
					return m, nil
				}

				err = m.credentialsManager.SaveCredentials(&creds)
				if err != nil {
					log.Fatalf("Not able to save credentials %v", err)
				}
				return m, navigation.LoginCmd(client)
			}
			// Move to password field when pressing enter in email field
			if m.emailInput.Focused() {
				m.emailInput.Blur()
				m.passwordInput.Focus()
				return m, textinput.Blink
			}

		case "ctrl+t": // Add theme toggle shortcut
			if err := m.themeManager.NextTheme(); err != nil {
				m.err = fmt.Errorf("failed to change theme: %v", err)
				return m, nil
			}
			m.styles = m.themeManager.Styles()
			// Update input styles
			m.emailInput.PromptStyle = m.styles.Input
			m.emailInput.TextStyle = m.styles.Input
			m.passwordInput.PromptStyle = m.styles.Input
			m.passwordInput.TextStyle = m.styles.Input
			return m, nil
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

func (m Model) View() string {
	// Create the form content
	var b strings.Builder

	// Use theme styles
	b.WriteString(m.styles.Title.Render("Welcome to Merlion!") + "\n\n")
	b.WriteString(m.styles.App.Render("Please enter your credentials") + "\n\n")

	inputContainer := m.styles.Input.Copy().Padding(0, 1)
	b.WriteString(inputContainer.Render("Email: "+m.emailInput.View()) + "\n")
	b.WriteString(inputContainer.Render("Password: "+m.passwordInput.View()) + "\n\n")

	if m.validating {
		b.WriteString(m.styles.Muted.Render("Validating credentials...") + "\n")
	} else {
		b.WriteString(m.styles.Muted.Render("(Press tab to switch fields, enter to submit)") + "\n")
	}

	if m.err != nil {
		b.WriteString(m.styles.Error.Render(fmt.Sprintf("Error: %v", m.err)) + "\n")
	}

	// Add signup notice
	b.WriteString(m.styles.Muted.Render("\nDon't have an account yet?") + "\n")
	b.WriteString(m.styles.Muted.Render("Visit https://merlion.dev to create one"))

	// Center the container in the terminal
	container := m.styles.Container.Render(b.String())

	return lipgloss.Place(
		m.width,
		m.height,
		lipgloss.Center,
		lipgloss.Center,
		container,
	)
}
