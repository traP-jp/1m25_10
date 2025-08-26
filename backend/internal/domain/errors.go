package domain

import "errors"

// Common domain errors
var (
	ErrNotFound         = errors.New("not found")
	ErrNoFieldsToUpdate = errors.New("no fields to update")
)
