package export

import (
	"fmt"
	"log"

	"merlion/internal/model"
	"merlion/internal/store/files"
	"merlion/internal/store/sqlite"
	"merlion/internal/store/sqlite/database"
)

func ExportCmd(args ...string) int {
	if len(args) == 0 || args[0] == "--help" || args[0] == "-h" {
		fmt.Println("Usage: merlion export <path>")
		return 1
	}

	fmt.Println("Exporting notes to Obsidian vault")
	fmt.Println("Path:", args[0])

	localDB, err := database.InitDB()
	if err != nil {
		log.Fatalf("Failed to init DB: %v", err)
	}
	defer localDB.Close()
	sqliteClient := sqlite.NewClient(localDB)
	fileClient, err := files.NewClient(args[0])
	if err != nil {
		log.Fatalf("Failed to init file client: %v", err)
	}

	notes, err := sqliteClient.ListNotes()
	if err != nil {
		log.Fatalf("Failed to list notes: %v", err)
	}

	for _, note := range notes {
		note, err := sqliteClient.GetNote(note.NoteID)
		if err != nil {
			log.Fatalf("Failed to get note: %v", err)
		}

		req := model.CreateNoteRequest{
			Title:       note.Title,
			Content:     note.Content,
			Tags:        note.Tags,
			IsFavorite:  &note.IsFavorite,
			IsWorkLog:   &note.IsWorkLog,
			IsPublic:    &note.IsPublic,
			CreatedAt:   &note.CreatedAt,
			UpdatedAt:   &note.UpdatedAt,
			WorkspaceID: note.WorkspaceID,
		}
		fileClient.CreateNote(req)
	}

	return 0
}
