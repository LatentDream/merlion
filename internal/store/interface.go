package store

import (
	"merlion/internal/model"

	"github.com/charmbracelet/log"
)

// Encapsulation to introduce on device local store
type Store interface {
	ListNotes() ([]model.Note, error)
	UpdateNote(noteId string, req model.CreateNoteRequest) *model.Note
	GetTags() []string
	GetNote(noteID string) (*model.Note, error)
	CreateNote(req model.CreateNoteRequest) (*model.Note, error)
	DeleteNote(noteID string) error
}

type Manager struct {
	activeStore Store
	notes       []model.Note
}

func NewManager(store Store) *Manager {
	notes, err := store.ListNotes()
	if err != nil {
		log.Fatalf("Not able to list notes at the moment")
	}
	return &Manager{
		activeStore: store,
		notes:       notes,
	}
}
