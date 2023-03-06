package repository

import (
	"context"
	"database/sql"

	"github.com/shota-tech/graphql/server/graph/model"
)

type (
	ITodoRepository interface {
		Store(context.Context, *model.Todo) error
		List(context.Context) ([]*model.Todo, error)
	}

	TodoRepository struct {
		db *sql.DB
	}
)

func NewTodoRepository(db *sql.DB) *TodoRepository {
	return &TodoRepository{db: db}
}

func (r *TodoRepository) Store(ctx context.Context, todo *model.Todo) error {
	query := "INSERT INTO todos (id, text, done, user_id) VALUES (?, ?, ?, ?) " +
		"ON DUPLICATE KEY UPDATE text = VALUES(text), done = VALUES(done), user_id = VALUES(user_id);"
	_, err := r.db.ExecContext(ctx, query, todo.ID, todo.Text, todo.Done, todo.UserID)
	return err
}

func (r *TodoRepository) List(ctx context.Context) ([]*model.Todo, error) {
	query := "SELECT id, text, done, user_id FROM todos;"
	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	todos := make([]*model.Todo, 0)
	for rows.Next() {
		todo := &model.Todo{}
		if err := rows.Scan(&todo.ID, &todo.Text, &todo.Done, &todo.UserID); err != nil {
			return nil, err
		}
		todos = append(todos, todo)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return todos, nil
}
