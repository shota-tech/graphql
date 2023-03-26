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
	ITodoRepository interface {
		Store(context.Context, *model.Todo) error
		Get(context.Context, string) (*model.Todo, error)
		ListByTaskIDs(context.Context, []string) ([]*model.Todo, error)
	}

	TodoRepository struct {
		db *sql.DB
	}
)

func NewTodoRepository(db *sql.DB) *TodoRepository {
	return &TodoRepository{db: db}
}

func (r *TodoRepository) Store(ctx context.Context, todo *model.Todo) error {
	if todo == nil {
		return errors.New("todo is required")
	}
	row := models.Todo{
		ID:     todo.ID,
		Text:   todo.Text,
		Done:   todo.Done,
		TaskID: todo.TaskID,
	}
	if err := row.Upsert(ctx, r.db, boil.Infer(), boil.Infer()); err != nil {
		return fmt.Errorf("failed to upsert record: %w", err)
	}
	return nil
}

func (r *TodoRepository) Get(ctx context.Context, id string) (*model.Todo, error) {
	row, err := models.Todos(models.TodoWhere.ID.EQ(id)).One(ctx, r.db)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("record not found")
		}
		return nil, fmt.Errorf("failed to get record: %w", err)
	}
	return &model.Todo{
		ID:     row.ID,
		Text:   row.Text,
		Done:   row.Done,
		TaskID: row.TaskID,
	}, nil
}

func (r *TodoRepository) ListByTaskIDs(ctx context.Context, taskIDs []string) ([]*model.Todo, error) {
	rows, err := models.Todos(models.TodoWhere.TaskID.IN(taskIDs)).All(ctx, r.db)
	if err != nil {
		return nil, fmt.Errorf("failed to get records: %w", err)
	}
	todos := make([]*model.Todo, len(rows))
	for i, row := range rows {
		todos[i] = &model.Todo{
			ID:     row.ID,
			Text:   row.Text,
			Done:   row.Done,
			TaskID: row.TaskID,
		}
	}
	return todos, nil
}
