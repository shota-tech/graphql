package loader

import (
	"context"
	"fmt"
	"log"

	"github.com/graph-gophers/dataloader/v7"
	"github.com/shota-tech/graphql/server/graph/model"
	"github.com/shota-tech/graphql/server/repository"
)

type TodoLoader struct {
	repository repository.ITodoRepository
}

func NewTodoLoader(repository repository.ITodoRepository) *TodoLoader {
	return &TodoLoader{
		repository: repository,
	}
}

func (l *TodoLoader) BulkGet(ctx context.Context, ids []string) []*dataloader.Result[*model.Todo] {
	todo, err := l.repository.List(ctx, ids)
	if err != nil {
		log.Printf("failed to list todos: %v", err)
		return nil
	}

	todoByID := make(map[string]*model.Todo, len(ids))
	for _, todo := range todo {
		todoByID[todo.ID] = todo
	}

	results := make([]*dataloader.Result[*model.Todo], len(ids))
	for i, key := range ids {
		todo, ok := todoByID[key]
		if ok {
			results[i] = &dataloader.Result[*model.Todo]{Data: todo}
		} else {
			results[i] = &dataloader.Result[*model.Todo]{Error: fmt.Errorf("todo not found: %s", key)}
		}
	}
	return results
}

func (l *TodoLoader) BulkGetByTaskIDs(ctx context.Context, taskIDs []string) []*dataloader.Result[[]*model.Todo] {
	todos, err := l.repository.ListByTaskIDs(ctx, taskIDs)
	if err != nil {
		log.Printf("failed to list todos: %v", err)
		return nil
	}

	todoByTaskID := make(map[string][]*model.Todo, len(taskIDs))
	for _, todo := range todos {
		todoByTaskID[todo.TaskID] = append(todoByTaskID[todo.TaskID], todo)
	}

	results := make([]*dataloader.Result[[]*model.Todo], len(taskIDs))
	for i, taskID := range taskIDs {
		results[i] = &dataloader.Result[[]*model.Todo]{
			Data: todoByTaskID[taskID],
		}
	}
	return results
}
