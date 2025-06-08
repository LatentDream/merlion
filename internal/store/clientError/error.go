package clientError

import "errors"

// ErrNoteNotFound is returned when a note with the given ID is not found
var ErrNoteNotFound = errors.New("note not found")

