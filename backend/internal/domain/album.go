package domain

import (
	"time"

	"github.com/google/uuid"
)

// Album represents an album entity in the domain
type Album struct {
	Id          uuid.UUID   `json:"id"`
	Title       string      `json:"title"`
	Description string      `json:"description"`
	Creator     string      `json:"creator"`
	Images      []uuid.UUID `json:"images"`
	CreatedAt   time.Time   `json:"created_at"`
	UpdatedAt   time.Time   `json:"updated_at"`
}

// AlbumItem represents a simplified album item for list views
type AlbumItem struct {
	Id      uuid.UUID `json:"id"`
	Title   string    `json:"title"`
	Creator string    `json:"creator"`
}

// AlbumFilter represents filtering criteria for albums
type AlbumFilter struct {
	CreatorID  *string
	BeforeDate *time.Time // Filter by created_at
	AfterDate  *time.Time // Filter by created_at
	Limit      *int
	Offset     *int
	//あとはIsFavorite(*bool)とか？
}

// PostAlbumParams represents parameters for creating a new album
type PostAlbumParams struct {
	Title       string      `json:"title"`
	Description string      `json:"description"`
	Creator     string      `json:"creator"`
	Images      []uuid.UUID `json:"images"`
}

// UpdateAlbumParams represents parameters for updating an album
type UpdateAlbumParams struct {
	Title       *string      `json:"title,omitempty"`
	Description *string      `json:"description,omitempty"`
	Images      *[]uuid.UUID `json:"images,omitempty"`
}
