package main

import (
	"net/http"
	_ "net/http/pprof"
	"os"
	"slices"

	"merlion/internal/context"
	"merlion/internal/store/cloud"
	"merlion/internal/store/sqlite/database"
	"merlion/internal/ui"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/log"
)

func startTUI(flags ...string) {
	// Setup profiling in dev env
	if os.Getenv("APP_ENV") == "dev" {
		log.Warn("Enabling pprof for profiling")
		go func() {
			log.Warn(http.ListenAndServe("localhost:6060", nil))
		}()
	}

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
	if slices.Contains(flags, "--favorites") {
		options = append(options, context.WithFavorityOpen())
	}
	if slices.Contains(flags, "--work-logs") {
		options = append(options, context.WithWorkLogOpen())
	}

	// App CTX
	ctx, err := context.NewContext(options...)
	if err != nil {
		log.Fatalf("Failed to initialize theme manager: %v", err)
	}

	// Initialize credentials manager
	credMgr, err := cloud.NewCredentialsManager()
	if err != nil {
		log.Fatalf("Failed to initialize credentials manager: %v", err)
	}

	// Initial local DB - I don't like have the DB where, but needed to have the Close
	// TODO: A ctx where closing funcs can be registered would be great
	localDB, err := database.InitDB()
	if err != nil {
		log.Fatalf("Failed to init DB: %v", err)
	}
	defer localDB.Close()

	// Init the file client from the MERLION_PATH env var
	var localPath *string = nil
	if os.Getenv("MERLION_PATH") != "" {
		root := os.Getenv("MERLION_PATH")
		localPath = &root
	}

	model, err := ui.NewModel(credMgr, localDB, localPath, ctx)
	if err != nil {
		log.Fatalf("Failed to create UI model: %v", err)
	}

	p := tea.NewProgram(
		model,
		tea.WithAltScreen(),
		tea.WithMouseCellMotion(),
	)

	if _, err := p.Run(); err != nil {
		log.Fatalf("Error running program: %v", err)
	}
}
