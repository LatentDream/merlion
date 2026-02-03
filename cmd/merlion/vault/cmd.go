// Package vault contains the logic on how to add new note storage clients
// Used with the `merlion new` command
package vault

import (
	"fmt"
	"os"
	"path/filepath"

	"merlion/cmd/merlion/parser"
	"merlion/internal/config"
	"merlion/internal/styles"
	"merlion/internal/ui/login"
	"merlion/internal/ui/vault"
	"merlion/internal/utils"
	"merlion/internal/vault/cloud"
	"merlion/internal/vault/files"
	"merlion/internal/vault/sqlite"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/log"
)

func Cmd(args ...string) int {
	if utils.Contains(args, "--help") || utils.Contains(args, "-h") {
		printVaultHelp(true)
	}

	if len(args) == 0 {
		return ChooseVault()
	}

	arg, args := parser.GetArg(args, printVaultHelp)
	switch arg {
	case "sqlite":
		return newSQLiteVault(args...)
	case "files":
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

	fmt.Println("Usage: merlion vault [<provider>]")
	fmt.Println("Provider: sqlite, files [<obsidian-vault-path>], cloud")
	fmt.Println("  - sqlite: create a new SQLite database")
	fmt.Println("  - files <obsidian-vault-path>: create a new local Obsidian vault")
	fmt.Println("  - cloud: create a new cloud storage provider")
	fmt.Println("")
	fmt.Println("To remove a vault")
	fmt.Println("   - sqlite & files -> edit ~/.merlion/config.json")
	fmt.Println("   - cloud -> merlion logout")

	if invalidArgs {
		os.Exit(1)
	}
	os.Exit(0)
}

func ChooseVault() int {
	p := tea.NewProgram(vault.NewModel())
	finalModel, err := p.Run()
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return 1
	}

	m := finalModel.(vault.Model)

	// Now call the appropriate function based on choice
	switch m.ChosenVault {
	case vault.SQLite:
		return newSQLiteVault()
	case vault.Obsidian:
		return newFilesVault(m.ChosenPath)
	case vault.Cloud:
		return newCloudVault()
	}

	return 0
}

// TODO: move the vault creation logic in the related package
// Will cause circular dependency => Let's deal with it later

func newCloudVault(args ...string) int {
	tm, err := styles.NewThemeManager()
	if err != nil {
		log.Fatalf("Failed to initialize theme manager: %v", err)
	}
	credMgr, err := cloud.NewCredentialsManager()
	if err != nil {
		log.Fatalf("Failed to initialize credentials manager: %v", err)
	}

	p := tea.NewProgram(login.NewModel(credMgr, tm))
	_, err = p.Run()
	if err != nil {
		log.Fatalf("Error running program: %v", err)
	}

	// TODO: could be more generic to have multiple accounts
	cfg := config.Load()
	for _, vault := range cfg.Vaults {
		if vault.Provider == cloud.Type {
			return 0
		}
	}
	cfg.Vaults = append(cfg.Vaults, config.Vault{
		Provider: cloud.Type,
		Name:     cloud.Name,
		Path:     "Not Used - credentials are stored in ~/.merlion/credentials.json",
	})
	cfg.Save()
	return 0
}

func newFilesVault(args ...string) int {
	root, _ := parser.GetArg(args, printVaultHelp)

	absPath, err := filepath.Abs(root)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: invalid path: %v\n", err)
		return 1
	}
	cleanPath := filepath.Clean(absPath)
	if cleanPath == "." || cleanPath == "" {
		fmt.Fprintf(os.Stderr, "Error: invalid path provided\n")
		return 1
	}

	cfg := config.Load()
	cfg.Vaults = append(cfg.Vaults, config.Vault{
		Provider: files.Type,
		Name:     cleanPath,
		Path:     cleanPath,
	})
	cfg.Save()
	return 0
}

func newSQLiteVault(args ...string) int {
	// TODO: could be more generic accepting a path for the db
	cfg := config.Load()
	cfg.Vaults = append(cfg.Vaults, config.Vault{
		Provider: sqlite.Type,
		Name:     sqlite.Name,
		Path:     "Not Used - Use MERLION_DB_PATH env var to override",
	})
	cfg.Save()
	return 0
}
