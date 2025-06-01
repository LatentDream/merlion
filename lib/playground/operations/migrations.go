package operations

import (
	"cmp"
	"database/sql"
	"embed"
	"errors"
	"fmt"
	"io/fs"
	"slices"
	"strconv"
	"strings"

	"github.com/glebarez/go-sqlite"
	"github.com/latentdream/merlion/lib/playground/assert"
)

const INITIAL_VERSION = 0

//go:embed migrations/*.sql
var migrationsContent embed.FS

func tableExist(db *sql.DB, tableName string) (bool, error) {
	query := `SELECT count(*) FROM sqlite_master WHERE type = 'table' AND name = ?;`

	var count int
	err := db.QueryRow(query, tableName).Scan(&count)
	if err != nil {
		return false, err
	}

	return count > 0, nil
}

func createVersionTable(db *sql.DB) (int, error) {
	query__createVersionTable := `
			CREATE TABLE IF NOT EXISTS merlion_version (
				id string PRIMARY KEY,
				version INTEGER NOT NULL
			)`
	_, err := db.Exec(query__createVersionTable)
	if err != nil {
		return -1, err
	}

	query__initialStamp := `INSERT INTO merlion_version VALUES("version", ?)`
	_, err = db.Exec(query__initialStamp, INITIAL_VERSION)
	if err != nil {
		var sqliteErr *sqlite.Error
		if errors.As(err, &sqliteErr) && sqliteErr.Code() == 1555 {
			return INITIAL_VERSION, nil
		}
		return -1, err
	}
	return INITIAL_VERSION, nil
}

func GetVersion(db *sql.DB) (int, error) {
	exist, err := tableExist(db, "merlion_version")
	if err != nil {
		return -1, err
	}
	if !exist {
		return createVersionTable(db)
	}

	query__getCurrVersion := `SELECT version FROM merlion_version`
	var version int
	db.QueryRow(query__getCurrVersion).Scan(&version)

	return version, nil
}

func getMigrationVersion(migrationFile fs.DirEntry) int {
	parts := strings.Split(migrationFile.Name(), "__")
	assert.Eq(len(parts), 2, "File format should be int__description.sql, like: 024__adding_something.sql but found:", migrationFile.Name())
	numberStr := parts[0]
	number, err := strconv.Atoi(numberStr)
	if err != nil {
		panic(fmt.Sprintf("Failed to parse migration number for: %s. Expect 011__descriotion.sql", migrationFile.Name()))
	}
	return number
}

func GetMigrationFiles() []fs.DirEntry {

	entries, err := migrationsContent.ReadDir("migrations")
	if err != nil {
		panic(fmt.Sprintln("Error listing embedded migrations:", err))
	}

	fmt.Println("Files in the embedded migrations directory:")

	sortedFunc := func(entry1 fs.DirEntry, entry2 fs.DirEntry) int {
		return cmp.Compare(getMigrationVersion(entry1), getMigrationVersion(entry2))
	}

	filterFunc := func(entry fs.DirEntry) bool {
		if entry.IsDir() {
			return true
		}
		parts := strings.Split(entry.Name(), ".")
		assert.Eq(len(parts), 2, "File format should be filename.ext, found:", entry.Name())
		if parts[1] != "sql" {
			return true
		}
		return false
	}

	entries = slices.DeleteFunc(entries, filterFunc)
	slices.SortFunc(entries, sortedFunc)

	return entries
}

func ApplyMigrations(db *sql.DB) error {

	currVersion, err := GetVersion(db)
	if err != nil {
		panic(fmt.Sprintln("FATAL: impossible to get version -", err))
	}
	fmt.Printf("DB currently stamp at: #%d\n", currVersion)

	migrationFiles := GetMigrationFiles()
	nbApply := 0
	for _, migrationFile := range migrationFiles {
		version := getMigrationVersion(migrationFile)
		if version <= currVersion {
			continue
		}

		filePath := "migrations/" + migrationFile.Name()
		data, err := migrationsContent.ReadFile(filePath)
		if err != nil {
			return fmt.Errorf("failed to read migration file %s: %w", migrationFile.Name(), err)
		}

		fmt.Printf("Applying migration %s (version %d)\n", migrationFile.Name(), version)

		_, err = db.Exec(string(data))
		if err != nil {
			return fmt.Errorf("failed to execute migration %s: %w", migrationFile.Name(), err)
		}

		updateVersionQuery := `UPDATE merlion_version SET version = ? WHERE id = "version"`
		_, err = db.Exec(updateVersionQuery, version)
		if err != nil {
			return fmt.Errorf("failed to update version after migration %s: %w", migrationFile.Name(), err)
		}

		nbApply += 1
	}

	fmt.Printf("Applied %d migrations\n", nbApply)
	return nil
}
