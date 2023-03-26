package loader

import (
	"context"
	"fmt"
	"log"

	"github.com/graph-gophers/dataloader/v7"
	"github.com/shota-tech/graphql/server/graph/model"
	"github.com/shota-tech/graphql/server/repository"
)

type TaskLoader struct {
	repository repository.ITaskRepository
}

func NewTaskLoader(repository repository.ITaskRepository) *TaskLoader {
	return &TaskLoader{
		repository: repository,
	}
}

func (l *TaskLoader) BulkGet(ctx context.Context, ids []string) []*dataloader.Result[*model.Task] {
	tasks, err := l.repository.List(ctx, ids)
	if err != nil {
		log.Printf("failed to list tasks: %v", err)
		return nil
	}

	taskByID := make(map[string]*model.Task, len(ids))
	for _, task := range tasks {
		taskByID[task.ID] = task
	}

	results := make([]*dataloader.Result[*model.Task], len(ids))
	for i, key := range ids {
		task, ok := taskByID[key]
		if ok {
			results[i] = &dataloader.Result[*model.Task]{Data: task}
		} else {
			results[i] = &dataloader.Result[*model.Task]{Error: fmt.Errorf("task not found: %s", key)}
		}
	}
	return results
}

func (l *TaskLoader) BulkGetByUserIDs(ctx context.Context, userIDs []string) []*dataloader.Result[[]*model.Task] {
	tasks, err := l.repository.ListByUserIDs(ctx, userIDs)
	if err != nil {
		log.Printf("failed to list tasks: %v", err)
		return nil
	}

	taskByUserID := make(map[string][]*model.Task, len(userIDs))
	for _, task := range tasks {
		taskByUserID[task.UserID] = append(taskByUserID[task.UserID], task)
	}

	results := make([]*dataloader.Result[[]*model.Task], len(userIDs))
	for i, userID := range userIDs {
		results[i] = &dataloader.Result[[]*model.Task]{
			Data: taskByUserID[userID],
		}
	}
	return results
}
