package main

import (
	"fmt"
	tea "github.com/charmbracelet/bubbletea"
	"log"
	"merlion/internal/api"
	"merlion/internal/auth"
	"merlion/internal/styles"
	"merlion/internal/ui"
	"os"
	"path/filepath"
)

func main() {
	// Initialize config directory
	userConfigDir, err := os.UserConfigDir()
	if err != nil {
		log.Fatalf("Failed to get user config directory: %v", err)
	}

	configDir := filepath.Join(userConfigDir, "merlion")
	if err := os.MkdirAll(configDir, 0700); err != nil {
		log.Fatalf("Failed to create config directory: %v", err)
	}

	// Initialize theme manager
	themeManager, err := styles.NewThemeManager(configDir)
	if err != nil {
		log.Fatalf("Failed to initialize theme manager: %v", err)
	}

	// Get current styles
	appStyles := themeManager.Styles()

	// Initialize credentials manager
	credMgr, err := auth.NewCredentialsManager()
	if err != nil {
		log.Fatalf("Failed to initialize credentials manager: %v", err)
	}

	// Try to load credentials
	creds, err := credMgr.LoadCredentials()
	if err != nil {
		// First time setup - prompt user for credentials
		creds, err = ui.GetCredentials(appStyles, themeManager)
		if err != nil {
			log.Fatalf("Failed to get credentials: %v", err)
		}

		// Save the credentials
		if err := credMgr.SaveCredentials(*creds); err != nil {
			log.Fatalf("Failed to save credentials: %v", err)
		}
	}

	// Initialize API client
	client, err := api.NewClient(creds)
	if err != nil {
		log.Fatalf("Failed to create API client: %v", err)
	}

	// Login
	if err := client.Login(); err != nil {
		log.Fatalf("Failed to login: %v", err)
	}

	// Now you can make authenticated requests
	notes, err := client.ListNotes()
	if err != nil {
		log.Fatalf("Failed to list notes: %v", err)
	}

	p := tea.NewProgram(
		initialModel(notes),
		tea.WithAltScreen(),       // Use alternate screen buffer
		tea.WithMouseCellMotion(), // Enable mouse support
	)

	if _, err := p.Run(); err != nil {
		fmt.Printf("Error running merlion: %v", err)
		os.Exit(1)
	}
}

type model struct {
	content string
	width   int
	height  int
}

func initialModel(notes []api.Note) model {
	return model{
		content: fmt.Sprintf("Welcome to Merlion - Found %d notes", len(notes)),
	}
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit
		}

	// Handle window size changes
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
	}

	return m, nil
}

func (m model) View() string {
	if m.width == 0 {
		return "Loading..."
	}

	return fmt.Sprintf(
		"Window size: %d x %d\n\n%s\n\n(press q to quit)",
		m.width, m.height, m.content,
	)
}
