package application

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
