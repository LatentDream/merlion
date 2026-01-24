package store

import (
	"strings"

	"merlion/internal/model"
	"merlion/internal/store/cloud"
	"merlion/internal/store/files"
	sqlite "merlion/internal/store/sqlite"
	"merlion/internal/utils"
	"merlion/internal/utils/assert"
	"github.com/charmbracelet/log"
)

const panic__consistency_msg = "Storage was changed, but note wasn't refresh"

// Manager handles note operations through an underlying store implementation.
// Note: When using ListNoteMetadata() or accessing Manager.notes directly, note content
// may not be populated (content field may be nil). For guaranteed access to note
// content, use GetFullNote() which will always return the complete note with content.
type Manager struct {
	activeStore          Store
	Name                 string
	stores               []Store
	internal__notesStore string
	// Notes contains a cached list of Notes, but their content field may be nil
	// Use GetNote() to retrieve the complete note with content
	Notes []model.Note
}

// NewManager creates a new manager with the given store implementation
// and initializes the notes metadata list (without full content).
func NewManager(cloudStore *cloud.Client, sqliteStore *sqlite.Client, filesStore *files.Client, defaultToCloud bool) *Manager {
	var defaultStore Store
	defaultStore = sqliteStore
	if cloudStore != nil && defaultToCloud {
		defaultStore = cloudStore
	}

	var stores []Store
	if cloudStore != nil {
		stores = append(stores, cloudStore)
	}
	if filesStore != nil {
		stores = append(stores, filesStore)
	}
	if sqliteStore != nil {
		stores = append(stores, sqliteStore)
	}

	if len(stores) == 0 {
		panic("No store found")
	}

	return &Manager{
		activeStore: defaultStore,
		Name:        defaultStore.Name(),
		stores:      stores,
	}
}

// UpdateCloudStore will swap the cloud storage for a new store
// Expect the store to be valid and functionnal, panic otherwise
func (m *Manager) UpdateCloudClient(client *cloud.Client) {
	found := false
	for i, store := range m.stores {
		if store.Name() == cloud.Name {
			m.stores[i] = client
			found = true
			break
		}
	}
	if !found {
		m.stores = append(m.stores, client)
		m.setActiveStore(client)
	}

	if m.activeStore.Name() == cloud.Name {
		m.activeStore = client
		_, err := m.ListNoteMetadata()
		if err != nil {
			log.Fatalf("Failed to fetch from new cloud store: %v", err)
		}
	}
}

// NextStore swap the current underlying storage with the next registered one
// Dev needs to call ListNoteMetadata after calling this, otherwise a panic occur
func (m *Manager) NextStore() error {
	for i, store := range m.stores {
		if store.Name() == m.activeStore.Name() {
			m.setActiveStore(m.stores[(i+1)%len(m.stores)])
			break
		}
	}

	return nil
}

func (m *Manager) setActiveStore(store Store) {
	m.activeStore = store
	m.Name = m.activeStore.Name()
	log.Infof("Set active store to %s", m.Name)
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
