package comment

import "errors"

var (
	// ErrNotFound is thrown when specified comment was not found in database.
	ErrNotFound = errors.New("specified comment was not found")
)
