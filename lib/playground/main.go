package main

import (
	"cmp"
	"database/sql"
	"embed"
	"fmt"
	"io/fs"
	"slices"
	"strconv"
	"strings"

	_ "github.com/glebarez/go-sqlite"
	"github.com/latentdream/merlion/lib/playground/operations"
)

//go:embed migrations/*.sql
var migrationsContent embed.FS

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

	// List all files in the embedded migrations directory
	entries, err := migrationsContent.ReadDir("migrations")
	if err != nil {
		fmt.Println("Error listing embedded migrations:", err)
		return
	}

	fmt.Println("Files in the embedded migrations directory:")

	sortedFunc := func(entry1 fs.DirEntry, entry2 fs.DirEntry) int {
		parse := func(entry fs.DirEntry) int {
			parts := strings.Split(entry.Name(), "__")
			AssertEq(len(parts), 2, "File format should be int__description.sql, like: 024__adding_something.sql but found:", entry.Name())
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
		AssertEq(len(parts), 2, "File format should be filename.ext, found:", entry.Name())
		if parts[1] != "sql" {
			return true
		}
		return false
	}

	entries = slices.DeleteFunc(entries, filterFunc)
	slices.SortFunc(entries, sortedFunc)

	for _, entry := range entries {
		fmt.Println(entry.Name())
	}
}

func AssertEq[T comparable](expected T, got T, message ...string) {
	if expected == got {
		return
	}

	panicMsg := fmt.Sprintf("Assertion failed: expected %v, got %v", expected, got)
	if len(message) > 0 {
		panicMsg = strings.Join(message, " ")
	}
	panic(panicMsg)
}
