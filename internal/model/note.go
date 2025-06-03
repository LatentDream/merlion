package model

import "time"

type Note struct {
	NoteID      string    `json:"note_id"`
	Title       string    `json:"title"`
	Content     *string   `json:"content"` // pointer for nullable string
	Tags        []string  `json:"tags"`
	IsFavorite  bool      `json:"is_favorite"`
	IsWorkLog   bool      `json:"is_work_log"`
	IsTrash     bool      `json:"is_trash"`    // New field
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	// WorkspaceID and IsPublic are removed
}

type CreateNoteRequest struct {
	Title      string   `json:"title"`
	Content    *string  `json:"content"` // pointer for nullable string
	Tags       []string `json:"tags,omitempty"`
	IsFavorite *bool    `json:"is_favorite,omitempty"`
	IsWorkLog  *bool    `json:"is_work_log,omitempty"`
	IsTrash    *bool    `json:"is_trash,omitempty"` // New field
	// WorkspaceID and IsPublic are removed
}

// Utilities

func (n Note) ToCreateRequest() CreateNoteRequest {
	return CreateNoteRequest{
		Title:      n.Title,
		Content:    n.Content,
		Tags:       n.Tags,
		IsFavorite: &n.IsFavorite,
		IsWorkLog:  &n.IsWorkLog,
		IsTrash:    &n.IsTrash,
	}
}
