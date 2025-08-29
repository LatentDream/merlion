package main

import (
	"fmt"
	version "merlion/cmd/merlion/version"
	_ "net/http/pprof"
	"os"
	"strings"
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
	fmt.Println("    https://merlion.dev/login")
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

// parseArgs separates flags from commands and returns both
func parseArgs(args []string) (flags []string, commands []string) {
	for _, arg := range args {
		if strings.HasPrefix(arg, "--") {
			flags = append(flags, arg)
		} else {
			commands = append(commands, arg)
		}
	}
	return flags, commands
}

func main() {
	flags, commands := parseArgs(os.Args[1:])

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
			if strings.ToLower(cmd.name) == strings.ToLower(command) {
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
