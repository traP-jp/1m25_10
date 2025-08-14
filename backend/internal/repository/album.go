package repository

import (
	"context"
	"database/sql"
	"time"

	"github.com/google/uuid"
)

type AlbumRepository interface {
	GetAlbums(ctx context.Context, filter AlbumFilter) ([]uuid.UUID, error)
	PostAlbum(ctx context.Context, params PostAlbumParams) (uuid.UUID, error)
	GetAlbum(ctx context.Context, albumID uuid.UUID) (*Album, error)
	DeleteAlbum(ctx context.Context, albumID uuid.UUID) error
	UpdateAlbum(ctx context.Context, albumID uuid.UUID, params UpdateAlbumParams) error
}

// TODO: Domainで適切に定義予定（issue #27）
type (
	Album struct {
		Id          uuid.UUID
		Title       string
		Description string
		Creator     uuid.UUID
		Images      []uuid.UUID
		CreatedAt   time.Time
		UpdatedAt   time.Time
	}

	AlbumFilter struct {
		CreatorID  *uuid.UUID
		IsFavorite *bool
		// other conditions can be added as needed
	}

	PostAlbumParams struct {
		Title       string
		Description string
		Creator     uuid.UUID
		Images      []uuid.UUID
	}

	UpdateAlbumParams struct {
		Title       *string
		Description *string
		Images      *[]uuid.UUID
	}
)

type sqlAlbumRepository struct {
	db *sql.DB
}

func NewAlbumRepository(db *sql.DB) AlbumRepository {
	return &sqlAlbumRepository{db: db}
}

//TODO: Implement the actual SQL logic for retrieving albums based on filter
func (r *sqlAlbumRepository) GetAlbums(ctx context.Context, filter AlbumFilter) ([]uuid.UUID, error) {
	return nil, nil
}

// TODO: Implement the actual SQL logic for inserting a new album
func (r *sqlAlbumRepository) PostAlbum(ctx context.Context, params PostAlbumParams) (uuid.UUID, error) {
	return uuid.Nil, nil
}

// TODO: Implement the actual SQL logic for retrieving an album by ID
func (r *sqlAlbumRepository) GetAlbum(ctx context.Context, albumID uuid.UUID) (*Album, error) {
	return nil, nil
}

// TODO: Implement the actual SQL logic for deleting an album by ID
func (r *sqlAlbumRepository) DeleteAlbum(ctx context.Context, albumID uuid.UUID) error {
	return nil
}

//TODO: Implement the actual SQL logic for updating an album
func (r *sqlAlbumRepository) UpdateAlbum(ctx context.Context, albumID uuid.UUID, params UpdateAlbumParams) error {
	return nil
}