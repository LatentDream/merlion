/* TODO: Refactoring in progress
 * - [x] New store package to encapsulate the Cloud behind an interface (will allow local storage)
 * - [ ] The backend should be able to receive nil ptr to the content & not delete the related content
 *       > remove the needs for internal/ui/manage/model.go:FetchNote
 * - [ ] Models shouldn't keep their own version of []model.Note, but should refer to the store
 */
package store

import (
	"merlion/internal/model"
	"strings"
)

// Manager handles note operations through an underlying store implementation.
// Note: When using ListNoteMetadata() or accessing Manager.notes directly, note content
// may not be populated (content field may be nil). For guaranteed access to note
// content, use GetFullNote() which will always return the complete note with content.
type Manager struct {
	activeStore store
	name        string
	inited      bool
	// notes contains a cached list of notes, but their content field may be nil
	// Use GetNote() to retrieve the complete note with content
	notes []model.Note
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
func (m *Manager) GetFullNote(noteID string) (*model.Note, error) {
	return m.activeStore.GetNote(noteID)
}

// ListNoteMetadata returns all notes from the store, but with potentially empty content fields.
// Note: Returned notes may not have their content field populated.
// To access a note's full content, use GetFullNote() with the note's ID.
func (m *Manager) ListNoteMetadata() ([]model.Note, error) {
	notes, err := m.activeStore.ListNotes()
	if err != nil {
		return notes, err
	}
	m.notes = notes
	return notes, nil
}

// UpdateNote modifies an existing note with the provided changes.
func (m *Manager) UpdateNote(noteId string, changes model.CreateNoteRequest) (*model.Note, error) {
	return m.activeStore.UpdateNote(noteId, changes)
}

// GetTags returns all available tags from the cached notes.
func (m *Manager) GetTags() []string {
	tagMap := make(map[string]bool)
	for _, note := range m.notes {
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
	return m.activeStore.CreateNote(req)
}

// DeleteNote removes a note by its ID.
func (m *Manager) DeleteNote(noteID string) error {
	return m.activeStore.DeleteNote(noteID)
}
