package Notes

import (
	"fmt"
	"merlion/internal/api"
	"os"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/log"
	"github.com/charmbracelet/x/editor"
)

// EDITOR ---

type editorFinishedMsg struct {
	err error
}

func (m *Model) openEditor(content string) tea.Cmd {
	// Create a temporary file for editing
	tmpfile, err := os.CreateTemp("", "note-*.md")
	if err != nil {
		return func() tea.Msg {
			return editorFinishedMsg{fmt.Errorf("could not create temp file: %w", err)}
		}
	}

	// Write content to temp file
	if _, err := tmpfile.WriteString(content); err != nil {
		os.Remove(tmpfile.Name())
		return func() tea.Msg {
			return editorFinishedMsg{fmt.Errorf("could not write to temp file: %w", err)}
		}
	}
	tmpfile.Close()

	// Create editor command
	cmd, err := editor.Cmd("Note", tmpfile.Name())
	if err != nil {
		os.Remove(tmpfile.Name())
		return func() tea.Msg {
			return editorFinishedMsg{fmt.Errorf("failed to create editor command: %w", err)}
		}
	}

	// Return command that will execute editor
	return tea.ExecProcess(cmd, func(err error) tea.Msg {
		defer os.Remove(tmpfile.Name())

		if err != nil {
			return editorFinishedMsg{fmt.Errorf("editor failed: %w", err)}
		}

		// Read the edited content
		newContent, err := os.ReadFile(tmpfile.Name())
		if err != nil {
			return editorFinishedMsg{fmt.Errorf("failed to read edited content: %w", err)}
		}

		// Update note content through your API
		if i := m.list.SelectedItem(); i != nil {
			note := i.(item).note
			content := string(newContent)
			note.Content = &content

			// Update the item in the model's list
			currentIndex := m.list.Index()
			items := m.list.Items()
			items[currentIndex] = item{note: note}
			m.list.SetItems(items)

			// Update the note to
			req := note.ToCreateRequest()
			_, err := m.client.UpdateNote(note.NoteID, req)
			if err != nil {
				log.Error("Not able to save the note %s", note.NoteID)
				return editorFinishedMsg{fmt.Errorf("failed to save the edited content: %w", err)}
			}
		}

		return editorFinishedMsg{nil}
	})
}

// Fetch Notes ---

type noteContentMsg string
type errMsg struct{ err error }

func fetchNoteContent(client *api.Client, noteId string) tea.Cmd {
	return func() tea.Msg {
		res, err := client.GetNote(noteId)
		if err != nil {
			return errMsg{err}
		}
		return noteContentMsg(*res.Content)
	}
}

type NotesLoadedMsg struct {
	Notes []api.Note
}

type notesLoadedMsg = NotesLoadedMsg

func (m Model) loadNotes() tea.Msg {
	items := make([]list.Item, len(m.list.Items()))
	for i, item := range m.list.Items() {
		items[i] = item
	}
	return notesLoadedMsg{Notes: nil}
}
