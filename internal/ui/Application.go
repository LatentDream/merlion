package ui

import (
	"merlion/internal/api"
	"merlion/internal/auth"
	"merlion/internal/styles"
	"merlion/internal/ui/create"
	"merlion/internal/ui/login"
	"merlion/internal/ui/navigation"
	NotesUI "merlion/internal/ui/notes"

	tea "github.com/charmbracelet/bubbletea"
)

type Model struct {
	state  navigation.CurrentUI
	views  map[navigation.CurrentUI]navigation.View
	client *api.Client
}

func NewModel(credentialsManager *auth.CredentialsManager, themeManager *styles.ThemeManager) (Model, error) {
	initialUI := navigation.LoginUI
	var client *api.Client = nil

	// Check credentials
	creds, _ := credentialsManager.LoadCredentials()
	if creds != nil {
		initialUI = navigation.NoteUI
		client, _ = api.NewClient(creds)
	}

	// Create views
	views := make(map[navigation.CurrentUI]navigation.View)
	views[navigation.LoginUI] = login.NewModel(credentialsManager, themeManager)
	views[navigation.NoteUI] = NotesUI.NewModel(client, themeManager)
	views[navigation.CreateUI] = create.NewModel(client, themeManager)

	return Model{
		state:  initialUI,
		views:  views,
		client: client,
	}, nil
}

func (m Model) Init() tea.Cmd {
	return m.views[m.state].Init()
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case navigation.SwitchUIMsg:
		m.state = msg.NewState
		return m, m.views[m.state].Init()
	case navigation.LoginMsg:
		m.client = msg.Client
		for _, view := range m.views {
			view.SetClient(msg.Client)
		}
		m.state = navigation.NoteUI
		return m, m.views[m.state].Init()
	}

	view, cmd := m.views[m.state].Update(msg)
	m.views[m.state] = view
	return m, cmd
}

func (m Model) View() string {
	return m.views[m.state].View()
}
