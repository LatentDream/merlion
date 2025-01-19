package main

import (
	"fmt"
	"merlion/internal/api"
	"merlion/internal/auth"
	"merlion/internal/styles"
	"merlion/internal/ui"
	NotesUI "merlion/internal/ui/notes"
	"merlion/internal/utils"
	"os"
	"path/filepath"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/log"
)

func main() {
	closer, err := utils.SetupLog()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// Initialize config directory
	userConfigDir, err := os.UserConfigDir()
	if err != nil {
		_ = closer()
		log.Fatalf("Failed to get user config directory: %v", err)
	}

	configDir := filepath.Join(userConfigDir, "merlion")
	if err := os.MkdirAll(configDir, 0700); err != nil {
		log.Fatalf("Failed to create config directory: %v", err)
	}

	// Initialize theme manager
	themeManager, err := styles.NewThemeManager(configDir)
	if err != nil {
		_ = closer()
		log.Fatalf("Failed to initialize theme manager: %v", err)
	}

	// Get current styles
	appStyles := themeManager.Styles()

	// Initialize credentials manager
	credMgr, err := auth.NewCredentialsManager()
	if err != nil {
		_ = closer()
		log.Fatalf("Failed to initialize credentials manager: %v", err)
	}

	// Try to load credentials
	creds, err := credMgr.LoadCredentials()
	if err != nil {
		// First time setup - prompt user for credentials
		creds, err = ui.GetCredentials(appStyles, themeManager)
		if err != nil {
			_ = closer()
			log.Fatalf("Failed to get credentials: %v", err)
		}

		// Save the credentials
		if err := credMgr.SaveCredentials(*creds); err != nil {
			_ = closer()
			log.Fatalf("Failed to save credentials: %v", err)
		}
	}

	// Initialize API client
	client, err := api.NewClient(creds)
	if err != nil {
		_ = closer()
		log.Fatalf("Failed to create API client: %v", err)
	}

	// Login
	if err := client.Login(); err != nil {
		_ = closer()
		log.Fatalf("Failed to login: %v", err)
	}

	// Create a channel for notes
	notesChan := make(chan []api.Note)

	model, err := NotesUI.NewModel([]api.Note{}, client, themeManager)
	if err != nil {
		_ = closer()
		log.Fatalf("Failed to create UI model: %v", err)
	}

	// Start loading notes in background
	go func() {
		notes, err := client.ListNotes()
		if err != nil {
			log.Printf("Failed to list notes: %v", err)
			notesChan <- []api.Note{}
			return
		}
		notesChan <- notes
	}()

	p := tea.NewProgram(
		model,
		tea.WithAltScreen(),       // Use alternate screen buffer
		tea.WithMouseCellMotion(), // Enable mouse support
	)

	// Start note loading
	go func() {
		notes := <-notesChan
		p.Send(NotesUI.NotesLoadedMsg{Notes: notes})
	}()

	if _, err := p.Run(); err != nil {
		_ = closer()
		log.Fatal("Error running merlion: %v", err)
	}
	_ = closer()
}
