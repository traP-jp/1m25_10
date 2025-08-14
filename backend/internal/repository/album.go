package repository

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
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

	AlbumItem struct {
		Id      uuid.UUID
		Title   string
		Creator uuid.UUID
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
	db *sqlx.DB
}

func NewAlbumRepository(db *sqlx.DB) AlbumRepository {
	return &sqlAlbumRepository{db: db}
}

type dbAlbum struct {
	Id          uuid.UUID   `db:"id"`
	Title       string      `db:"title"`
	Description string      `db:"description"`
	Creator     uuid.UUID   `db:"creator"`
	Images      []uuid.UUID `db:"images"`
	CreatedAt   time.Time   `db:"created_at"`
	UpdatedAt   time.Time   `db:"updated_at"`
}

// TODO: Implement the actual SQL logic for retrieving albums based on filter
func (r *sqlAlbumRepository) GetAlbums(ctx context.Context, filter AlbumFilter) ([]uuid.UUID, error) {
	return nil, nil
}

// PostAlbum creates a new album and returns its ID
func (r *sqlAlbumRepository) PostAlbum(ctx context.Context, params PostAlbumParams) (uuid.UUID, error) {
	query := `
		INSERT INTO albums (id, title, description, creator, created_at, updated_at)
		VALUES (:id, :title, :description, :creator, :created_at, :updated_at)
	`
	now := time.Now()
	albumToInsert := dbAlbum{
		Id:          uuid.New(),
		Title:       params.Title,
		Description: params.Description,
		Creator:     params.Creator,
		CreatedAt:   now,
		UpdatedAt:   now,
	}
	_, err := r.db.NamedExecContext(ctx, query, albumToInsert)
	if err != nil {
		return uuid.Nil, err
	}

	return albumToInsert.Id, nil
}

// TODO: Implement the actual SQL logic for retrieving an album by ID
func (r *sqlAlbumRepository) GetAlbum(ctx context.Context, albumID uuid.UUID) (*Album, error) {
	return nil, nil
}

// TODO: Implement the actual SQL logic for deleting an album by ID
func (r *sqlAlbumRepository) DeleteAlbum(ctx context.Context, albumID uuid.UUID) error {
	return nil
}

// TODO: Implement the actual SQL logic for updating an album
func (r *sqlAlbumRepository) UpdateAlbum(ctx context.Context, albumID uuid.UUID, params UpdateAlbumParams) error {
	return nil
}
