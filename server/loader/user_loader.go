package loader

import (
	"context"
	"fmt"
	"log"

	"github.com/graph-gophers/dataloader/v7"
	"github.com/shota-tech/graphql/server/graph/model"
	"github.com/shota-tech/graphql/server/repository"
)

type UserLoader struct {
	repository repository.IUserRepository
}

func NewUserLoader(repository repository.IUserRepository) *UserLoader {
	return &UserLoader{
		repository: repository,
	}
}

func (l *UserLoader) BulkGet(ctx context.Context, ids []string) []*dataloader.Result[*model.User] {
	users, err := l.repository.List(ctx, ids)
	if err != nil {
		log.Printf("failed to list users: %v", err)
		return nil
	}

	userByID := make(map[string]*model.User, len(ids))
	for _, user := range users {
		userByID[user.ID] = user
	}

	results := make([]*dataloader.Result[*model.User], len(ids))
	for i, key := range ids {
		user, ok := userByID[key]
		if ok {
			results[i] = &dataloader.Result[*model.User]{Data: user}
		} else {
			results[i] = &dataloader.Result[*model.User]{Error: fmt.Errorf("user not found: %s", key)}
		}
	}
	return results
}
