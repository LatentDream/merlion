package store

import (
	"merlion/internal/model"
)

// Encapsulation to introduce on device local store
type store interface {
	Name() string
	ListNotes() ([]model.Note, error)
	UpdateNote(string, model.CreateNoteRequest) (*model.Note, error)
	GetNote(noteID string) (*model.Note, error)
	CreateNote(req model.CreateNoteRequest) (*model.Note, error)
	DeleteNote(noteID string) error
}
