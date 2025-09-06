package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/google/uuid"
)

type ImageRepository interface {
	PostImage(ctx context.Context, ImageID uuid.UUID) (*uuid.UUID, error)
	GetImage(ctx context.Context, imageID uuid.UUID) (*uuid.UUID, error)
}

// PostImage stores a new image in the database.
func (r *sqlRepositoryImpl) PostImage(ctx context.Context, ImageID uuid.UUID) (*uuid.UUID, error) {
	query := `INSERT INTO images (id) VALUES (?)`
	_, err := r.db.ExecContext(ctx, query, ImageID)
	if err != nil {
		return nil, fmt.Errorf("failed to insert image: %w", err)
	}

	return &ImageID, nil
}

// GetImage retrieves an image by its ID.
func (r *sqlRepositoryImpl) GetImage(ctx context.Context, imageID uuid.UUID) (*uuid.UUID, error) {
	var id uuid.UUID
	query := `SELECT id FROM images WHERE id = ?`
	err := r.db.QueryRowContext(ctx, query, imageID).Scan(&id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("failed to get image: %w", err)
	}

	return &id, nil
}
