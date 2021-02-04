package post

import "errors"

var (
	// ErrNotFound is thrown when specified post was not found in database.
	ErrNotFound = errors.New("specified post was not found")
)
