package domain

import "github.com/google/uuid"

// Image represents an image entity in the domain
type Image struct {
	ID      uuid.UUID `json:"id"`
	Creator uuid.UUID `json:"creator"`
	Post    Post      `json:"post"`
}

// Post represents a post entity in the domain
type Post struct {
	ID      uuid.UUID `json:"id"`
	Content string    `json:"content"`
}
