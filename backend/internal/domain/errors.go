package domain

import "errors"

// Domain errors
var (
	ErrNotFound         = errors.New("not found")
	ErrNoFieldsToUpdate = errors.New("no fields to update")
)
