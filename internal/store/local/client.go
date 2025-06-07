package local

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
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

func NewLocalClient(db *sql.DB) *LocalClient {
	return &LocalClient{
		db:   db,
		name: "Local Storage",
	}
}

func (c *LocalClient) Name() string {
	return c.name
}

func (c *LocalClient) ListNotes() ([]model.Note, error) {
	rows, err := c.db.Query(`
		SELECT note_id, title, content, tags, is_favorite,
			   is_work_log, created_at, updated_at
		FROM notes
		WHERE is_trash = false
	`)
	if err != nil {
		return nil, fmt.Errorf("failed to query notes: %w", err)
	}
	defer rows.Close()

	var notes []model.Note
	for rows.Next() {
		note, err := scanNote(rows)
		if err != nil {
			// TODO: One bad row should stop the whole list.
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
		SELECT note_id, title, content, tags, is_favorite,
			   is_work_log, created_at, updated_at
		FROM notes WHERE note_id = ?
	`, noteID)

	note, err := scanNote(row)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, store.ErrNoteNotFound
		}
		return nil, fmt.Errorf("failed to scan note: %w", err)
	}
	return note, nil
}

func (c *LocalClient) CreateNote(req model.CreateNoteRequest) (*model.Note, error) {
	noteID := uuid.New().String()
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

	// title, content, tags, is_favorite, is_work_log, created_at, updated_at
	stmt, err := c.db.Prepare(`
		INSERT INTO notes (
			note_id, title, content, tags,
			is_favorite, is_work_log, created_at, updated_at
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?)
	`)
	if err != nil {
		return nil, fmt.Errorf("failed to prepare statement for create: %w", err)
	}
	defer stmt.Close()

	_, err = stmt.Exec(
		noteID, req.Title, req.Content, string(tagsJSON),
		isFavorite, isWorkLog, now, now,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to execute statement for create: %w", err)
	}

	return &model.Note{
		NoteID:     noteID,
		Title:      req.Title,
		Content:    req.Content,
		Tags:       req.Tags,
		IsFavorite: isFavorite,
		IsWorkLog:  isWorkLog,
		CreatedAt:  now,
		UpdatedAt:  now,
	}, nil
}

func (c *LocalClient) UpdateNote(noteID string, req model.CreateNoteRequest) (*model.Note, error) {
	now := time.Now()

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

	stmt, err := c.db.Prepare(`
		UPDATE notes
		SET title = ?, content = ?, tags = ?,
		    is_favorite = ?, is_work_log = ?,
		    updated_at = ?
		WHERE note_id = ?
	`)
	if err != nil {
		return nil, fmt.Errorf("failed to prepare update statement: %w", err)
	}
	defer stmt.Close()

	res, err := stmt.Exec(
		req.Title, req.Content, string(tagsJSON),
		isFavorite, isWorkLog,
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
		return nil, store.ErrNoteNotFound
	}

	return c.GetNote(noteID)
}

func (c *LocalClient) DeleteNote(noteID string) error {
	stmt, err := c.db.Prepare(`UPDATE notes SET is_trash = ? WHERE note_id = ?`)
	if err != nil {
		return fmt.Errorf("failed to prepare update statement: %w", err)
	}
	defer stmt.Close()

	res, err := stmt.Exec(true, noteID)
	if err != nil {
		return fmt.Errorf("failed to execute update statement: %w", err)
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return store.ErrNoteNotFound
	}

	return nil
}

func scanNote(row interface{ Scan(...interface{}) error }) (*model.Note, error) {
	var note model.Note
	var tagsJSON string
	var content sql.NullString

	err := row.Scan(
		&note.NoteID,
		&note.Title,
		&content,
		&tagsJSON,
		&note.IsFavorite,
		&note.IsWorkLog,
		&note.CreatedAt,
		&note.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	if content.Valid {
		note.Content = &content.String
	} else {
		emptyContent := ""
		note.Content = &emptyContent
	}

	if err := json.Unmarshal([]byte(tagsJSON), &note.Tags); err != nil {
		// Just to be sure
		if tagsJSON == "" || tagsJSON == "null" {
			note.Tags = []string{}
		} else {
			return nil, fmt.Errorf("failed to unmarshal tags: %w", err)
		}
	}
	return &note, nil
}
