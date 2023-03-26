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
	ITaskRepository interface {
		Store(context.Context, *model.Task) error
		Get(context.Context, string) (*model.Task, error)
		List(context.Context, []string) ([]*model.Task, error)
		ListByUserID(context.Context, string) ([]*model.Task, error)
		ListByUserIDs(context.Context, []string) ([]*model.Task, error)
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
	row := models.Task{
		ID:     task.ID,
		Text:   task.Text,
		Status: task.Status.String(),
		UserID: task.UserID,
	}
	if err := row.Upsert(ctx, r.db, boil.Infer(), boil.Infer()); err != nil {
		return fmt.Errorf("failed to upsert record: %w", err)
	}
	return nil
}

func (r *TaskRepository) Get(ctx context.Context, id string) (*model.Task, error) {
	row, err := models.Tasks(models.TaskWhere.ID.EQ(id)).One(ctx, r.db)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("record not found")
		}
		return nil, fmt.Errorf("failed to get record: %w", err)
	}
	return &model.Task{
		ID:     row.ID,
		Text:   row.Text,
		Status: model.Status(row.Status),
		UserID: row.UserID,
	}, nil
}

func (r *TaskRepository) List(ctx context.Context, ids []string) ([]*model.Task, error) {
	rows, err := models.Tasks(models.TaskWhere.ID.IN(ids)).All(ctx, r.db)
	if err != nil {
		return nil, fmt.Errorf("failed to get records: %w", err)
	}
	tasks := make([]*model.Task, len(rows))
	for i, row := range rows {
		tasks[i] = &model.Task{
			ID:     row.ID,
			Text:   row.Text,
			Status: model.Status(row.Status),
			UserID: row.UserID,
		}
	}
	return tasks, nil
}

func (r *TaskRepository) ListByUserID(ctx context.Context, userID string) ([]*model.Task, error) {
	rows, err := models.Tasks(models.TaskWhere.UserID.EQ(userID)).All(ctx, r.db)
	if err != nil {
		return nil, fmt.Errorf("failed to get records: %w", err)
	}
	tasks := make([]*model.Task, len(rows))
	for i, row := range rows {
		tasks[i] = &model.Task{
			ID:     row.ID,
			Text:   row.Text,
			Status: model.Status(row.Status),
			UserID: row.UserID,
		}
	}
	return tasks, nil
}

func (r *TaskRepository) ListByUserIDs(ctx context.Context, userIDs []string) ([]*model.Task, error) {
	rows, err := models.Tasks(models.TaskWhere.UserID.IN(userIDs)).All(ctx, r.db)
	if err != nil {
		return nil, fmt.Errorf("failed to get records: %w", err)
	}
	tasks := make([]*model.Task, len(rows))
	for i, row := range rows {
		tasks[i] = &model.Task{
			ID:     row.ID,
			Text:   row.Text,
			Status: model.Status(row.Status),
			UserID: row.UserID,
		}
	}
	return tasks, nil
}
