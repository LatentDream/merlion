package store

import (
	"fmt"
	"log"
	"merlion/internal/model"
	"merlion/internal/store/cloud"
	"merlion/internal/store/local"
	"merlion/internal/utils"
	"merlion/internal/utils/assert"
	"strings"
)

const panic__consistency_msg = "Storage was changed, but note wasn't refresh"

// Manager handles note operations through an underlying store implementation.
// Note: When using ListNoteMetadata() or accessing Manager.notes directly, note content
// may not be populated (content field may be nil). For guaranteed access to note
// content, use GetFullNote() which will always return the complete note with content.
type Manager struct {
	activeStore          Store
	Name                 string
	inited               bool
	internal__cloudStore *cloud.Client
	internal__localStore *local.Client
	internal__notesStore string
	// Notes contains a cached list of Notes, but their content field may be nil
	// Use GetNote() to retrieve the complete note with content
	Notes []model.Note
}

// NewManager creates a new manager with the given store implementation
// and initializes the notes metadata list (without full content).
func NewManager(cloudStore *cloud.Client, localStore *local.Client) *Manager {
	var defaultStore Store
	defaultStore = localStore
	if cloudStore != nil {
		defaultStore = cloudStore
	}

	return &Manager{
		activeStore:          defaultStore,
		Name:                 defaultStore.Name(),
		inited:               false,
		internal__cloudStore: cloudStore,
		internal__localStore: localStore,
		internal__notesStore: defaultStore.Name(),
	}
}

// UpdateCloudStore will swap the cloud storage for a new store
// Expect the store to be valid and functionnal, panic otherwise
func (m *Manager) UpdateCloudClient(client *cloud.Client) {
	if m.activeStore.Name() == cloud.Name {
		m.activeStore = client
		_, err := m.ListNoteMetadata()
		if err != nil {
			log.Fatalf("Failed to fetch from new cloud store: %v", err)
		}
	}
	m.internal__cloudStore = client
}

// NextStore swap the current underlying storage with the next registered one
// Dev needs to call ListNoteMetadata after calling this, otherwise a panic occur
func (m *Manager) NextStore() error {
	if m.Name == cloud.Name {
		m.activeStore = m.internal__localStore
	} else {
		if m.internal__cloudStore == nil {
			return fmt.Errorf("Cloud store is nil")
		}
		m.activeStore = m.internal__cloudStore
	}
	m.Name = m.activeStore.Name()

	return nil
}

// GetFullNote retrieves a specific note by ID with its complete content.
// This method guarantees that the returned note will have its content field populated.
func (m *Manager) GetFullNote(noteId string) (*model.Note, error) {
	assert.Eq(m.internal__notesStore, m.activeStore.Name(), panic__consistency_msg)

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
	m.internal__notesStore = m.activeStore.Name()
	return notes, nil
}

func (m *Manager) SearchById(noteId string) *model.Note {
	assert.Eq(m.internal__notesStore, m.activeStore.Name(), panic__consistency_msg)

	for _, note := range m.Notes {
		if note.NoteID == noteId {
			return &note
		}
	}
	return nil
}

func (m *Manager) SearchByTitle(title string) *model.Note {
	assert.Eq(m.internal__notesStore, m.activeStore.Name(), panic__consistency_msg)

	standardize := func(s string) string {
		return strings.TrimSpace(strings.ToLower(s))
	}
	searchTitle := standardize(title)
	for _, note := range m.Notes {
		currTitle := standardize(note.Title)
		if currTitle == searchTitle {
			return &note
		}
	}
	return nil

}

// GetTags returns all available tags from the cached notes.
func (m *Manager) GetTags() []string {
	assert.Eq(m.internal__notesStore, m.activeStore.Name(), panic__consistency_msg)

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
	assert.Eq(m.internal__notesStore, m.activeStore.Name(), panic__consistency_msg)

	note, err := m.activeStore.CreateNote(req)
	if err != nil {
		return nil, err
	}
	m.Notes = append(m.Notes, *note)
	return note, nil
}

// UpdateNote modifies an existing note with the provided changes and updates the metadata cache.
func (m *Manager) UpdateNote(noteId string, changes model.CreateNoteRequest) (*model.Note, error) {
	assert.Eq(m.internal__notesStore, m.activeStore.Name(), panic__consistency_msg)

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
	assert.Eq(m.internal__notesStore, m.activeStore.Name(), panic__consistency_msg)

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
