package main

import (
	"fmt"
	"merlion/internal/store/cloud"
	"merlion/internal/styles"
	"merlion/internal/ui"
	"merlion/internal/utils"
	"net/http"
	_ "net/http/pprof"
	"os"
	"path/filepath"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/log"
)

func main() {
	startTime := time.Now()
	closer, err := utils.SetupLog()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	log.Debug("Starting main function", "time", startTime)
	log.Debug("Log setup complete", "elapsed", time.Since(startTime))

	if os.Getenv("APP_ENV") == "dev" {
		log.Info("Enabling pprof for profiling")
		go func() {
			log.Info(http.ListenAndServe("localhost:6060", nil))
		}()
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
	log.Debug("Config directory initialized", "elapsed", time.Since(startTime))

	// Initialize theme manager
	themeManager, err := styles.NewThemeManager(configDir)
	if err != nil {
		_ = closer()
		log.Fatalf("Failed to initialize theme manager: %v", err)
	}
	log.Debug("Theme manager initialized", "elapsed", time.Since(startTime))

	// Initialize credentials manager
	credMgr, err := cloud.NewCredentialsManager()
	if err != nil {
		_ = closer()
		log.Fatalf("Failed to initialize credentials manager: %v", err)
	}
	log.Debug("Credentials manager initialized", "elapsed", time.Since(startTime))

	model, err := ui.NewModel(credMgr, themeManager)
	if err != nil {
		_ = closer()
		log.Fatalf("Failed to create UI model: %v", err)
	}
	log.Debug("UI model created", "elapsed", time.Since(startTime))

	p := tea.NewProgram(
		model,
		tea.WithAltScreen(),
		tea.WithMouseCellMotion(),
	)
	log.Debug("Tea program created", "elapsed", time.Since(startTime))

	log.Debug("Starting Tea program", "elapsed", time.Since(startTime))
	if _, err := p.Run(); err != nil {
		_ = closer()
		log.Fatalf("Error running program: %v", err)
	}
	_ = closer()
	log.Debug("Main function completed", "total_time", time.Since(startTime))
}
