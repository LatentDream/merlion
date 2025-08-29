package main

import (
	"fmt"
	"merlion/internal/context"
	"merlion/internal/store/cloud"
	"merlion/internal/store/local/database"
	"merlion/internal/ui"
	"merlion/internal/utils"
	"net/http"
	_ "net/http/pprof"
	"os"
	"path/filepath"
	"slices"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/log"
)

func startTUI(flags ...string) {
	startTime := time.Now()

	// Setup Logging
	closer, err := utils.SetupLog()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer closer()
	log.Debug("Starting main function", "time", startTime)
	log.Debug("Log setup complete", "elapsed", time.Since(startTime))

	// Setup profiling in dev env
	if os.Getenv("APP_ENV") == "dev" {
		log.Warn("Enabling pprof for profiling")
		go func() {
			log.Warn(http.ListenAndServe("localhost:6060", nil))
		}()
	}

	// Initialize config directory
	userConfigDir, err := os.UserConfigDir()
	if err != nil {
		log.Fatalf("Failed to get user config directory: %v", err)
	}
	configDir := filepath.Join(userConfigDir, "merlion")
	if err := os.MkdirAll(configDir, 0700); err != nil {
		log.Fatalf("Failed to create config directory: %v", err)
	}
	log.Debug("Config directory initialized", "elapsed", time.Since(startTime))

	// Initialize theme manager
	options := []context.ContextOption{}
	if slices.Contains(flags, "--compact") {
		options = append(options, context.WithCompactViewStart(true))
	}
	if slices.Contains(flags, "--no-save") {
		options = append(options, context.WithSaveOnChange(false))
	}
	if slices.Contains(flags, "--local") {
		options = append(options, context.WithLocalFirst(true))
	} else if slices.Contains(flags, "--remote") {
		options = append(options, context.WithLocalFirst(false))
	}

	// App CTX
	ctx, err := context.NewContext(configDir, options...)
	if err != nil {
		log.Fatalf("Failed to initialize theme manager: %v", err)
	}
	log.Debug("Theme manager initialized", "elapsed", time.Since(startTime))

	// Initialize credentials manager
	credMgr, err := cloud.NewCredentialsManager()
	if err != nil {
		log.Fatalf("Failed to initialize credentials manager: %v", err)
	}
	log.Debug("Credentials manager initialized", "elapsed", time.Since(startTime))

	// Initial local DB - I don't like have the DB where, but needed to have the Close
	// TODO: A ctx where closing funcs can be registered would be great
	localDB, err := database.InitDB()
	if err != nil {
		log.Fatalf("Failed to init DB: %v", err)
	}
	defer localDB.Close()

	model, err := ui.NewModel(credMgr, localDB, ctx)
	if err != nil {
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
		log.Fatalf("Error running program: %v", err)
	}
	log.Debug("Main function completed", "total_time", time.Since(startTime))
}
