package loader

import (
	"github.com/graph-gophers/dataloader/v7"
	"github.com/shota-tech/graphql/server/graph/model"
)

type Loaders struct {
	UserLoader         dataloader.Interface[string, *model.User]
	TaskLoader         dataloader.Interface[string, *model.Task]
	TodoLoader         dataloader.Interface[string, *model.Todo]
	TaskLoaderByUserID dataloader.Interface[string, []*model.Task]
	TodoLoaderByTaskID dataloader.Interface[string, []*model.Todo]
}

func NewLoaders(
	userLoader *UserLoader,
	taskLoader *TaskLoader,
	todoLoader *TodoLoader,
) *Loaders {
	return &Loaders{
		UserLoader: dataloader.NewBatchedLoader(
			userLoader.BulkGet,
			dataloader.WithCache[string, *model.User](
				&dataloader.NoCache[string, *model.User]{},
			),
		),
		TaskLoader: dataloader.NewBatchedLoader(
			taskLoader.BulkGet,
			dataloader.WithCache[string, *model.Task](
				&dataloader.NoCache[string, *model.Task]{},
			),
		),
		TodoLoader: dataloader.NewBatchedLoader(
			todoLoader.BulkGet,
			dataloader.WithCache[string, *model.Todo](
				&dataloader.NoCache[string, *model.Todo]{},
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
