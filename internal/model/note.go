package model

import "time"

type Note struct {
	NoteID      string    `json:"note_id"`
	Title       string    `json:"title"`
	Content     *string   `json:"content"`      // pointer for nullable string
	WorkspaceID *string   `json:"workspace_id"` // pointer for nullable UUID
	Tags        []string  `json:"tags"`
	IsFavorite  bool      `json:"is_favorite"`
	IsWorkLog   bool      `json:"is_work_log"`
	IsPublic    bool      `json:"is_public"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type CreateNoteRequest struct {
	Title       string   `json:"title"`
	Content     *string  `json:"content"`      // pointer for nullable string
	WorkspaceID *string  `json:"workspace_id"` // pointer for nullable UUID
	Tags        []string `json:"tags,omitempty"`
	IsFavorite  *bool    `json:"is_favorite,omitempty"`
	IsWorkLog   *bool    `json:"is_work_log,omitempty"`
	IsPublic    *bool    `json:"is_public,omitempty"`
}

// Utilities

func (n Note) ToCreateRequest() CreateNoteRequest {
	return CreateNoteRequest{
		Title:       n.Title,
		Content:     n.Content,
		WorkspaceID: n.WorkspaceID,
		Tags:        n.Tags,
		IsFavorite:  &n.IsFavorite,
		IsWorkLog:   &n.IsWorkLog,
		IsPublic:    &n.IsPublic,
	}
}
