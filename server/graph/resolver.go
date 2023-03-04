package graph

import "github.com/shota-tech/graphql/server/repository"

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct {
	TodoRepository repository.ITodoRepository
	UserRepository repository.IUserRepository
}
