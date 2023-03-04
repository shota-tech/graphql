package repository

import (
	"context"
	"fmt"

	"github.com/shota-tech/graphql/server/graph/model"
)

type (
	IUserRepository interface {
		Store(context.Context, *model.User) error
		Get(context.Context, string) (*model.User, error)
	}

	UserRepository struct {
		data map[string]*model.User
	}
)

func NewUserRepository() *UserRepository {
	return &UserRepository{
		data: make(map[string]*model.User),
	}
}

func (r *UserRepository) Store(_ context.Context, user *model.User) error {
	r.data[user.ID] = user
	return nil
}

func (r *UserRepository) Get(_ context.Context, id string) (*model.User, error) {
	user, ok := r.data[id]
	if !ok {
		return nil, fmt.Errorf("user %s not found", id)
	}
	return user, nil
}
