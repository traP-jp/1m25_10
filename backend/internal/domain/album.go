package domain

import (
	"time"

	"github.com/google/uuid"
)

// Album represents an album entity in the domain
type Album struct {
	Id          uuid.UUID
	Title       string
	Description string
	Creator     uuid.UUID
	Images      []uuid.UUID
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

// AlbumItem represents a simplified album item for listing
type AlbumItem struct {
	Id      uuid.UUID
	Title   string
	Creator uuid.UUID
}

// AlbumFilter represents filter criteria for album queries
type AlbumFilter struct {
	CreatorID  *uuid.UUID
	BeforeDate *time.Time // Filter by created_at
	AfterDate  *time.Time // Filter by created_at
	Limit      *int
	Offset     *int
	//あとはIsFavorite(*bool)とか？
}

// PostAlbumParams represents parameters for creating a new album
type PostAlbumParams struct {
	Title       string
	Description string
	Creator     uuid.UUID
	Images      []uuid.UUID
}

// UpdateAlbumParams represents parameters for updating an album
type UpdateAlbumParams struct {
	Title       *string
	Description *string
	Images      *[]uuid.UUID
}
