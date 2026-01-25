// Package export implements the export command allowing user to export notes from one store to another
package export

import (
	"fmt"
	"os"

	"merlion/cmd/merlion/parser"
	"merlion/internal/model"
	"merlion/internal/store"
	"merlion/internal/store/cloud"
	"merlion/internal/store/files"
	"merlion/internal/store/sqlite"
	"merlion/internal/utils"

	"github.com/charmbracelet/log"
)

func initSqliteDB() (*sqlite.Client, error) {
	sqliteClient := sqlite.NewClient()
	return sqliteClient, nil
}

func initFileClient(path string) (*files.Client, error) {
	fileClient, err := files.NewClient(path, "")
	if err != nil {
		return nil, err
	}
	return fileClient, nil
}

func initCloudClient() (*cloud.Client, error) {
	credentialsManager, err := cloud.NewCredentialsManager()
	if err != nil {
		log.Fatalf("Failed to initialize credentials manager: %v", err)
	}
	creds, err := credentialsManager.LoadCredentials()
	if err != nil {
		log.Fatalf("Failed to load credentials: %v", err)
	}
	if creds == nil {
		log.Fatalf("You need to login to use Cloud")
	}

	cloudClient, err := cloud.NewClient(creds)
	if err != nil {
		log.Fatalf("Failed to init cloud client: %v", err)
	}

	return cloudClient, nil
}

func printHelp(invalidArgs bool) {
	if invalidArgs {
		fmt.Println("Invalid arguments")
	}

	fmt.Println("Usage: merlion export <from-provider> <to-provider>")
	fmt.Println("Provider: sqlite, file <obsidian-vault-path>, cloud")
	fmt.Println("  - sqlite: export to a local SQLite database")
	fmt.Println("  - file <obsidian-vault-path>: export to a local Obsidian vault")
	fmt.Println("  - cloud: export to a cloud storage provider")
	fmt.Println("Examples:")
	fmt.Println("  merlion export sqlite files ~/notes")
	fmt.Println("  merlion export sqlite cloud")

	if invalidArgs {
		os.Exit(1)
	}
	os.Exit(0)
}

func parseStore(args []string) (store.Store, []string) {
	arg, args := parser.GetArg(args, printHelp)
	switch arg {
	case "sqlite":
		client, err := initSqliteDB()
		if err != nil {
			log.Fatalf("Failed to init SQLite client: %v", err)
		}
		return client, args
	case "file":
		arg, args = parser.GetArg(args, printHelp)
		client, err := initFileClient(arg)
		if err != nil {
			log.Fatalf("Failed to init file client: %v", err)
		}
		return client, args
	case "cloud":
		client, err := initCloudClient()
		if err != nil {
			log.Fatalf("Failed to init cloud client: %v", err)
		}
		return client, args
	default:
		printHelp(true)
	}
	return nil, nil
}

func Cmd(args ...string) int {
	if len(args) == 0 || utils.Contains(args, "--help") || utils.Contains(args, "-h") {
		printHelp(false)
	}

	var fromStore, toStore store.Store = nil, nil
	fromStore, args = parseStore(args)
	toStore, _ = parseStore(args)

	fmt.Printf("Exporting %s -> %s\n", fromStore.Name(), toStore.Name())

	notes, err := fromStore.ListNotes()
	if err != nil {
		log.Fatalf("Failed to list notes: %v", err)
	}

	nbErrors := 0
	for _, note := range notes {
		note, err := fromStore.GetNote(note.NoteID)
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
		_, err = toStore.CreateNote(req)
		if err != nil {
			fmt.Printf("Failed to create note '%s': %v\n", note.Title, err)
			nbErrors++
		}
	}

	fmt.Printf("Exported %d notes\n", len(notes)-nbErrors)

	return 0
}
