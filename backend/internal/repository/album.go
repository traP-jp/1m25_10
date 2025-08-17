package repository

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type AlbumRepository interface {
	GetAlbums(ctx context.Context, filter AlbumFilter) ([]uuid.UUID, error)
	PostAlbum(ctx context.Context, params PostAlbumParams) (*Album, error)
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

// uuidの配列をdbに保存するための型(json形式で保存)
type dbUUIDs []uuid.UUID

// dbUUIDsをjsonに変換
func (u dbUUIDs) Value() (driver.Value, error) {
	if u == nil {
		return nil, nil
	}
	return json.Marshal(u)
}

// jsonからdbUUIDsを復元
func (u *dbUUIDs) Scan(value interface{}) error {
	if value == nil {
		*u = nil
		return nil
	}
	b, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}
	return json.Unmarshal(b, u)
}

type dbAlbum struct {
	Id          uuid.UUID `db:"id"`
	Title       string    `db:"title"`
	Description string    `db:"description"`
	Creator     uuid.UUID `db:"creator"`
	Images      dbUUIDs   `db:"images"`
	CreatedAt   time.Time `db:"created_at"`
	UpdatedAt   time.Time `db:"updated_at"`
}

// TODO: Implement the actual SQL logic for retrieving albums based on filter
func (r *sqlAlbumRepository) GetAlbums(ctx context.Context, filter AlbumFilter) ([]uuid.UUID, error) {
	return nil, nil
}

// PostAlbum creates a new album and returns its ID
func (r *sqlAlbumRepository) PostAlbum(ctx context.Context, params PostAlbumParams) (*Album, error) {
	query := `
		INSERT INTO albums (id, title, description, creator, images, created_at, updated_at)
		VALUES (:id, :title, :description, :creator, :images, :created_at, :updated_at)
	`
	now := time.Now()
	newAlbum := dbAlbum{
		Id:          uuid.New(),
		Title:       params.Title,
		Description: params.Description,
		Creator:     params.Creator,
		Images:      dbUUIDs(params.Images),
		CreatedAt:   now,
		UpdatedAt:   now,
	}
	_, err := r.db.NamedExecContext(ctx, query, newAlbum)
	if err != nil {
		return nil, err
	}

	return &Album{
		Id:          newAlbum.Id,
		Title:       newAlbum.Title,
		Description: newAlbum.Description,
		Creator:     newAlbum.Creator,
		Images:      []uuid.UUID(newAlbum.Images),
		CreatedAt:   newAlbum.CreatedAt,
		UpdatedAt:   newAlbum.UpdatedAt,
	}, nil
}

var ErrNotFound = errors.New("not found")

// TODO: Implement the actual SQL logic for retrieving an album by ID
func (r *sqlAlbumRepository) GetAlbum(ctx context.Context, albumID uuid.UUID) (*Album, error) {
	var dbAlbumModel dbAlbum
	query := `
		SELECT id, title, description, creator, images, created_at, updated_at
		FROM albums
		Where id = ?
		`
	query = r.db.Rebind(query)
	err := r.db.GetContext(ctx, &dbAlbumModel, query, albumID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("failed to get album (id=%s) : %w", albumID, err)

	}
	return &Album{
		Id:          dbAlbumModel.Id,
		Title:       dbAlbumModel.Title,
		Description: dbAlbumModel.Description,
		Creator:     dbAlbumModel.Creator,
		Images:      dbAlbumModel.Images,
		CreatedAt:   dbAlbumModel.CreatedAt,
		UpdatedAt:   dbAlbumModel.UpdatedAt,
	}, nil
}

// DeleteAlbum deletes an album by its ID
func (r *sqlAlbumRepository) DeleteAlbum(ctx context.Context, albumID uuid.UUID) error {
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

var ErrNoFieldsToUpdate = errors.New("no fields to update")

// UpdateAlbum updates an album with the given parameters
func (r *sqlAlbumRepository) UpdateAlbum(ctx context.Context, albumID uuid.UUID, params UpdateAlbumParams) error {
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
		sets = append(sets, "images = ?")
		args = append(args, dbUUIDs(*params.Images))
	}

	if len(sets) == 0 {
		return ErrNoFieldsToUpdate
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
