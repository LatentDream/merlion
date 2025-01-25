package navigation

import (
	"merlion/internal/api"

	tea "github.com/charmbracelet/bubbletea"
)

type CurrentUI int

const (
	NoteUI CurrentUI = iota
	CreateUI
	LoginUI
)

type SwitchUIMsg struct {
	NewState CurrentUI
}

type LoginMsg struct {
	Client *api.Client
}

type View interface {
	Init() tea.Cmd
	Update(tea.Msg) (View, tea.Cmd)
	View() string
	SetClient(*api.Client) tea.Cmd
}

// Global CMD to switch View

func SwitchUICmd(newState CurrentUI) tea.Cmd {
	return func() tea.Msg {
		return SwitchUIMsg{NewState: newState}
	}
}

func LoginCmd(client *api.Client) tea.Cmd {
	return func() tea.Msg {
		return LoginMsg{Client: client}
	}
}
