package local

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid" // Added
	_ "github.com/mattn/go-sqlite3"
	"merlion/internal/model"
	"merlion/internal/store"
)

// Ensure LocalClient implements the store.Store interface.
var _ store.Store = (*LocalClient)(nil)

type LocalClient struct {
	db   *sql.DB
	name string
}

func NewLocalClient(db *sql.DB, name string) *LocalClient {
	return &LocalClient{
		db:   db,
		name: name,
	}
}

func (c *LocalClient) Name() string {
	return c.name
}

func (c *LocalClient) ListNotes() ([]model.Note, error) {
	rows, err := c.db.Query(`
		SELECT note_id, title, content, tags,
		       is_favorite, is_work_log, is_trash, created_at, updated_at
		FROM notes
	`)
	if err != nil {
		return nil, fmt.Errorf("failed to query notes: %w", err)
	}
	defer rows.Close()

	var notes []model.Note
	for rows.Next() {
		note, err := scanNote(rows)
		if err != nil {
			// TODO: Decide if one bad row should stop the whole list.
			// For now, it does. Consider collecting errors or skipping.
			return nil, fmt.Errorf("failed to scan note during list: %w", err)
		}
		notes = append(notes, *note)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("rows iteration error: %w", err)
	}

	return notes, nil
}

func (c *LocalClient) GetNote(noteID string) (*model.Note, error) {
	row := c.db.QueryRow(`
		SELECT note_id, title, content, tags,
		       is_favorite, is_work_log, is_trash, created_at, updated_at
		FROM notes WHERE note_id = ?
	`, noteID)

	note, err := scanNote(row)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, store.ErrNotFound
		}
		return nil, fmt.Errorf("failed to scan note: %w", err)
	}
	return note, nil
}

func (c *LocalClient) CreateNote(req model.CreateNoteRequest) (*model.Note, error) {
	noteID := uuid.New().String() // New NoteID generation
	now := time.Now()

	var tagsJSON []byte
	var err error
	if req.Tags == nil || len(req.Tags) == 0 {
		tagsJSON = []byte("[]")
	} else {
		tagsJSON, err = json.Marshal(req.Tags)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal tags: %w", err)
		}
	}

	isFavorite := false
	if req.IsFavorite != nil {
		isFavorite = *req.IsFavorite
	}
	isWorkLog := false
	if req.IsWorkLog != nil {
		isWorkLog = *req.IsWorkLog
	}
	isTrash := false // Default for IsTrash
	if req.IsTrash != nil {
		isTrash = *req.IsTrash
	}

	// title, content, tags, is_favorite, is_work_log, is_trash, created_at, updated_at
	stmt, err := c.db.Prepare(`
		INSERT INTO notes (
			note_id, title, content, tags,
			is_favorite, is_work_log, is_trash, created_at, updated_at
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)
	`)
	if err != nil {
		return nil, fmt.Errorf("failed to prepare statement for create: %w", err)
	}
	defer stmt.Close()

	_, err = stmt.Exec(
		noteID, req.Title, req.Content, string(tagsJSON),
		isFavorite, isWorkLog, isTrash, now, now,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to execute statement for create: %w", err)
	}

	return &model.Note{
		NoteID:     noteID,
		Title:      req.Title,
		Content:    req.Content,
		Tags:       req.Tags, // Return original tags (might be nil, client should handle)
		IsFavorite: isFavorite,
		IsWorkLog:  isWorkLog,
		IsTrash:    isTrash,
		CreatedAt:  now,
		UpdatedAt:  now,
	}, nil
}

func (c *LocalClient) UpdateNote(noteID string, req model.CreateNoteRequest) (*model.Note, error) {
	now := time.Now()

	// Build the update query dynamically based on provided fields
	// For simplicity, current schema has NOT NULL DEFAULTS, so we can update all relevant fields.
	// However, if we want true partial updates (only set fields that are non-nil in req),
	// the query construction would be more complex.
	// Given the model changes, we update all user-settable fields from CreateNoteRequest.

	var tagsJSON []byte
	var err error
	if req.Tags == nil || len(req.Tags) == 0 {
		tagsJSON = []byte("[]")
	} else {
		tagsJSON, err = json.Marshal(req.Tags)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal tags for update: %w", err)
		}
	}

	isFavorite := false // Default if not provided
	if req.IsFavorite != nil {
		isFavorite = *req.IsFavorite
	}
	isWorkLog := false // Default if not provided
	if req.IsWorkLog != nil {
		isWorkLog = *req.IsWorkLog
	}
	isTrash := false // Default if not provided
	if req.IsTrash != nil {
		isTrash = *req.IsTrash
	}

	stmt, err := c.db.Prepare(`
		UPDATE notes
		SET title = ?, content = ?, tags = ?,
		    is_favorite = ?, is_work_log = ?, is_trash = ?,
		    updated_at = ?
		WHERE note_id = ?
	`)
	if err != nil {
		return nil, fmt.Errorf("failed to prepare update statement: %w", err)
	}
	defer stmt.Close()

	res, err := stmt.Exec(
		req.Title, req.Content, string(tagsJSON),
		isFavorite, isWorkLog, isTrash,
		now,
		noteID,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to execute update statement: %w", err)
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return nil, fmt.Errorf("failed to get rows affected: %w", err)
	}
	if rowsAffected == 0 {
		return nil, store.ErrNotFound
	}

	// We need the original CreatedAt time. We should fetch it.
	// For simplicity now, we return what we have, but GetNote would be better.
	// However, GetNote doesn't return a Note with all fields from CreateNoteRequest.
	// Let's call GetNote(noteID) to ensure consistency.
	return c.GetNote(noteID)
}

func (c *LocalClient) DeleteNote(noteID string) error {
	stmt, err := c.db.Prepare(`DELETE FROM notes WHERE note_id = ?`)
	if err != nil {
		return fmt.Errorf("failed to prepare delete statement: %w", err)
	}
	defer stmt.Close()

	res, err := stmt.Exec(noteID)
	if err != nil {
		return fmt.Errorf("failed to execute delete statement: %w", err)
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return store.ErrNotFound
	}

	return nil
}

func scanNote(row interface{ Scan(...interface{}) error }) (*model.Note, error) {
	var note model.Note
	var tagsJSON string
	var content sql.NullString
	// workspaceID and isPublic removed

	err := row.Scan(
		&note.NoteID,
		&note.Title,
		&content,
		&tagsJSON,
		&note.IsFavorite,
		&note.IsWorkLog,
		&note.IsTrash, // Added
		&note.CreatedAt,
		&note.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	if content.Valid {
		note.Content = &content.String
	} else {
		note.Content = nil
	}
	// workspaceID logic removed

	if err := json.Unmarshal([]byte(tagsJSON), &note.Tags); err != nil {
		// Handle cases where tagsJSON might be empty or null in DB if that's possible
		if tagsJSON == "" || tagsJSON == "null" {
			note.Tags = []string{}
		} else {
			return nil, fmt.Errorf("failed to unmarshal tags: %w", err)
		}
	}
	return &note, nil
}
