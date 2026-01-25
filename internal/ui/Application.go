package ui

import (
	"merlion/internal/config"
	"merlion/internal/context"
	"merlion/internal/vault"
	"merlion/internal/vault/cloud"
	"merlion/internal/ui/create"
	"merlion/internal/ui/dialog"
	"merlion/internal/ui/manage"
	"merlion/internal/ui/navigation"
	NotesUI "merlion/internal/ui/notes"

	tea "github.com/charmbracelet/bubbletea"
)

type Model struct {
	state navigation.CurrentUI
	views map[navigation.CurrentUI]navigation.View
	store *vault.Manager
}

func NewModel(cfg *config.UserConfig, credentialsManager *cloud.CredentialsManager, ctx *context.Context) (Model, error) {
	initialUI := navigation.NoteUI

	manager := vault.NewManager(cfg, credentialsManager, ctx.DefaultToCloud)

	// Create views
	views := make(map[navigation.CurrentUI]navigation.View)
	views[navigation.CreateUI] = create.NewModel(manager, ctx.ThemeManager)
	views[navigation.NoteUI] = NotesUI.NewModel(manager, ctx.ThemeManager, ctx.FirstTab)
	views[navigation.DialogUI] = dialog.NewModel(manager, ctx.ThemeManager)
	views[navigation.ManageUI] = manage.NewModel(manager, ctx.ThemeManager)

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
		return m, m.views[m.state].Init(msg.Args...)

	case navigation.LoginMsg:
		m.store.UpdateCloudClient(msg.Client)
		var cmds []tea.Cmd
		for i, view := range m.views {
			cloudClient := msg.Client
			updatedView := view.SetCloudClient(cloudClient)
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
