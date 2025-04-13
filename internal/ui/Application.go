package ui

import (
	"merlion/internal/api"
	"merlion/internal/auth"
	"merlion/internal/store"
	"merlion/internal/styles"
	"merlion/internal/ui/create"
	"merlion/internal/ui/dialog"
	"merlion/internal/ui/login"
	"merlion/internal/ui/manage"
	"merlion/internal/ui/navigation"
	NotesUI "merlion/internal/ui/notes"

	tea "github.com/charmbracelet/bubbletea"
)

type Model struct {
	state navigation.CurrentUI
	views map[navigation.CurrentUI]navigation.View
	store *store.Manager
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

	manager := store.NewManager(client)

	// Create views
	views := make(map[navigation.CurrentUI]navigation.View)
	views[navigation.CreateUI] = create.NewModel(manager, themeManager)
	views[navigation.LoginUI] = login.NewModel(credentialsManager, themeManager)
	views[navigation.NoteUI] = NotesUI.NewModel(manager, themeManager)
	views[navigation.DialogUI] = dialog.NewModel(manager, themeManager)
	views[navigation.ManageUI] = manage.NewModel(manager, themeManager)

	return Model{
		state: initialUI,
		views: views,
		store: manager,
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
		m.store = msg.Manager
		var cmds []tea.Cmd
		for _, view := range m.views {
			storeManager := msg.Manager
			if cmd := view.SetClient(storeManager); cmd != nil {
				cmds = append(cmds, cmd)
			}
		}
		m.state = navigation.NoteUI
		cmds = append(cmds, m.views[m.state].Init())
		return m, tea.Batch(cmds...)

	case navigation.OpenDialogMsg:
		m.state = navigation.DialogUI
		view, cmd := m.views[m.state].Update(msg)
		m.views[m.state] = view
		return m, tea.Batch(cmd, tea.WindowSize())

	case navigation.OpenManageMsg:
		m.state = navigation.ManageUI
		view, cmd := m.views[m.state].Update(msg)
		m.views[m.state] = view
		return m, tea.Batch(cmd, tea.WindowSize())
	}

	view, cmd := m.views[m.state].Update(msg)
	m.views[m.state] = view
	return m, cmd
}

func (m Model) View() string {
	return m.views[m.state].View()
}
