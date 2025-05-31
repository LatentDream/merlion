package main

import (
	"database/sql"
	"fmt"

	_ "github.com/glebarez/go-sqlite"
	"github.com/latentdream/merlion/lib/playground/operations"
)

func main() {
	// Connect to the SQLite database
	db, err := sql.Open("sqlite", "./my.db")
	if err != nil {
		fmt.Println(err)
		return
	}

	defer db.Close()
	fmt.Println("Connected to the SQLite database successfully.")

	// Get the version of SQLite
	var sqliteVersion string
	err = db.QueryRow("select sqlite_version()").Scan(&sqliteVersion)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(sqliteVersion)

	// Read the migrations folder
	version, err := operations.GetVersion(db)
	if err != nil {
		panic(fmt.Sprintln("FATAL: impossible to get version -", err))
	}
	fmt.Println("Current DB version: ", version)

	// Apply all migrations files which havn't been executed yet

	// Stamp the migration

}
