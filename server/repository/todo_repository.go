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
	query := "INSERT INTO todos (id, text, done, task_id) VALUES (?, ?, ?, ?) " +
		"ON DUPLICATE KEY UPDATE text = VALUES(text), done = VALUES(done), task_id = VALUES(task_id);"
	_, err := r.db.ExecContext(ctx, query, todo.ID, todo.Text, todo.Done, todo.TaskID)
	if err != nil {
		return fmt.Errorf("failed to upsert record: %w", err)
	}
	return nil
}

func (r *TodoRepository) Get(ctx context.Context, id string) (*model.Todo, error) {
	todo := &model.Todo{}
	query := "SELECT id, text, done, task_id FROM todos WHERE id = ?;"
	err := r.db.QueryRowContext(ctx, query, id).Scan(&todo.ID, &todo.Text, &todo.Done, &todo.TaskID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("record not found")
		}
		return nil, fmt.Errorf("failed to get record: %w", err)
	}
	return todo, nil
}

func (r *TodoRepository) ListByTaskIDs(ctx context.Context, taskIDs []string) ([]*model.Todo, error) {
	args := make([]any, len(taskIDs))
	for i, taskID := range taskIDs {
		args[i] = taskID
	}
	query := "SELECT id, text, done, task_id FROM todos WHERE task_id IN (?" + strings.Repeat(",?", len(taskIDs)-1) + ");"
	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to get records: %w", err)
	}
	defer rows.Close()

	todos := make([]*model.Todo, 0)
	for rows.Next() {
		todo := &model.Todo{}
		if err := rows.Scan(&todo.ID, &todo.Text, &todo.Done, &todo.TaskID); err != nil {
			return nil, fmt.Errorf("failed to scan record: %w", err)
		}
		todos = append(todos, todo)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("failed to scan records: %w", err)
	}
	return todos, nil
}
