package repository

import (
	"github.com/jmoiron/sqlx"
)

type Repository interface {
	AlbumRepository
}

type sqlRepositoryImpl struct {
	db *sqlx.DB
}

func New(db *sqlx.DB) Repository {
	return &sqlRepositoryImpl{db: db}
}
