package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"github.com/shota-tech/graphql/server/graph/model"
)

type (
	IUserRepository interface {
		Store(context.Context, *model.User) error
		Get(context.Context, string) (*model.User, error)
		List(context.Context, []string) ([]*model.User, error)
	}

	UserRepository struct {
		db *sql.DB
	}
)

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) Store(ctx context.Context, user *model.User) error {
	if user == nil {
		return errors.New("user is required")
	}
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
	err := r.db.QueryRowContext(ctx, query, id).Scan(&user.ID, &user.Name)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("record not found")
		}
		return nil, fmt.Errorf("failed to get record: %w", err)
	}
	return user, nil
}

func (r *UserRepository) List(ctx context.Context, ids []string) ([]*model.User, error) {
	args := make([]any, len(ids))
	for i, id := range ids {
		args[i] = id
	}
	query := "SELECT id, name FROM users " +
		"WHERE id IN (?" + strings.Repeat(",?", len(ids)-1) + ");"
	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to get records: %w", err)
	}
	defer rows.Close()

	users := make([]*model.User, 0)
	for rows.Next() {
		user := &model.User{}
		if err := rows.Scan(&user.ID, &user.Name); err != nil {
			return nil, fmt.Errorf("failed to scan record: %w", err)
		}
		users = append(users, user)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("failed to scan records: %w", err)
	}
	return users, nil
}
