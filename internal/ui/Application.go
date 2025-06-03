package ui

import (
	"merlion/internal/store"
	"merlion/internal/store/cloud"
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
	store *store.Manager // This manager is now directly passed in
	// credentialsManager *cloud.CredentialsManager // Keep if needed for login attempts, or pass to LoginUI directly
	// themeManager *styles.ThemeManager // Keep for theming
}

func NewModel(manager *store.Manager, themeManager *styles.ThemeManager, credentialsManager *cloud.CredentialsManager) (Model, error) {
	// Default to NoteUI since we expect a local store to be available.
	// LoginUI could be triggered by an explicit logout action or if manager is somehow nil (though it shouldn't be with local store).
	initialUI := navigation.NoteUI
	if manager == nil || manager.Name() == "unknown" { // "unknown" if activeStore was nil during manager creation
		// This case should ideally not happen if local store is always provided.
		// If it does, or if we want to explicitly offer login for cloud:
		initialUI = navigation.LoginUI
		// Consider logging a warning if manager is nil here, as it's unexpected with a default local store.
		// log.Warn("Store manager is nil or uninitialized at NewModel, defaulting to LoginUI")
	} else if manager.Name() == "local-sqlite" {
		// If it's the local sqlite store, we might still want to check if the user *wants* to log in
		// For now, we'll just go to NoteUI. A later feature could be to prompt for login.
		// Or, check if credentialsManager has creds, and if so, perhaps prompt for cloud sync or switch.
		// For this step, we simplify: if local store exists, use it.
		initialUI = navigation.NoteUI
	}
	// The original logic for checking credentialsManager.LoadCredentials()
	// and potentially switching to NoteUI if creds exist is now partly handled by
	// the fact that we *start* with a local store.
	// If we want to auto-login to cloud if creds are present, that logic would need to be re-introduced,
	// perhaps by trying to switch the manager after initialization.

	// Create views, passing the provided manager
	views := make(map[navigation.CurrentUI]navigation.View)
	views[navigation.CreateUI] = create.NewModel(manager, themeManager)
	// LoginUI might not need the global manager initially if it's for cloud login,
	// but it takes credentialsManager.
	views[navigation.LoginUI] = login.NewModel(credentialsManager, themeManager)
	views[navigation.NoteUI] = NotesUI.NewModel(manager, themeManager)
	views[navigation.DialogUI] = dialog.NewModel(manager, themeManager)
	views[navigation.ManageUI] = manage.NewModel(manager, themeManager)

	return Model{
		state: initialUI,
		views: views,
		store: manager, // Use the passed-in manager
	}, nil
}

func (m Model) Init() tea.Cmd {
	return m.views[m.state].Init()
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case navigation.SwitchUIMsg:
		m.state = msg.NewState
		return m, m.views[m.state].Init(msg.Args...)

	case navigation.LoginMsg:
		m.store = msg.Manager
		var cmds []tea.Cmd
		for i, view := range m.views {
			storeManager := msg.Manager
			updatedView := view.SetClient(storeManager)
			m.views[i] = updatedView
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
