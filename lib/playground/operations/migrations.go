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

func TableExist(db *sql.DB, tableName string) (bool, error) {
	query := `SELECT count(*) FROM sqlite_master WHERE type = 'name' AND name='?';`

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
	exist, err := TableExist(db, "version")
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

func ApplyMigrations() []fs.DirEntry {

	// List all files in the embedded migrations directory
	entries, err := migrationsContent.ReadDir("migrations")
	if err != nil {
		panic(fmt.Sprintln("Error listing embedded migrations:", err))
	}

	fmt.Println("Files in the embedded migrations directory:")

	sortedFunc := func(entry1 fs.DirEntry, entry2 fs.DirEntry) int {
		parse := func(entry fs.DirEntry) int {
			parts := strings.Split(entry.Name(), "__")
			assert.Eq(len(parts), 2, "File format should be int__description.sql, like: 024__adding_something.sql but found:", entry.Name())
			numberStr := parts[0]
			number, err := strconv.Atoi(numberStr)
			if err != nil {
				panic(fmt.Sprintf("Failed to parse migration number for: %s. Expect 011__descriotion.sql", entry.Name()))
			}
			return number
		}
		return cmp.Compare(parse(entry1), parse(entry2))
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
