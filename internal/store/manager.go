/* TODO: Refactoring in progress
 * - [x] New store package to encapsulate the Cloud behind an interface (will allow local storage)
 * - [ ] The backend should be able to receive nil ptr to the content & not delete the related content
 *       > remove the needs for internal/ui/manage/model.go:FetchNote
 * - [ ] Models shouldn't keep their own version of []model.Note, but should refer to the store
 */
package store

import (
	"errors"
	"log"
	"merlion/internal/model"
	"merlion/internal/utils"
	"strings"
)

// ErrNoteNotFound is returned when a note with the given ID is not found
var ErrNoteNotFound = errors.New("note not found")

// Manager handles note operations through an underlying store implementation.
// Note: When using ListNoteMetadata() or accessing Manager.notes directly, note content
// may not be populated (content field may be nil). For guaranteed access to note
// content, use GetFullNote() which will always return the complete note with content.
type Manager struct {
	activeStore store
	name        string
	inited      bool
	// Notes contains a cached list of Notes, but their content field may be nil
	// Use GetNote() to retrieve the complete note with content
	Notes []model.Note
}

// NewManager creates a new manager with the given store implementation
// and initializes the notes metadata list (without full content).
func NewManager(store store) *Manager {
	return &Manager{
		activeStore: store,
		name:        store.Name(),
		inited:      false,
	}
}

// GetFullNote retrieves a specific note by ID with its complete content.
// This method guarantees that the returned note will have its content field populated.
func (m *Manager) GetFullNote(noteId string) (*model.Note, error) {
	note, err := m.activeStore.GetNote(noteId)
	if err != nil {
		return nil, err
	}

	found := false
	for i, cachedNote := range m.Notes {
		if cachedNote.NoteID == noteId {
			m.Notes[i] = *note
			found = true
			break
		}
	}
	if !found {
		log.Fatalf("User was able to get an undefined note - Should not happen")
	}
	return note, nil
}

// ListNoteMetadata returns all notes from the store, but with potentially empty content fields.
// Note: Returned notes may not have their content field populated.
// To access a note's full content, use GetFullNote() with the note's ID.
func (m *Manager) ListNoteMetadata() ([]model.Note, error) {
	notes, err := m.activeStore.ListNotes()
	if err != nil {
		return notes, err
	}
	m.Notes = notes
	return notes, nil
}

func (m *Manager) SearchById(noteId string) (*model.Note, error) {
	for _, note := range m.Notes {
		if note.NoteID == noteId {
			return &note, nil
		}
	}
	return nil, ErrNoteNotFound
}

func (m *Manager) SearchByTitle(title string) (*model.Note, error) {
	standardize := func(s string) string {
		return strings.TrimSpace(strings.ToLower(s))
	}
	searchTitle := standardize(title)
	for _, note := range m.Notes {
		currTitle := standardize(note.Title)
		if currTitle == searchTitle {
			return &note, nil
		}
	}
	return nil, ErrNoteNotFound

}

// GetTags returns all available tags from the cached notes.
func (m *Manager) GetTags() []string {
	tagMap := make(map[string]bool)
	for _, note := range m.Notes {
		for _, tag := range note.Tags {
			// If not in map, add it
			tagMap[strings.ToLower(tag)] = true
		}
	}
	tags := make([]string, 0, len(tagMap))
	for tag := range tagMap {
		tags = append(tags, tag)
	}
	return tags
}

// CreateNote creates a new note with the provided request data.
func (m *Manager) CreateNote(req model.CreateNoteRequest) (*model.Note, error) {
	note, err := m.activeStore.CreateNote(req)
	if err != nil {
		return nil, err
	}
	m.Notes = append(m.Notes, *note)
	return note, nil
}

// UpdateNote modifies an existing note with the provided changes and updates the metadata cache.
func (m *Manager) UpdateNote(noteId string, changes model.CreateNoteRequest) (*model.Note, error) {
	note, err := m.activeStore.UpdateNote(noteId, changes)
	if err != nil {
		return nil, err
	}

	found := false
	for i, cachedNote := range m.Notes {
		if cachedNote.NoteID == noteId {
			m.Notes[i] = *note
			found = true
			break
		}
	}
	if !found {
		log.Fatalf("User was able to update an undefined note - Should not happen")
	}
	return note, nil
}

// DeleteNote removes a note by its ID.
func (m *Manager) DeleteNote(noteId string) error {
	err := m.activeStore.DeleteNote(noteId)
	if err != nil {
		return err
	}
	idx := -1
	for i, cachedNote := range m.Notes {
		if cachedNote.NoteID == noteId {
			idx = i
			break
		}
	}
	if idx == -1 {
		log.Fatalf("Deleted a note which wasn't cached locally - Should not happen")
	}
	utils.Remove(m.Notes, idx)
	return nil
}
