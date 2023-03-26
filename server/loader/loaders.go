package loader

import (
	"github.com/graph-gophers/dataloader/v7"
	"github.com/shota-tech/graphql/server/graph/model"
)

type Loaders struct {
	TodoLoaderByTaskID dataloader.Interface[string, []*model.Todo]
}

func NewLoaders(todoLoader *TodoLoader) *Loaders {
	return &Loaders{
		TodoLoaderByTaskID: dataloader.NewBatchedLoader(
			todoLoader.BulkGetByTaskIDs,
			dataloader.WithCache[string, []*model.Todo](&dataloader.NoCache[string, []*model.Todo]{}),
		),
	}
}
