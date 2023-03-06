package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/shota-tech/graphql/server/graph/model"
)

type (
	IUserRepository interface {
		Store(context.Context, *model.User) error
		Get(context.Context, string) (*model.User, error)
	}

	UserRepository struct {
		db *sql.DB
	}
)

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) Store(ctx context.Context, user *model.User) error {
	query := "INSERT INTO users (id, name) VALUES (?, ?) ON DUPLICATE KEY UPDATE name = VALUES(name);"
	_, err := r.db.ExecContext(ctx, query, user.ID, user.Name)
	if err != nil {
		return fmt.Errorf("failed to upsert record: %w", err)
	}
	return nil
}

func (r *UserRepository) Get(ctx context.Context, id string) (*model.User, error) {
	user := &model.User{}
	query := "SELECT id, name FROM users WHERE id = ?;"
	if err := r.db.QueryRowContext(ctx, query, id).Scan(&user.ID, &user.Name); err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("record not found")
		}
		return nil, fmt.Errorf("failed to get records: %w", err)
	}
	return user, nil
}
