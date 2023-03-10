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
		Get(context.Context, string) (*model.Task, error)
		ListByUserID(context.Context, string) ([]*model.Task, error)
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

func (r *TaskRepository) Get(ctx context.Context, id string) (*model.Task, error) {
	task := &model.Task{}
	query := "SELECT id, text, status, user_id FROM tasks WHERE id = ?;"
	err := r.db.QueryRowContext(ctx, query, id).Scan(&task.ID, &task.Text, &task.Status, &task.UserID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("record not found")
		}
		return nil, fmt.Errorf("failed to get record: %w", err)
	}
	return task, nil
}

func (r *TaskRepository) ListByUserID(ctx context.Context, userID string) ([]*model.Task, error) {
	query := "SELECT id, text, status, user_id FROM tasks WHERE user_id = ?;"
	rows, err := r.db.QueryContext(ctx, query, userID)
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
