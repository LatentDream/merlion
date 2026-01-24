// Package files implements a local storage for notes.
// Uses the same format as Obsidian (https://obsidian.md/)
package files

import (
	"encoding/json"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/charmbracelet/log"
	"merlion/internal/model"
)

type Client struct {
	root string
}

func NewClient(root string) (*Client, error) {
	baseFolder, err := validatePath(root)
	if err != nil {
		return nil, err
	}

	err = ensureDirectoryExists(baseFolder)
	if err != nil {
		return nil, err
	}

	return &Client{baseFolder}, nil
}

func (c *Client) Name() string {
	return "File System"
}

func (c *Client) ListNotes() ([]model.Note, error) {
	log.Debug("Starting to list notes", "root", c.root)
	var notes []model.Note

	err := filepath.WalkDir(c.root, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			log.Error("Error walking directory", "path", path, "error", err)
			return err
		}

		if d.IsDir() {
			return nil
		}

		if strings.ToLower(filepath.Ext(path)) != ".md" {
			return nil
		}

		note, err := c.parseNoteFile(path)
		if err != nil {
			log.Error("Error parsing note file", "path", path, "error", err)
			return err
		}

		notes = append(notes, *note)
		return nil
	})

	log.Debug("Finished listing notes", "count", len(notes), "error", err)
	return notes, err
}

func (c *Client) GetNote(noteID string) (*model.Note, error) {
	notePath := filepath.Join(c.root, noteID+".md")
	log.Debug("Getting note", "noteID", noteID, "path", notePath)

	if _, err := os.Stat(notePath); os.IsNotExist(err) {
		log.Warn("Note not found", "noteID", noteID, "path", notePath)
		return nil, fmt.Errorf("note not found: %s", noteID)
	}

	note, err := c.parseNoteFile(notePath)
	if err != nil {
		log.Error("Error parsing note", "noteID", noteID, "error", err)
	} else {
		log.Debug("Successfully retrieved note", "noteID", noteID, "title", note.Title)
	}
	return note, err
}

func (c *Client) CreateNote(req model.CreateNoteRequest) (*model.Note, error) {
	log.Debug("Creating note", "title", req.Title, "workspaceID", req.WorkspaceID, keyTags, req.Tags)

	// Verify that the title is valid
	forbiddenChars := []string{"/", "\\", ":", "*", "?", "\"", "<", ">", "|"}
	for _, char := range forbiddenChars {
		if strings.Contains(req.Title, char) {
			log.Error("Invalid title with forbidden characters", "title", req.Title, "char", char)
			return nil, fmt.Errorf("invalid title: %s", req.Title)
		}
	}

	noteID := req.Title
	notePath := filepath.Join(c.root, noteID+".md")
	log.Debug("Checking if note already exists", "noteID", noteID, "path", notePath)

	if _, err := os.Stat(notePath); err == nil {
		log.Warn("Note already exists", "noteID", noteID)
		return nil, fmt.Errorf("note already exists: %s", noteID)
	}

	now := time.Now()
	createdAt := now
	updatedAt := now
	if req.CreatedAt != nil {
		createdAt = *req.CreatedAt
	}
	if req.UpdatedAt != nil {
		updatedAt = *req.UpdatedAt
	}

	note := model.Note{
		NoteID:      noteID,
		Title:       req.Title,
		Content:     req.Content,
		WorkspaceID: req.WorkspaceID,
		Tags:        req.Tags,
		IsFavorite:  getBoolOrDefault(req.IsFavorite, false),
		IsWorkLog:   getBoolOrDefault(req.IsWorkLog, false),
		IsPublic:    getBoolOrDefault(req.IsPublic, false),
		CreatedAt:   createdAt,
		UpdatedAt:   updatedAt,
	}

	err := c.writeNoteFile(notePath, note)
	if err != nil {
		log.Error("Failed to write note file", "noteID", noteID, "error", err)
		return nil, err
	}

	log.Info("Successfully created note", "noteID", noteID, "title", note.Title)
	return &note, nil
}

