package database

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetDBPath_Default(t *testing.T) {
	// Ensure MERLION_DB_PATH is unset for this test.
	// t.Setenv is available from Go 1.17. If using older, manual os.Unsetenv and os.Setenv is needed.
	originalEnv, isSet := os.LookupEnv("MERLION_DB_PATH")
	if isSet {
		os.Unsetenv("MERLION_DB_PATH")
		defer os.Setenv("MERLION_DB_PATH", originalEnv) // Restore after test
	}

	path, err := GetDBPath()
	require.NoError(t, err)

	homeDir, err := os.UserHomeDir()
	require.NoError(t, err, "Failed to get user home directory for assertion")
	expectedPath := filepath.Join(homeDir, ".merlion", "notes.db")

	assert.Equal(t, expectedPath, path, "Default DB path does not match expected")

	// Check if the directory was created
	dbDir := filepath.Dir(expectedPath)
	_, err = os.Stat(dbDir)
	assert.NoError(t, err, "Database directory should have been created: %s", dbDir)
	// Basic cleanup: remove .merlion/notes.db and .merlion if empty.
	// More robust cleanup might be needed if other tests rely on this dir.
	// For now, GetDBPath creates it, so it's a side effect we check.
	// Consider removing the directory if it's safe for other tests or if tests run in isolated environments.
}

func TestGetDBPath_EnvVar(t *testing.T) {
	customPath := filepath.Join(t.TempDir(), "custom_merlion_test.db")
	// t.TempDir() automatically creates a temporary directory for the test and cleans it up.

	t.Setenv("MERLION_DB_PATH", customPath)

	path, err := GetDBPath()
	require.NoError(t, err)
	assert.Equal(t, customPath, path, "DB path from env var does not match expected")

	// Check if the directory for the custom path was created (t.TempDir() creates the base)
	// GetDBPath should ensure the immediate parent dir of customPath exists.
	// Since customPath is directly in t.TempDir(), its parent is t.TempDir() itself, which is already created.
	// If customPath was /tmp/newdir/file.db, then newdir should be created.
	// In this case, filepath.Dir(customPath) is t.TempDir().
	dbDir := filepath.Dir(customPath)
	_, err = os.Stat(dbDir)
	assert.NoError(t, err, "Custom database directory should exist (created by t.TempDir): %s", dbDir)
}

func TestGetDBPath_EnvVarWithDirCreation(t *testing.T) {
	// Create a path within a non-existent subdirectory of t.TempDir()
	baseTmpDir := t.TempDir()
	nonExistentSubDir := filepath.Join(baseTmpDir, "nonexistentdir")
	customPath := filepath.Join(nonExistentSubDir, "custom_merlion_test.db")

	t.Setenv("MERLION_DB_PATH", customPath)

	path, err := GetDBPath()
	require.NoError(t, err)
	assert.Equal(t, customPath, path, "DB path from env var with subdir does not match expected")

	// Check if the non-existent subdirectory was created by GetDBPath
	_, err = os.Stat(nonExistentSubDir)
	assert.NoError(t, err, "The subdirectory specified in MERLION_DB_PATH should have been created: %s", nonExistentSubDir)
}

// TestInitDB is more of an integration test.
// For now, we are focusing on GetDBPath.
// A simple InitDB test could check if it returns a non-nil DB object and no error
// when using an in-memory database, and that migrations are attempted.
/*
func TestInitDB_InMemory(t *testing.T) {
	// Override GetDBPath for this test to return an in-memory DSN
	// This is tricky without DI for GetDBPath in InitDB.
	// Alternatively, set MERLION_DB_PATH to ":memory:" if GetDBPath supports it (it currently doesn't, expects file path).

	// For a real test of InitDB, you'd likely want to:
	// 1. Set up MERLION_DB_PATH to a temporary file.
	// 2. Call InitDB().
	// 3. Check if the db file was created.
	// 4. Check if migrations ran (e.g., by querying sqlite_master for expected tables).
	// 5. Close and remove the temp db file.

	// Example using a temp file:
	tmpFile := filepath.Join(t.TempDir(), "test_init.db")
	t.Setenv("MERLION_DB_PATH", tmpFile)

	db, err := InitDB()
	require.NoError(t, err)
	require.NotNil(t, db)
	defer db.Close()

	// Verify migrations ran - check for 'notes' table
	var tableName string
	queryErr := db.QueryRow("SELECT name FROM sqlite_master WHERE type='table' AND name='notes';").Scan(&tableName)
	require.NoError(t, queryErr, "Failed to query for notes table, migrations might not have run.")
	assert.Equal(t, "notes", tableName, "'notes' table should exist after migrations.")
}
*/
// Note: The TestInitDB_InMemory is commented out as per instructions to focus on GetDBPath.
// The actual implementation of such a test would require careful handling of GetDBPath or using
// a temporary database file. The current GetDBPath also creates directories, which isn't applicable for ":memory:".
// The example provided within the commented block shows a file-based approach for testing InitDB.
// The logic in internal/database/sqlite.go for GetDBPath was also updated to create the dir for MERLION_DB_PATH.
// This makes the TestGetDBPath_EnvVarWithDirCreation more relevant.
// The TestGetDBPath_Default's cleanup of ~/.merlion might be too aggressive if other things use this dir.
// For isolated test runs (like in CI or with t.TempDir), this is less of an issue.
// For local dev, care should be taken. I will remove the cleanup for the default path for now.
// The creation check is sufficient.
