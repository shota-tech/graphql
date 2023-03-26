package loader

import (
	"github.com/graph-gophers/dataloader/v7"
	"github.com/shota-tech/graphql/server/graph/model"
)

type Loaders struct {
	TaskLoader         dataloader.Interface[string, *model.Task]
	TaskLoaderByUserID dataloader.Interface[string, []*model.Task]
	TodoLoaderByTaskID dataloader.Interface[string, []*model.Todo]
}

func NewLoaders(
	taskLoader *TaskLoader,
	todoLoader *TodoLoader,
) *Loaders {
	return &Loaders{
		TaskLoader: dataloader.NewBatchedLoader(
			taskLoader.BulkGet,
			dataloader.WithCache[string, *model.Task](
				&dataloader.NoCache[string, *model.Task]{},
			),
		),
		TaskLoaderByUserID: dataloader.NewBatchedLoader(
			taskLoader.BulkGetByUserIDs,
			dataloader.WithCache[string, []*model.Task](
				&dataloader.NoCache[string, []*model.Task]{},
			),
		),
		TodoLoaderByTaskID: dataloader.NewBatchedLoader(
			todoLoader.BulkGetByTaskIDs,
			dataloader.WithCache[string, []*model.Todo](
				&dataloader.NoCache[string, []*model.Todo]{},
			),
		),
	}
}
