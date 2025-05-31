package operations

import (
	"database/sql"
	"errors"

	"github.com/glebarez/go-sqlite"
)

const INITIAL_VERSION = 0

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
