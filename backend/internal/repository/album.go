package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/traP-jp/1m25_10/backend/internal/domain"
)

type AlbumRepository interface {
	GetAlbums(ctx context.Context, filter domain.AlbumFilter) ([]domain.AlbumItem, error)
	PostAlbum(ctx context.Context, params domain.PostAlbumParams) (*domain.Album, error)
	GetAlbum(ctx context.Context, albumID uuid.UUID) (*domain.Album, error)
	DeleteAlbum(ctx context.Context, albumID uuid.UUID) error
	UpdateAlbum(ctx context.Context, albumID uuid.UUID, params domain.UpdateAlbumParams) error
}

// AlbumImage represents the relationship between albums and images (repository-specific)
type AlbumImage struct {
	Id      uuid.UUID `db:"id"`
	AlbumID uuid.UUID `db:"album_id"`
	ImageID uuid.UUID `db:"image_id"`
}

type dbAlbum struct {
	Id          uuid.UUID `db:"id"`
	Title       string    `db:"title"`
	Description string    `db:"description"`
	Creator     uuid.UUID `db:"creator"`
	CreatedAt   time.Time `db:"created_at"`
	UpdatedAt   time.Time `db:"updated_at"`
}

type dbAlbumItem struct {
	Id      uuid.UUID `db:"id"`
	Title   string    `db:"title"`
	Creator uuid.UUID `db:"creator"`
}

// GetAlbums retrieves albums based on the provided filter.
func (r *sqlRepositoryImpl) GetAlbums(ctx context.Context, filter domain.AlbumFilter) ([]domain.AlbumItem, error) {
	query := `SELECT id, title, creator FROM albums WHERE 1=1`
	args := []interface{}{}

	if filter.CreatorID != nil {
		query += " AND creator = ?"
		args = append(args, *filter.CreatorID)
	}
	if filter.AfterDate != nil {
		query += " AND created_at >= ?"
		args = append(args, *filter.AfterDate)
	}
	if filter.BeforeDate != nil {
		query += " AND created_at <= ?"
		args = append(args, *filter.BeforeDate)
	}

	// created_atに基づき並べ替え
	query += " ORDER BY created_at DESC"

	const maxLimit = 100
	lim := 20 // Default limit
	if filter.Limit != nil {
		if *filter.Limit > 0 && *filter.Limit < maxLimit {
			lim = *filter.Limit
		} else {
			lim = maxLimit
		}
	}
	query += " LIMIT ?"
	args = append(args, lim)

	if filter.Offset != nil {
		query += " OFFSET ?"
		args = append(args, *filter.Offset)
	}

	query = r.db.Rebind(query)

	var dbItems []dbAlbumItem
	if err := r.db.SelectContext(ctx, &dbItems, query, args...); err != nil {
		return nil, fmt.Errorf("failed to select album items: %w", err)
	}

	items := make([]domain.AlbumItem, len(dbItems))
	for i, d := range dbItems {
		items[i] = domain.AlbumItem(d)
	}

	return items, nil
}

// PostAlbum creates a new album and returns its ID
func (r *sqlRepositoryImpl) PostAlbum(ctx context.Context, params domain.PostAlbumParams) (*domain.Album, error) {
	query := `
		INSERT INTO albums (id, title, description, creator, created_at, updated_at)
		VALUES (:id, :title, :description, :creator, :created_at, :updated_at)
	`
	now := time.Now()
	newAlbum := dbAlbum{
		Id:          uuid.New(),
		Title:       params.Title,
		Description: params.Description,
		Creator:     params.Creator,
		CreatedAt:   now,
		UpdatedAt:   now,
	}
	_, err := r.db.NamedExecContext(ctx, query, newAlbum)
	if err != nil {
		return nil, err
	}

	query = `
		INSERT INTO album_images (id, album_id, image_id)
		VALUES (:id, :album_id, :image_id)
	`
	for _, imgID := range params.Images {
		newAlbumImage := AlbumImage{
			Id:      uuid.New(),
			AlbumID: newAlbum.Id,
			ImageID: imgID,
		}
		_, err := r.db.NamedExecContext(ctx, query, newAlbumImage)
		if err != nil {
			return nil, err
		}
	}

	return &domain.Album{
		Id:          newAlbum.Id,
		Title:       newAlbum.Title,
		Description: newAlbum.Description,
		Creator:     newAlbum.Creator,
		Images:      params.Images,
		CreatedAt:   newAlbum.CreatedAt,
		UpdatedAt:   newAlbum.UpdatedAt,
	}, nil
}

