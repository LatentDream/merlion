package ui

import (
	NotesUI "merlion/internal/ui/notes"
	"merlion/internal/ui/notes/create"

	tea "github.com/charmbracelet/bubbletea"
)

type state int

const (
	showNotes state = iota
	showCreate
)

type Model struct {
	state  state
	notes  NotesUI.Model
	create create.Model
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch m.state {

	case showNotes:
		var cmd tea.Cmd
		notesModel, cmd := m.notes.Update(msg)
		m.notes = notesModel.(NotesUI.Model)
		return m, cmd

	case showCreate:
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
	case showNotes:
		return m.notes.View()
	default:
		return m.create.View()
	}
}
