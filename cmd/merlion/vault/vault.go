// Package vault contains the logic on how to add new note storage clients
// Used with the `merlion new` command
package vault

import (
	"fmt"
	"merlion/cmd/merlion/parser"
	"merlion/internal/utils"
	"os"
)


func VaultCmd(args ...string) int {
	if len(args) == 0 || len(args) > 2 || utils.Contains(args, "--help") || utils.Contains(args, "-h") {
		printVaultHelp(true)
	}

	arg, args := parser.GetArg(args, printVaultHelp)
	switch arg {
	case "sqlite":
		return newSQLiteVault(args...)
	case "file":
		return newFilesVault(args...)
	case "cloud":
		return newCloudVault(args...)
	default:
		printVaultHelp(true)
	}
	return 0
}

func printVaultHelp(invalidArgs bool) {
	if invalidArgs {
		fmt.Println("Invalid arguments")
	}

	fmt.Println("Usage: merlion vault <provider>")
	fmt.Println("Provider: sqlite, file [<obsidian-vault-path>], cloud")
	fmt.Println("  - sqlite: create a new SQLite database")
	fmt.Println("  - file <obsidian-vault-path>: create a new local Obsidian vault")
	fmt.Println("  - cloud: create a new cloud storage provider")
	fmt.Println("Examples:")
	fmt.Println("  merlion vault sqlite")
	fmt.Println("  merlion vault file ~/notes")
	fmt.Println("  merlion vault cloud")

	if invalidArgs {
		os.Exit(1)
	}
	os.Exit(0)
}

func newCloudVault(args ...string) int {
	return 0
}

func newFilesVault(args ...string) int {
	return 0
}

func newSQLiteVault(args ...string) int {
	return 0
}
