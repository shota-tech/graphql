package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/shota-tech/graphql/server/graph/model"
	"github.com/shota-tech/graphql/server/repository/models"
	"github.com/volatiletech/sqlboiler/v4/boil"
)

type (
	IUserRepository interface {
		Store(context.Context, *model.User) error
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
	row := models.User{
		ID:   user.ID,
		Name: user.Name,
	}
	if err := row.Upsert(ctx, r.db, boil.Infer(), boil.Infer()); err != nil {
		return fmt.Errorf("failed to upsert record: %w", err)
	}
	return nil
}

func (r *UserRepository) List(ctx context.Context, ids []string) ([]*model.User, error) {
	rows, err := models.Users(models.UserWhere.ID.IN(ids)).All(ctx, r.db)
	if err != nil {
		return nil, fmt.Errorf("failed to get records: %w", err)
	}
	users := make([]*model.User, len(rows))
	for i, row := range rows {
		users[i] = &model.User{
			ID:   row.ID,
			Name: row.Name,
		}
	}
	return users, nil
}
