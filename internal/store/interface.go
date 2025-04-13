package store

import (
	"merlion/internal/model"
)

// Encapsulation to introduce on device local store
type Store interface {
	ListNotes() ([]model.Note, error)
	UpdateNote(string, model.CreateNoteRequest) (*model.Note, error)
	GetTags() []string
	GetNote(noteID string) (*model.Note, error)
	CreateNote(req model.CreateNoteRequest) (*model.Note, error)
	DeleteNote(noteID string) error
}
