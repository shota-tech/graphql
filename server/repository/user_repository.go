package repository

import (
	"context"
	"database/sql"
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
	stmt, err := r.db.PrepareContext(ctx, query)
	if err != nil {
		return err
	}
	defer stmt.Close()
	_, err = stmt.ExecContext(ctx, user.ID, user.Name)
	return err
}

func (r *UserRepository) Get(ctx context.Context, id string) (*model.User, error) {
	user := &model.User{}
	stmt, err := r.db.PrepareContext(ctx, "SELECT id, name FROM users WHERE id = ?;")
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	if err := stmt.QueryRowContext(ctx, id).Scan(&user.ID, &user.Name); err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("user %s not found", id)
		}
		return nil, err
	}
	return user, nil
}
