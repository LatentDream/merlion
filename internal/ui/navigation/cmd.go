package navigation

import (
	"merlion/internal/store"

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
	Args     []any
}

type LoginMsg struct {
	Manager *store.Manager
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
	SetClient(*store.Manager) tea.Cmd
}

// Global CMD to switch View

func SwitchUICmd(newState CurrentUI) tea.Cmd {
	return func() tea.Msg {
		return SwitchUIMsg{NewState: newState}
	}
}

func LoginCmd(storeManager *store.Manager) tea.Cmd {
	return func() tea.Msg {
		return LoginMsg{Manager: storeManager}
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
