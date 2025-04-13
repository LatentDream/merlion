package store

import (
	"log"
	"merlion/internal/model"
)

type Manager struct {
	activeStore Store
	name        string
	notes       []model.Note
}

func NewManager(store Store) *Manager {
	notes, err := store.ListNotes()
	if err != nil {
		log.Fatalf("Not able to list notes at the moment")
	}
	return &Manager{
		activeStore: store,
		name:        store.Name(),
		notes:       notes,
	}
}

func (m *Manager) GetNote(noteID string) (*model.Note, error) {
	return m.activeStore.GetNote(noteID)
}

func (m *Manager) ListNotes() ([]model.Note, error) {
	return m.activeStore.ListNotes()
}
func (m *Manager) UpdateNote(noteId string, changes model.CreateNoteRequest) (*model.Note, error) {
	return m.activeStore.UpdateNote(noteId, changes)
}
func (m *Manager) GetTags() []string {
	return m.activeStore.GetTags()
}
func (m *Manager) CreateNote(req model.CreateNoteRequest) (*model.Note, error) {
	return m.activeStore.CreateNote(req)
}
func (m *Manager) DeleteNote(noteID string) error {
	return m.activeStore.DeleteNote(noteID)
}
