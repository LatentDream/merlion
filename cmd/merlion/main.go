package main

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/charmbracelet/log"
	_ "github.com/mattn/go-sqlite3" // SQLite driver
	version "merlion/cmd/merlion/version"
	"merlion/internal/database" // Added
	"merlion/internal/store"
	"merlion/internal/store/local"
	// "merlion/internal/store/local/operations" // Removed if no longer directly used
	_ "net/http/pprof"
)

// var db *sql.DB // Removed global db variable

// getDBPath() and initDB() functions are now moved to internal/database/sqlite.go

type Command struct {
	name        string
	description string
	run         func(args ...string) int
}

var COMMANDS []Command

func helpCmd(args ...string) int {
	fmt.Println("Merlion - Help")
	fmt.Println("")
	fmt.Println("TUI Note-taking application")
	fmt.Println("Run `merlion` without args to start in TUI mode")
	fmt.Println("For remote storage capability, create an account at:")
	fmt.Println("    https://merlion.dev/login")
	fmt.Println("\nAvailable commands:")
	for _, cmd := range COMMANDS {
		fmt.Printf("  %-10s %s\n", cmd.name, cmd.description)
	}
	return 0
}

func init() {
	COMMANDS = []Command{
		{
			name:        "help",
			description: "Show help information",
			run:         helpCmd,
		},
		{
			name:        "version",
			description: "Show version information",
			run:         version.VersionCmd,
		},
		{
			name:        "logout",
			description: "Removed the cached credentials",
			run:         logoutCmd,
		},
	}
}

func main() {
	// Call the new InitDB function from the database package
	db, err := database.InitDB()
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	// Defer close on the db instance returned by InitDB
	if db != nil { // Should always be non-nil if err is nil, but good practice
		defer func() {
			if errClose := db.Close(); errClose != nil {
				log.Errorf("Failed to close database: %v", errClose)
			} else {
				log.Info("Database connection closed.")
			}
		}()
	}

	if len(os.Args) > 1 {
		command := os.Args[1]
		args := os.Args[2:]
		for _, cmd := range COMMANDS {
			if strings.ToLower(cmd.name) == strings.ToLower(command) {
				os.Exit(cmd.run(args...))
			}
		}
		fmt.Printf("Unknown command: %s\n", command)
		helpCmd()
		os.Exit(1)
	}

	localStore := local.NewLocalClient(db, "local-sqlite")
	storeManager := store.NewManager(localStore)

	startTUI(storeManager)
}