func (c *Client) UpdateNote(noteID string, req model.CreateNoteRequest) (*model.Note, error) {
	notePath := filepath.Join(c.root, noteID+".md")
	log.Debug("Updating note", "noteID", noteID, "path", notePath)

	existingNote, err := c.parseNoteFile(notePath)
	if err != nil {
		log.Error("Failed to parse existing note for update", "noteID", noteID, "error", err)
		return nil, err
	}

	updatedNote := *existingNote
	updatedNote.Title = req.Title
	updatedNote.Content = req.Content
	updatedNote.WorkspaceID = req.WorkspaceID
	updatedNote.Tags = req.Tags
	updatedNote.IsFavorite = getBoolOrDefault(req.IsFavorite, existingNote.IsFavorite)
	updatedNote.IsWorkLog = getBoolOrDefault(req.IsWorkLog, existingNote.IsWorkLog)
	updatedNote.IsPublic = getBoolOrDefault(req.IsPublic, existingNote.IsPublic)
	updatedNote.UpdatedAt = time.Now()

	err = c.writeNoteFile(notePath, updatedNote)
	if err != nil {
		log.Error("Failed to write updated note file", "noteID", noteID, "error", err)
		return nil, err
	}

	return &updatedNote, nil
}

func (c *Client) DeleteNote(noteID string) error {
	notePath := filepath.Join(c.root, noteID+".md")
	log.Debug("Deleting note", "noteID", noteID, "path", notePath)

	if _, err := os.Stat(notePath); os.IsNotExist(err) {
		log.Warn("Note not found for deletion", "noteID", noteID, "path", notePath)
		return fmt.Errorf("note not found: %s", noteID)
	}

	err := moveToTrash(notePath)
	if err != nil {
		log.Error("Failed to move note to trash", "noteID", noteID, "error", err)
	} else {
		log.Info("Successfully deleted note", "noteID", noteID)
	}
	return err
}

func (c *Client) parseNoteFile(path string) (*model.Note, error) {
	content, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read note file: %w", err)
	}

	fileContent := string(content)
	frontMatter, noteContent, err := splitFrontMatterContent(fileContent)
	if err != nil {
		return nil, fmt.Errorf("failed to split front matter: %w", err)
	}

	filename, err := filepath.Rel(c.root, path)
	if err != nil {
		log.Error("Error getting relative path", "path", path, "error", err)
		return nil, err
	}
	noteID := strings.TrimSuffix(filename, filepath.Ext(filename))

	osCreatedTime, osUpdatedTime, err := getFileTimes(path)
	if err != nil {
		return nil, fmt.Errorf("failed to get file times: %w", err)
	}

	tags := frontMatterGetList(frontMatter, keyTags)
	isFavorite := frontMatterGetBool(frontMatter, keyIsFavorite, false)
	isWorkLog := frontMatterGetBool(frontMatter, keyIsWorkLog, false)
	createdTime := frontMatterGetTime(frontMatter, keyCreatedAt, osCreatedTime)
	updatedTime := frontMatterGetTime(frontMatter, keyUpdatedAt, osUpdatedTime)

	note := model.Note{
		NoteID:     noteID,
		Title:      noteID,
		Content:    &noteContent, // We could return a nil -> App will call GetNote to get the content (optional readContent bool arg)
		Tags:       tags,
		IsFavorite: isFavorite,
		IsWorkLog:  isWorkLog,
		IsPublic:   false,
		CreatedAt:  createdTime,
		UpdatedAt:  updatedTime,
	}

	if wsID, ok := frontMatter[keyWorkspace].(string); ok && wsID != "" {
		note.WorkspaceID = &wsID
	}

	return &note, nil
}

func (c *Client) writeNoteFile(path string, note model.Note) error {
	frontMatter := map[string]any{
		keyTags:       note.Tags,
		keyIsFavorite: note.IsFavorite,
		keyIsWorkLog:  note.IsWorkLog,
		keyCreatedAt:  note.CreatedAt.Format(time.RFC3339),
		keyUpdatedAt:  note.UpdatedAt.Format(time.RFC3339),
	}

	if note.WorkspaceID != nil {
		frontMatter[keyWorkspace] = *note.WorkspaceID
	}

	fmBytes, err := json.MarshalIndent(frontMatter, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal front matter: %w", err)
	}

	var fileContent strings.Builder
	fileContent.WriteString("---\n")
	fileContent.Write(fmBytes)
	fileContent.WriteString("\n---\n")

	if note.Content != nil {
		fileContent.WriteString(*note.Content)
	}

	return os.WriteFile(path, []byte(fileContent.String()), 0o644)
}
