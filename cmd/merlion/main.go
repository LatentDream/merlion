package main

import (
	"fmt"
	_ "net/http/pprof"
	"os"
	"strings"

	"merlion/cmd/merlion/export"
	"merlion/cmd/merlion/logout"
	"merlion/cmd/merlion/parser"
	"merlion/cmd/merlion/vault"
	version "merlion/cmd/merlion/version"
	"merlion/internal/utils"
)

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
	fmt.Println("    https://note.merlion.dev/login")
	fmt.Println("\nAvailable commands:")
	for _, cmd := range COMMANDS {
		fmt.Printf("  %-10s %s\n", cmd.name, cmd.description)
	}
	fmt.Println("\nAvailable flags for TUI:")
	fmt.Println("  --compact         Start in compact mode")
	fmt.Println("  --no-save         Disable autosave of config")
	fmt.Println("  --local           Show local notes when opening the app")
	fmt.Println("  --remote          Show remote notes when opening the app")
	fmt.Println("  --favorites       Show favorites tab when opening the app")
	fmt.Println("  --work-logs       Show work logs tab when opening the app")
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
			name:        "vault",
			description: "Create a new note vault storage, switch between them using `)`",
			run:         vault.Cmd,
		},
		{
			name:        "version",
			description: "Show version information",
			run:         version.Cmd,
		},
		{
			name:        "export",
			description: "Export the SQLite database to a Obsidian vault.",
			run:         export.Cmd,
		},
		{
			name:        "logout",
			description: "Removed the cached credentials",
			run:         logout.Cmd,
		},
	}
}

func main() {
	// Setup Logging
	closer, err := utils.SetupLog()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer closer()

	flags, commands := parser.SplitCmdsAndFlags(os.Args[1:])

	// Handle --help flag
	for _, flag := range flags {
		if flag == "--help" {
			commands = append(commands, "help")
		}
	}

	// Handle commands
	if len(commands) > 0 {
		command := commands[0]
		args := commands[1:]

		for _, cmd := range COMMANDS {
			if strings.EqualFold(cmd.name, command) {
				os.Exit(cmd.run(args...))
			}
		}

		fmt.Printf("Unknown command: %s\n", command)
		helpCmd()
		os.Exit(1)
	}

	// Default: start TUI
	startTUI(flags...)
}
