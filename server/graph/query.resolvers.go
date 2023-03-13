package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.
// Code generated by github.com/99designs/gqlgen version v0.17.25

import (
	"context"

	"github.com/shota-tech/graphql/server/graph/model"
)

// FetchTasks is the resolver for the fetchTasks field.
func (r *queryResolver) FetchTasks(ctx context.Context, userID string) ([]*model.Task, error) {
	return r.TaskRepository.ListByUserID(ctx, userID)
}

// FetchUser is the resolver for the fetchUser field.
func (r *queryResolver) FetchUser(ctx context.Context, id string) (*model.User, error) {
	return r.UserRepository.Get(ctx, id)
}

// Query returns QueryResolver implementation.
func (r *Resolver) Query() QueryResolver { return &queryResolver{r} }

type queryResolver struct{ *Resolver }
