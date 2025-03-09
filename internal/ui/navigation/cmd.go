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
}

type LoginMsg struct {
	Client *api.Client
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
