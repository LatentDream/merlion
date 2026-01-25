package navigation

import (
	"merlion/internal/vault/cloud"

	tea "github.com/charmbracelet/bubbletea"
)

type CurrentUI int

const (
	NoteUI CurrentUI = iota
	CreateUI
	ManageUI
	DialogUI
)

type Level int

const (
	InfoLvl Level = iota
	DangerLvl
)

type SwitchUIMsg struct {
	NewState CurrentUI
	Args     []any
}

type LoginMsg struct {
	Client *cloud.Client
}

type OpenDialogMsg struct {
	Title     string
	Subtitle  string
	Level     Level
	OnConfirm func()
	ReturnUI  CurrentUI
}

type OpenManageMsg struct {
	NoteId string
}

type View interface {
	Init(...any) tea.Cmd
	Update(tea.Msg) (View, tea.Cmd)
	View() string
	SetCloudClient(*cloud.Client) View
}

// Global CMD to switch View

func SwitchUICmd(newState CurrentUI, args []any) tea.Cmd {
	return func() tea.Msg {
		return SwitchUIMsg{NewState: newState, Args: args}
	}
}

func AskConfirmationCmd(title string, subtitle string, level Level, onConfirm func(), returnUI CurrentUI) tea.Cmd {
	return func() tea.Msg {
		return OpenDialogMsg{Title: title, Subtitle: subtitle, Level: level, OnConfirm: onConfirm, ReturnUI: returnUI}
	}
}

func OpenManageViewCmd(noteId string) tea.Cmd {
	return func() tea.Msg {
		return OpenManageMsg{NoteId: noteId}
	}
}
