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
	"github.com/charmbracelet/log"
)

type state int

const (
	noteUI state = iota
	createUI
	loginUI
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

	case loginUI:
		var cmd tea.Cmd
		loginModel, cmd := m.login.Update(msg)
		m.login = loginModel.(login.Model)
		return m, cmd

	case noteUI:
		var cmd tea.Cmd
		notesModel, cmd := m.notes.Update(msg)
		m.notes = notesModel.(NotesUI.Model)
		return m, cmd

	case createUI:
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
	case loginUI:
		return m.login.View()
	case noteUI:
		return m.notes.View()
	default:
		return m.create.View()
	}
}

func NewModel(credentialsManager *auth.CredentialsManager, themeManager *styles.ThemeManager) (Model, error) {
	// Verify user credentials
	initialUI := loginUI
	var client *api.Client = nil
	var err error = nil

	creds, _ := credentialsManager.LoadCredentials()
	if creds != nil {
		initialUI = noteUI
		client, err = api.NewClient(creds)
		if err != nil {
			log.Error("Failed to create API client: %v", err)
			return Model{}, err
		}
	}

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
		state:  initialUI,
		login:  loginModel,
		notes:  notesModel,
		create: createModel,
	}, nil
}
