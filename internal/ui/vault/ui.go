package vault

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

type vaultType string

const (
	Cloud    vaultType = "Cloud"
	SQLite   vaultType = "SQLite"
	Obsidian vaultType = "Obsidian Vault"
)

// Model states
type state int

const (
	stateChooseVault state = iota
	stateObsidianVault
	stateDone
)

type Model struct {
	state       state
	cursor      int
	choices     []vaultType
	textInput   textinput.Model
	ChosenVault vaultType
	ChosenPath  string
}

func NewModel() Model {
	ti := textinput.New()
	ti.Focus()

	return Model{
		state:     stateChooseVault,
		cursor:    0,
		choices:   []vaultType{Cloud, SQLite, Obsidian},
		textInput: ti,
	}
}

func (m Model) Init() tea.Cmd {
	return textinput.Blink
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
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

func (m Model) handleEnter() (tea.Model, tea.Cmd) {
	switch m.state {
	case stateChooseVault:
		m.ChosenVault = m.choices[m.cursor]

		switch m.ChosenVault {
		case SQLite, Cloud:
			// SQLite and Cloud don't need input, go straight to done
			m.state = stateDone
			return m, tea.Quit
		case Obsidian:
			m.state = stateObsidianVault
			homeDir, _ := os.UserHomeDir()
			defaultPath := filepath.Join(homeDir, "notes")
			m.textInput.SetValue(defaultPath)
		}

	case stateObsidianVault:
		m.ChosenPath = m.textInput.Value()
		m.state = stateDone
		return m, tea.Quit
	}

	return m, nil
}

func (m Model) View() string {
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
		switch m.ChosenVault {
		case SQLite:
			s.WriteString("Using default SQLite database\n")
		case Obsidian:
			s.WriteString(fmt.Sprintf("Obsidian vault: %s\n", m.ChosenPath))
		case Cloud:
			s.WriteString("Cloud vault will be configured\n")
		}
	}

	return s.String()
}
