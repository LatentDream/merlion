package main

import (
	"net/http"
	_ "net/http/pprof"
	"os"
	"slices"

	"merlion/cmd/merlion/vault"
	"merlion/internal/config"
	"merlion/internal/context"
	"merlion/internal/store/cloud"
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

	// Load config
	cfg := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	if cfg.Vaults == nil || len(cfg.Vaults) == 0 {
		ok := vault.ChooseVault()
		if ok != 0 {
			os.Exit(1)
		}
	}

	model, err := ui.NewModel(cfg, credMgr, ctx)
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
