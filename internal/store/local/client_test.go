package local

import (
	"database/sql"
	"encoding/json"
	"os"
	"testing"
	"time"

	"github.com/google/uuid" // Added for NoteID validation
	"merlion/internal/model"
	"merlion/internal/store"
	"merlion/internal/store/local/operations" // For migrations

	_ "github.com/mattn/go-sqlite3" // SQLite driver
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// setupTestDB initializes an in-memory SQLite database and applies migrations.
// It returns the database connection and a teardown function.
func setupTestDB(t *testing.T) (*sql.DB, func()) {
	t.Helper() // Marks this function as a test helper

	// Open an in-memory SQLite database.
	// ":memory:" is a special DSN for in-memory SQLite.
	// Options like ?_foreign_keys=on can be added if needed.
	db, err := sql.Open("sqlite3", ":memory:?_foreign_keys=on")
	require.NoError(t, err, "Failed to open in-memory database")

	// Apply migrations
	// The migrations path needs to be correctly located if operations.ApplyMigrations expects it.
	// Assuming operations.ApplyMigrations can work with the *sql.DB directly without needing a path for embedded migrations.
	// If ApplyMigrations needs a path to SQL files, this setup will need adjustment
	// (e.g., by copying migrations to a temp dir or using a relative path if tests run from a specific location).
	// For now, assuming ApplyMigrations uses embedded migrations or a mechanism that doesn't rely on file paths from os.Getwd().
	// The actual ApplyMigrations function in the codebase will determine this.
	// Let's assume it's: err = operations.ApplyMigrations(db)
	err = operations.ApplyMigrations(db)
	require.NoError(t, err, "Failed to apply migrations")

	teardown := func() {
		err := db.Close()
		assert.NoError(t, err, "Failed to close database") // Use assert for teardown, require might stop other teardowns
	}

	return db, teardown
}

// TestMain can be used for package-level setup/teardown if needed,
// e.g., setting up a logger or other global states for tests.
func TestMain(m *testing.M) {
	// Example: log.SetOutput(ioutil.Discard) to suppress logs during tests
	os.Exit(m.Run())
}

// Placeholder for the first test
func TestLocalClient_CreateAndGetNote(t *testing.T) {
	db, teardown := setupTestDB(t)
	defer teardown()

	client := NewLocalClient(db, "test-sqlite")
	require.NotNil(t, client)
	assert.Equal(t, "test-sqlite", client.Name())

	// 1. Create a note
	contentStr := "This is the note content."
	trueBool := true
	falseBool := false

	createReq := model.CreateNoteRequest{
		Title:      "Test Note 1",
		Content:    &contentStr,
		Tags:       []string{"tag1", "tag2"},
		IsFavorite: &trueBool,
		IsWorkLog:  &falseBool,
		IsTrash:    &falseBool, // New field
	}

	createdNote, err := client.CreateNote(createReq)
	require.NoError(t, err)
	require.NotNil(t, createdNote)

	// Assertions for the created note
	_, err = uuid.Parse(createdNote.NoteID) // Verify NoteID is a valid UUID
	assert.NoError(t, err, "NoteID should be a valid UUID")
	assert.Equal(t, createReq.Title, createdNote.Title)
	require.NotNil(t, createdNote.Content)
	assert.Equal(t, *createReq.Content, *createdNote.Content)
	// WorkspaceID removed
	assert.ElementsMatch(t, createReq.Tags, createdNote.Tags) // Use ElementsMatch for slices
	assert.Equal(t, *createReq.IsFavorite, createdNote.IsFavorite)
	assert.Equal(t, *createReq.IsWorkLog, createdNote.IsWorkLog)
	// IsPublic removed
	assert.Equal(t, *createReq.IsTrash, createdNote.IsTrash) // Verify IsTrash

	assert.WithinDuration(t, time.Now(), createdNote.CreatedAt, 2*time.Second, "CreatedAt should be recent")
	assert.WithinDuration(t, time.Now(), createdNote.UpdatedAt, 2*time.Second, "UpdatedAt should be recent")
	assert.Equal(t, createdNote.CreatedAt, createdNote.UpdatedAt, "CreatedAt and UpdatedAt should be same on creation")

	// Synced and RevisionID removed

	// 2. Get the note
	retrievedNote, err := client.GetNote(createdNote.NoteID)
	require.NoError(t, err)
	require.NotNil(t, retrievedNote)
	assert.Equal(t, *createdNote, *retrievedNote, "Retrieved note should match created note")

	// Specifically check Tags deserialization
	assert.ElementsMatch(t, createReq.Tags, retrievedNote.Tags, "Tags should match")


	// 3. Test GetNote for a non-existent ID
	_, err = client.GetNote("non-existent-id")
	assert.ErrorIs(t, err, store.ErrNotFound, "Expected ErrNotFound for non-existent note ID")
}

func TestLocalClient_ListNotes(t *testing.T) {
	db, teardown := setupTestDB(t)
	defer teardown()

	client := NewLocalClient(db, "test-sqlite")
	require.NotNil(t, client)

	// 1. List notes on an empty database
	notes, err := client.ListNotes()
	require.NoError(t, err)
	assert.Empty(t, notes, "Expected no notes in an empty database")

	// 2. Create a couple of notes
	content1 := "content1"
	content2 := "content2"
	createReq1 := model.CreateNoteRequest{Title: "Note A", Content: &content1}
	createReq2 := model.CreateNoteRequest{Title: "Note B", Content: &content2}

	noteA, err := client.CreateNote(createReq1)
	require.NoError(t, err)
	require.NotNil(t, noteA)

	noteB, err := client.CreateNote(createReq2)
	require.NoError(t, err)
	require.NotNil(t, noteB)

	// 3. List notes again
	notes, err = client.ListNotes()
	require.NoError(t, err)
	assert.Len(t, notes, 2, "Expected two notes")

	// Verify the content (shallow check)
	expectedNotes := map[string]string{
		noteA.NoteID: noteA.Title,
		noteB.NoteID: noteB.Title,
	}
	actualNotes := make(map[string]string)
	for _, n := range notes {
		actualNotes[n.NoteID] = n.Title
	}
	assert.Equal(t, expectedNotes, actualNotes, "Listed notes do not match expected")
}

func TestLocalClient_UpdateNote(t *testing.T) {
	db, teardown := setupTestDB(t)
	defer teardown()

	client := NewLocalClient(db, "test-sqlite")
	require.NotNil(t, client)

	// 1. Create an initial note
	initialContent := "Initial content"
	initialTags := []string{"initial"}
	trueBool := true
	falseBool := false
	initialReq := model.CreateNoteRequest{
		Title:      "Initial Title",
		Content:    &initialContent,
		Tags:       initialTags,
		IsFavorite: &trueBool,
		IsTrash:    &falseBool,
	}
	createdNote, err := client.CreateNote(initialReq)
	require.NoError(t, err)
	require.NotNil(t, createdNote)
	originalUpdatedAt := createdNote.UpdatedAt
	originalNoteID := createdNote.NoteID // Store original NoteID to check it doesn't change

	// Introduce a small delay to ensure UpdatedAt timestamp can change
	time.Sleep(50 * time.Millisecond)

	// 2. Create an update request
	updatedContent := "Updated Content"
	updatedTags := []string{"initial", "updated"}
	updatedIsTrash := true
	updateReq := model.CreateNoteRequest{
		Title:      "Updated Title",
		Content:    &updatedContent,
		Tags:       updatedTags,
		IsFavorite: &falseBool,       // Changed value
		IsTrash:    &updatedIsTrash, // Changed value
		// IsWorkLog will be default false as it's not provided in updateReq
	}

	// 3. Call UpdateNote
	updatedNote, err := client.UpdateNote(createdNote.NoteID, updateReq)
	require.NoError(t, err)
	require.NotNil(t, updatedNote)

	// 4. Verify UpdatedAt and that NoteID did not change
	assert.True(t, updatedNote.UpdatedAt.After(originalUpdatedAt), "UpdatedAt should be newer")
	assert.Equal(t, originalNoteID, updatedNote.NoteID, "NoteID should not change on update")
	// RevisionID and Synced removed

	// 5. Retrieve the note and assert all fields
	retrievedNote, err := client.GetNote(createdNote.NoteID)
	require.NoError(t, err)
	require.NotNil(t, retrievedNote)

	assert.Equal(t, updateReq.Title, retrievedNote.Title)
	require.NotNil(t, retrievedNote.Content)
	assert.Equal(t, *updateReq.Content, *retrievedNote.Content)
	assert.ElementsMatch(t, updateReq.Tags, retrievedNote.Tags)
	assert.Equal(t, *updateReq.IsFavorite, retrievedNote.IsFavorite)
	// IsPublic removed
	assert.Equal(t, *updateReq.IsTrash, retrievedNote.IsTrash) // Verify IsTrash
	assert.Equal(t, createdNote.CreatedAt, retrievedNote.CreatedAt, "CreatedAt should not change on update")
	assert.Equal(t, updatedNote.UpdatedAt, retrievedNote.UpdatedAt)

	// IsWorkLog was not in initialReq or updateReq, should default to false in the DB.
	// The client.UpdateNote currently sets it to false if not provided in req.
	// So, retrievedNote.IsWorkLog should be false.
	assert.False(t, retrievedNote.IsWorkLog, "IsWorkLog should be false as it was not set in update")


	// 6. Test UpdateNote for a non-existent ID
	_, err = client.UpdateNote("non-existent-id", updateReq)
	assert.ErrorIs(t, err, store.ErrNotFound, "Expected ErrNotFound for updating non-existent note ID")
}

func TestLocalClient_DeleteNote(t *testing.T) {
	db, teardown := setupTestDB(t)
	defer teardown()

	client := NewLocalClient(db, "test-sqlite")
	require.NotNil(t, client)

	// 1. Create a note
	content := "Content to be deleted"
	createReq := model.CreateNoteRequest{Title: "Delete Me", Content: &content}
	createdNote, err := client.CreateNote(createReq)
	require.NoError(t, err)
	require.NotNil(t, createdNote)

	// 2. Delete the note
	err = client.DeleteNote(createdNote.NoteID)
	require.NoError(t, err, "Expected no error when deleting an existing note")

	// 3. Try to Get the deleted note
	_, err = client.GetNote(createdNote.NoteID)
	assert.ErrorIs(t, err, store.ErrNotFound, "Expected ErrNotFound after deleting the note")

	// 4. Test DeleteNote for a non-existent ID
	err = client.DeleteNote("non-existent-id")
	assert.ErrorIs(t, err, store.ErrNotFound, "Expected ErrNotFound when trying to delete non-existent note ID")

	// 5. Test DeleteNote for an already deleted ID
	err = client.DeleteNote(createdNote.NoteID)
	assert.ErrorIs(t, err, store.ErrNotFound, "Expected ErrNotFound when trying to delete an already deleted note ID")
}

func TestLocalClient_TagsSerialization(t *testing.T) {
	db, teardown := setupTestDB(t)
	defer teardown()

	client := NewLocalClient(db, "test-sqlite")
	require.NotNil(t, client)

	// 1. Create a note with multiple tags
	tags1 := []string{"apple", "banana", "cherry"}
	createReq1 := model.CreateNoteRequest{Title: "Tags Test 1", Tags: tags1}
	note1, err := client.CreateNote(createReq1)
	require.NoError(t, err)
	retrieved1, err := client.GetNote(note1.NoteID)
	require.NoError(t, err)
	assert.ElementsMatch(t, tags1, retrieved1.Tags, "Failed to retrieve multiple tags")

	// 2. Create a note with empty tags (from request)
	var tags2 []string // empty, not nil
	createReq2 := model.CreateNoteRequest{Title: "Tags Test 2", Tags: tags2}
	note2, err := client.CreateNote(createReq2)
	require.NoError(t, err)
	retrieved2, err := client.GetNote(note2.NoteID)
	require.NoError(t, err)
	assert.Equal(t, []string{}, retrieved2.Tags, "Expected empty slice for empty tags, got nil or non-empty")
	// The scanNote function in client.go correctly initializes Tags to []string{} if JSON is null or empty string.

	// 3. Create a note with nil tags (from request)
	createReqNil := model.CreateNoteRequest{Title: "Tags Test Nil", Tags: nil}
	noteNil, err := client.CreateNote(createReqNil)
	require.NoError(t, err)
	retrievedNil, err := client.GetNote(noteNil.NoteID)
	require.NoError(t, err)
	assert.Equal(t, []string{}, retrievedNil.Tags, "Expected empty slice for nil tags from request, got non-empty")


	// 4. Update a note to have different tags
	tags3 := []string{"dog", "elephant"}
	updateReq3 := model.CreateNoteRequest{Title: "Tags Test 1 Updated", Tags: tags3}
	_, err = client.UpdateNote(note1.NoteID, updateReq3)
	require.NoError(t, err)
	retrieved3, err := client.GetNote(note1.NoteID)
	require.NoError(t, err)
	assert.ElementsMatch(t, tags3, retrieved3.Tags, "Failed to update tags")

	// 5. Update a note to have empty tags
	var tags4 []string
	updateReq4 := model.CreateNoteRequest{Title: "Tags Test 1 Emptied", Tags: tags4}
	_, err = client.UpdateNote(note1.NoteID, updateReq4)
	require.NoError(t, err)
	retrieved4, err := client.GetNote(note1.NoteID)
	require.NoError(t, err)
	assert.Equal(t, []string{}, retrieved4.Tags, "Expected empty slice for updated empty tags")

	// 6. Update a note to have nil tags (from request)
	updateReqNil := model.CreateNoteRequest{Title: "Tags Test 1 Nil Updated", Tags: nil}
	_, err = client.UpdateNote(note1.NoteID, updateReqNil)
	require.NoError(t, err)
	retrievedNilUpdate, err := client.GetNote(note1.NoteID)
	require.NoError(t, err)
	assert.Equal(t, []string{}, retrievedNilUpdate.Tags, "Expected empty slice for updated nil tags")
}

func TestLocalClient_ContentNullability(t *testing.T) { // Renamed
	db, teardown := setupTestDB(t)
	defer teardown()

	client := NewLocalClient(db, "test-sqlite")
	require.NotNil(t, client)

	// 1. Create with nil Content
	createReqNil := model.CreateNoteRequest{Title: "Nil Content"} // Content is nil by default
	noteNil, err := client.CreateNote(createReqNil)
	require.NoError(t, err)
	retrievedNil, err := client.GetNote(noteNil.NoteID)
	require.NoError(t, err)
	assert.Nil(t, retrievedNil.Content, "Expected nil Content when created with nil")
	// WorkspaceID parts removed

	// 2. Create with empty string Content
	emptyStr := ""
	createReqEmpty := model.CreateNoteRequest{Title: "Empty String Content", Content: &emptyStr}
	noteEmpty, err := client.CreateNote(createReqEmpty)
	require.NoError(t, err)
	retrievedEmpty, err := client.GetNote(noteEmpty.NoteID)
	require.NoError(t, err)
	require.NotNil(t, retrievedEmpty.Content, "Content should not be nil when created with empty string")
	assert.Equal(t, "", *retrievedEmpty.Content, "Expected empty string Content")
	// WorkspaceID parts removed

	// 3. Update existing note (retrievedEmpty) to have nil Content
	updateReqToNil := model.CreateNoteRequest{Title: "Updated To Nil Content", Content: nil}
	_, err = client.UpdateNote(retrievedEmpty.NoteID, updateReqToNil)
	require.NoError(t, err)
	updatedToNilNote, err := client.GetNote(retrievedEmpty.NoteID)
	require.NoError(t, err)
	assert.Nil(t, updatedToNilNote.Content, "Expected nil Content after update to nil")
	assert.Equal(t, updateReqToNil.Title, updatedToNilNote.Title) // Ensure title updated

	// 4. Update note with nil fields (retrievedNil) to have actual Content value
	newContent := "New Content"
	updateReqToVal := model.CreateNoteRequest{Title: "Updated To Value Content", Content: &newContent}
	_, err = client.UpdateNote(retrievedNil.NoteID, updateReqToVal)
	require.NoError(t, err)
	updatedToValNote, err := client.GetNote(retrievedNil.NoteID)
	require.NoError(t, err)
	require.NotNil(t, updatedToValNote.Content)
	assert.Equal(t, newContent, *updatedToValNote.Content)
	assert.Equal(t, updateReqToVal.Title, updatedToValNote.Title)
}
