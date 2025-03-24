package Notes

import (
	"fmt"
	"merlion/internal/api"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/log"
	"github.com/charmbracelet/x/editor"
)

// EDITOR ---

type editorFinishedMsg struct {
	err error
}

func (m *Model) openEditor() tea.Cmd {
	note := m.getCurrentNote(true)
	if note == nil {
		log.Fatalf("Trying to edit a nil note")
	}
	if note.Content == nil {
		log.Fatalf("Trying to open a note with content == nil")
	}
	log.Info("Editing:", note.NoteID)

	// Create a temporary file for editing
	tmpfile, err := os.CreateTemp("", "note-*.md")
	if err != nil {
		return func() tea.Msg {
			return editorFinishedMsg{fmt.Errorf("could not create temp file: %w", err)}
		}
	}

	if _, err := tmpfile.WriteString(*note.Content); err != nil {
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
		if err != nil {
			// Don't delete the file, instead show where it is
			return editorFinishedMsg{fmt.Errorf("editor failed: %w\nRecovery file location: %s", err, tmpfile.Name())}
		}

		// Read the edited content
		newContent, err := os.ReadFile(tmpfile.Name())
		if err != nil {
			return editorFinishedMsg{fmt.Errorf("failed to read edited content: %w\nRecovery file location: %s", err, tmpfile.Name())}
		}

		// Update note content through your API
		content := string(newContent)
		note.Content = &content

		// Update the main note list
		updated := false
		for _, groundTruthNote := range m.allNotes {
			if groundTruthNote.NoteID == note.NoteID {
				*groundTruthNote.Content = content
				updated = true
				break
			}
		}
		if !updated {
			log.Fatalf("Master list didn't get updated after Editor Finish")
		}

		// Update the note to backend
		req := note.ToCreateRequest()
		_, err = m.client.UpdateNote(note.NoteID, req)
		if err != nil {
			log.Errorf("Not able to save the note %s", note.NoteID)
			return editorFinishedMsg{fmt.Errorf("failed to save the edited content: %w\nRecovery file location: %s", err, tmpfile.Name())}
		}

		// Only remove the file on successful completion
		if err := os.Remove(tmpfile.Name()); err != nil {
			log.Errorf("Failed to remove temporary file %s: %v", tmpfile.Name(), err)
		}
		return editorFinishedMsg{nil}
	})
}

// Fetch Notes ---

type noteContentMsg struct {
	NoteId  string
	Content string
}
type errMsg struct{ err error }

func fetchNoteContent(client *api.Client, noteId string) tea.Cmd {
	return func() tea.Msg {
		res, err := client.GetNote(noteId)
		if err != nil {
			return errMsg{err}
		}
		return noteContentMsg{NoteId: res.NoteID, Content: *res.Content}
	}
}

type NotesLoadedMsg struct {
	Notes []api.Note
	Err   error
}

type notesLoadedMsg = NotesLoadedMsg

func (m Model) loadNotes() tea.Cmd {
	return func() tea.Msg {
		if m.client == nil {
			log.Fatalf("In CMD - Trying to load notes without any client")
		}
		notes, err := m.client.ListNotes()
		return notesLoadedMsg{Notes: notes, Err: err}
	}
}
