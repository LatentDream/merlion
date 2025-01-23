package ui

import (
	"fmt"
	"merlion/internal/api"
	"merlion/internal/auth"
	"merlion/internal/styles"
	"merlion/internal/ui/login"
	NotesUI "merlion/internal/ui/notes"
	"merlion/internal/ui/notes/create"

	tea "github.com/charmbracelet/bubbletea"
)

type state int

const (
	showNotes state = iota
	showCreate
	showLogin
)

type Model struct {
	state  state
	notes  NotesUI.Model
	create create.Model
	login  login.Model
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch m.state {

	case showLogin:
		var cmd tea.Cmd
		loginModel, cmd := m.login.Update(msg)
		m.login = loginModel.(login.Model)
		return m, cmd

	case showNotes:
		var cmd tea.Cmd
		notesModel, cmd := m.notes.Update(msg)
		m.notes = notesModel.(NotesUI.Model)
		return m, cmd

	case showCreate:
		var cmd tea.Cmd
		createModel, cmd := m.create.Update(msg)
		m.create = createModel.(create.Model)
		return m, cmd

	default:
		return m, nil
	}
}

func (m Model) View() string {
	switch m.state {
	case showLogin:
		return m.login.View()
	case showNotes:
		return m.notes.View()
	default:
		return m.create.View()
	}
}

func NewModel(credentialsManager *auth.CredentialsManager, themeManager *styles.ThemeManager) (Model, error) {
	// Verify user credentials

	// TODO
	
	var client *api.Client = nil

	// Init all model
	loginModel, err := login.NewModel(themeManager)
	if err != nil {
		return Model{}, fmt.Errorf("failed to create login model: %w", err)
	}

	notesModel, err := NotesUI.NewModel(client, themeManager)
	if err != nil {
		return Model{}, fmt.Errorf("failed to create notes model: %w", err)
	}

	createModel, err := create.NewModel(client, themeManager)
	if err != nil {
		return Model{}, fmt.Errorf("failed to create create model: %w", err)
	}

	// Determine which UI to start

	return Model{
		state:  showLogin,
		login:  loginModel,
		notes:  notesModel,
		create: createModel,
	}, nil
}
