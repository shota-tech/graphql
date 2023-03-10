package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/shota-tech/graphql/server/graph/model"
)

type (
	ITaskRepository interface {
		Store(context.Context, *model.Task) error
		List(context.Context) ([]*model.Task, error)
	}

	TaskRepository struct {
		db *sql.DB
	}
)

func NewTaskRepository(db *sql.DB) *TaskRepository {
	return &TaskRepository{db: db}
}

func (r *TaskRepository) Store(ctx context.Context, task *model.Task) error {
	if task == nil {
		return errors.New("task is required")
	}
	query := "INSERT INTO tasks (id, text, status, user_id) VALUES (?, ?, ?, ?) " +
		"ON DUPLICATE KEY UPDATE text = VALUES(text), status = VALUES(status), user_id = VALUES(user_id);"
	_, err := r.db.ExecContext(ctx, query, task.ID, task.Text, task.Status.String(), task.UserID)
	if err != nil {
		return fmt.Errorf("failed to upsert record: %w", err)
	}
	return nil
}

func (r *TaskRepository) List(ctx context.Context) ([]*model.Task, error) {
	query := "SELECT id, text, status, user_id FROM tasks;"
	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to get records: %w", err)
	}
	defer rows.Close()

	tasks := make([]*model.Task, 0)
	for rows.Next() {
		task := &model.Task{}
		if err := rows.Scan(&task.ID, &task.Text, &task.Status, &task.UserID); err != nil {
			return nil, fmt.Errorf("failed to scan record: %w", err)
		}
		tasks = append(tasks, task)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("failed to scan records: %w", err)
	}
	return tasks, nil
}
