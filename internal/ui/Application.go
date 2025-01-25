package ui

import (
	"fmt"
	"merlion/internal/api"
	"merlion/internal/auth"
	"merlion/internal/styles"
	"merlion/internal/ui/login"
	"merlion/internal/ui/navigation"
	NotesUI "merlion/internal/ui/notes"
	"merlion/internal/ui/create"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/log"
)

type Model struct {
	state  navigation.CurrentUI
	notes  NotesUI.Model
	create create.Model
	login  login.Model
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {

	// Global Signal to swtich UI
	switch msg := msg.(type) {
	case navigation.SwitchUIMsg:
		m.state = msg.NewState
		return m, nil
	case navigation.LoginMsg:
		m.notes.SetClient(msg.Client)
		m.create.SetClient(msg.Client)
		m.state = navigation.NoteUI
		return m, nil
	}

	// Display Update
	switch m.state {
	case navigation.LoginUI:
		var cmd tea.Cmd
		loginModel, cmd := m.login.Update(msg)
		m.login = loginModel.(login.Model)
		return m, cmd
	case navigation.NoteUI:
		var cmd tea.Cmd
		notesModel, cmd := m.notes.Update(msg)
		m.notes = notesModel.(NotesUI.Model)
		return m, cmd
	case navigation.CreateUI:
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
	case navigation.LoginUI:
		return m.login.View()
	case navigation.NoteUI:
		return m.notes.View()
	default:
		return m.create.View()
	}
}

func NewModel(credentialsManager *auth.CredentialsManager, themeManager *styles.ThemeManager) (Model, error) {
	// Verify user credentials
	initialUI := navigation.LoginUI
	var client *api.Client = nil
	var err error = nil

	creds, _ := credentialsManager.LoadCredentials()
	if creds != nil {
		initialUI = navigation.NoteUI
		client, err = api.NewClient(creds)
		if err != nil {
			log.Error("Failed to create API client: %v", err)
			return Model{}, err
		}
	}

	// Init all model
	loginModel, err := login.NewModel(credentialsManager, themeManager)
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

	return Model{
		state:  initialUI,
		login:  loginModel,
		notes:  notesModel,
		create: createModel,
	}, nil
}
