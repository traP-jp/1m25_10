package repository

import (
	"context"
	"fmt"

	"github.com/google/uuid"
)

type UserRepository interface {
	GetUsers(ctx context.Context) ([]*User, error)
	CreateUser(ctx context.Context, params CreateUserParams) (uuid.UUID, error)
	GetUser(ctx context.Context, userID uuid.UUID) (*User, error)
}

type (
	// users table
	User struct {
		ID    uuid.UUID `db:"id"`
		Name  string    `db:"name"`
		Email string    `db:"email"`
	}

	CreateUserParams struct {
		Name  string
		Email string
	}
)

func (r *sqlRepositoryImpl) GetUsers(ctx context.Context) ([]*User, error) {
	users := []*User{}
	if err := r.db.SelectContext(ctx, &users, "SELECT * FROM users"); err != nil {
		return nil, fmt.Errorf("select users: %w", err)
	}

	return users, nil
}

func (r *sqlRepositoryImpl) CreateUser(ctx context.Context, params CreateUserParams) (uuid.UUID, error) {
	userID := uuid.New()
	if _, err := r.db.ExecContext(ctx, "INSERT INTO users (id, name, email) VALUES (?, ?, ?)", userID, params.Name, params.Email); err != nil {
		return uuid.Nil, fmt.Errorf("insert user: %w", err)
	}

	return userID, nil
}

func (r *sqlRepositoryImpl) GetUser(ctx context.Context, userID uuid.UUID) (*User, error) {
	user := &User{}
	if err := r.db.GetContext(ctx, user, "SELECT * FROM users WHERE id = ?", userID); err != nil {
		return nil, fmt.Errorf("select user: %w", err)
	}

	return user, nil
}
