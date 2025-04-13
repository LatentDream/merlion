package main

import (
	"fmt"
	"merlion/internal/store/cloud"
	"merlion/internal/styles"
	"merlion/internal/ui"
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

	// Initialize credentials manager
	credMgr, err := cloud.NewCredentialsManager()
	if err != nil {
		_ = closer()
		log.Fatalf("Failed to initialize credentials manager: %v", err)
	}

	model, err := ui.NewModel(credMgr, themeManager)
	if err != nil {
		_ = closer()
		log.Fatalf("Failed to create UI model: %v", err)
	}

	p := tea.NewProgram(
		model,
		tea.WithAltScreen(),
		tea.WithMouseCellMotion(),
	)

	if _, err := p.Run(); err != nil {
		_ = closer()
		log.Fatalf("Error running program: %v", err)
	}
	_ = closer()
}
