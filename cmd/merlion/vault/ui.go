package vault

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"merlion/cmd/merlion/parser"
	"merlion/internal/config"
	"merlion/internal/store/cloud"
	"merlion/internal/store/files"
	"merlion/internal/store/sqlite"
	"merlion/internal/styles"
	"merlion/internal/ui/login"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/log"
)

type vaultType string

const (
	vaultCloud    vaultType = "Cloud"
	vaultSQLite   vaultType = "SQLite"
	vaultObsidian vaultType = "Obsidian Vault"
)

// Model states
type state int

const (
	stateChooseVault state = iota
	stateObsidianVault
	stateDone
)

type model struct {
	state       state
	cursor      int
	choices     []vaultType
	textInput   textinput.Model
	chosen      vaultType
	obsidianDir string
}

func initialModel() model {
	ti := textinput.New()
	ti.Focus()

	return model{
		state:     stateChooseVault,
		cursor:    0,
		choices:   []vaultType{vaultCloud, vaultSQLite, vaultObsidian},
		textInput: ti,
	}
}

func (m model) Init() tea.Cmd {
	return textinput.Blink
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "up", "k":
			if m.state == stateChooseVault && m.cursor > 0 {
				m.cursor--
			}
		case "down", "j":
			if m.state == stateChooseVault && m.cursor < len(m.choices)-1 {
				m.cursor++
			}
		case "enter":
			return m.handleEnter()
		}
	}

	// Update text input when in input states
	if m.state != stateChooseVault && m.state != stateDone {
		var cmd tea.Cmd
		m.textInput, cmd = m.textInput.Update(msg)
		return m, cmd
	}

	return m, nil
}

func (m model) handleEnter() (tea.Model, tea.Cmd) {
	switch m.state {
	case stateChooseVault:
		m.chosen = m.choices[m.cursor]

		switch m.chosen {
		case vaultSQLite, vaultCloud:
			// SQLite and Cloud don't need input, go straight to done
			m.state = stateDone
			return m, tea.Quit
		case vaultObsidian:
			m.state = stateObsidianVault
			homeDir, _ := os.UserHomeDir()
			defaultPath := filepath.Join(homeDir, "notes")
			m.textInput.SetValue(defaultPath)
		}

	case stateObsidianVault:
		m.obsidianDir = m.textInput.Value()
		m.state = stateDone
		return m, tea.Quit
	}

	return m, nil
}

func (m model) View() string {
	var s strings.Builder

	switch m.state {
	case stateChooseVault:
		s.WriteString("You don't have any place to store your notes, please choose one:\n")
		for i, choice := range m.choices {
			cursor := " "
			if m.cursor == i {
				cursor = ">"
			}
			s.WriteString(fmt.Sprintf("%s %s\n", cursor, choice))
		}
		s.WriteString("\n\033[37mTo add new provider in the future, you can use \033[36;4mmerlion vault\033[0m\n")

	case stateObsidianVault:
		s.WriteString("Select a vault, or create a new one:\n")
		s.WriteString("Root Folder: ")
		s.WriteString(m.textInput.View())
		s.WriteString("\n")

	case stateDone:
		s.WriteString("✓ Vault configured!\n\n")
		switch m.chosen {
		case vaultSQLite:
			s.WriteString("Using default SQLite database\n")
		case vaultObsidian:
			s.WriteString(fmt.Sprintf("Obsidian vault: %s\n", m.obsidianDir))
		case vaultCloud:
			s.WriteString("Cloud vault will be configured\n")
		}
	}

	return s.String()
}

// Integration functions
func ChooseVault() int {
	p := tea.NewProgram(initialModel())
	finalModel, err := p.Run()
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return 1
	}

	m := finalModel.(model)

	// Now call the appropriate function based on choice
	switch m.chosen {
	case vaultSQLite:
		return newSQLiteVault()
	case vaultObsidian:
		return newFilesVault(m.obsidianDir)
	case vaultCloud:
		return newCloudVault()
	}

	return 0
}

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
		if vault.Provider == cloud.Name {
			return 0
		}
	}
	cfg.Vaults = append(cfg.Vaults, config.Vault{
		Provider: cloud.Name,
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
		Provider: sqlite.Name,
		Name:     sqlite.Name,
		Path:     "Not Used - Use MERLION_DB_PATH env var to override",
	})
	cfg.Save()
	return 0
}