var ErrNotFound = domain.ErrNotFound

// GetAlbum retrieves an album by its ID
func (r *sqlRepositoryImpl) GetAlbum(ctx context.Context, albumID uuid.UUID) (*domain.Album, error) {
	var dbAlbumModel dbAlbum
	query := `
		SELECT id, title, description, creator, created_at, updated_at
		FROM albums
		WHERE id = ?
		`
	query = r.db.Rebind(query)
	err := r.db.GetContext(ctx, &dbAlbumModel, query, albumID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("failed to get album (id=%s) : %w", albumID, err)

	}

	var imageIDs []struct {
		ImageID uuid.UUID `db:"image_id"`
	}
	query = `
		SELECT image_id
		FROM album_images
		WHERE album_id = ?
	`
	query = r.db.Rebind(query)
	err = r.db.SelectContext(ctx, &imageIDs, query, albumID)
	if err != nil {
		return nil, fmt.Errorf("failed to get album images (album_id=%s) : %w", albumID, err)
	}
	images := make([]uuid.UUID, len(imageIDs))
	for i, id := range imageIDs {
		images[i] = id.ImageID
	}

	return &domain.Album{
		Id:          dbAlbumModel.Id,
		Title:       dbAlbumModel.Title,
		Description: dbAlbumModel.Description,
		Creator:     dbAlbumModel.Creator,
		Images:      images,
		CreatedAt:   dbAlbumModel.CreatedAt,
		UpdatedAt:   dbAlbumModel.UpdatedAt,
	}, nil
}

// DeleteAlbum deletes an album by its ID
func (r *sqlRepositoryImpl) DeleteAlbum(ctx context.Context, albumID uuid.UUID) error {
	if albumID == uuid.Nil {
		return fmt.Errorf("invalid album id")
	}

	query := `
		DELETE FROM albums
		WHERE id = ?
	`
	query = r.db.Rebind(query)
	result, err := r.db.ExecContext(ctx, query, albumID)
	if err != nil {
		return fmt.Errorf("failed to delete album (id=%s) : %w", albumID, err)
	}

	ra, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected (id=%s) : %w", albumID, err)
	}
	if ra == 0 {
		return ErrNotFound
	}
	return nil
}

// UpdateAlbum updates an album with the given parameters
func (r *sqlRepositoryImpl) UpdateAlbum(ctx context.Context, albumID uuid.UUID, params domain.UpdateAlbumParams) error {
	if albumID == uuid.Nil {
		return fmt.Errorf("invalid album id")
	}

	sets := []string{}
	args := []interface{}{}

	if params.Title != nil {
		sets = append(sets, "title = ?")
		args = append(args, *params.Title)
	}
	if params.Description != nil {
		sets = append(sets, "description = ?")
		args = append(args, *params.Description)
	}
	if params.Images != nil {
		// Imageは差分更新できるようにしてもいいかもしれない
		// 以下のコードは全て置き換える実装
		// 既存の関係を削除
		delQuery := `DELETE FROM album_images WHERE album_id = ?`
		delQuery = r.db.Rebind(delQuery)
		_, err := r.db.ExecContext(ctx, delQuery, albumID)
		if err != nil {
			return fmt.Errorf("failed to delete existing album images (album_id=%s): %w", albumID, err)
		}

		// 新しい関係を挿入
		insQuery := `
			INSERT INTO album_images (id, album_id, image_id)
			VALUES (:id, :album_id, :image_id)
		`
		for _, imgID := range *params.Images {
			newAlbumImage := AlbumImage{
				Id:      uuid.New(),
				AlbumID: albumID,
				ImageID: imgID,
			}
			_, err := r.db.NamedExecContext(ctx, insQuery, newAlbumImage)
			if err != nil {
				return fmt.Errorf("failed to insert new album image (album_id=%s, image_id=%s): %w", albumID, imgID, err)
			}
		}
	}

	if len(sets) == 0 {
		return domain.ErrNoFieldsToUpdate
	}

	sets = append(sets, "updated_at = ?")
	args = append(args, time.Now())

	args = append(args, albumID)

	query := "UPDATE albums SET " + strings.Join(sets, ", ") + " WHERE id = ?"
	query = r.db.Rebind(query)

	result, err := r.db.ExecContext(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("failed to update album (id=%s): %w", albumID, err)
	}
	ra, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("rows affected error (id=%s): %w", albumID, err)
	}
	if ra == 0 {
		return ErrNotFound
	}
	return nil
}
