package graph

import "github.com/shota-tech/graphql/server/repository"

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct {
	TaskRepository repository.ITaskRepository
	UserRepository repository.IUserRepository
	TodoRepository repository.ITodoRepository
}
