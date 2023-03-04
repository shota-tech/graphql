package repository

import (
	"context"

	"github.com/shota-tech/graphql/server/graph/model"
)

type (
	ITodoRepository interface {
		Create(context.Context, *model.Todo) error
		List(context.Context) ([]*model.Todo, error)
	}

	TodoRepository struct {
		data []*model.Todo
	}
)

func NewTodoRepository() *TodoRepository {
	return &TodoRepository{
		data: make([]*model.Todo, 0),
	}
}

func (r *TodoRepository) Create(_ context.Context, todo *model.Todo) error {
	r.data = append(r.data, todo)
	return nil
}

func (r *TodoRepository) List(_ context.Context) ([]*model.Todo, error) {
	return r.data, nil
}
