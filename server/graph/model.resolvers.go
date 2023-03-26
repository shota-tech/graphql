package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.
// Code generated by github.com/99designs/gqlgen version v0.17.25

import (
	"context"
	"errors"

	"github.com/shota-tech/graphql/server/graph/model"
	"github.com/shota-tech/graphql/server/middleware/auth"
)

// User is the resolver for the user field.
func (r *taskResolver) User(ctx context.Context, obj *model.Task) (*model.User, error) {
	token := auth.TokenFromContext(ctx)
	claims := token.CustomClaims.(*auth.CustomClaims)
	if !claims.HasScope(auth.ScopeReadUser) {
		return nil, errors.New("invalid scope")
	}
	return r.UserRepository.Get(ctx, obj.UserID)
}

// Todos is the resolver for the todos field.
func (r *taskResolver) Todos(ctx context.Context, obj *model.Task) ([]*model.Todo, error) {
	token := auth.TokenFromContext(ctx)
	claims := token.CustomClaims.(*auth.CustomClaims)
	if !claims.HasScope(auth.ScopeReadTasks) {
		return nil, errors.New("invalid scope")
	}
	thunk := r.Loaders.TodoLoaderByTaskID.Load(ctx, obj.ID)
	return thunk()
}

// Task is the resolver for the task field.
func (r *todoResolver) Task(ctx context.Context, obj *model.Todo) (*model.Task, error) {
	token := auth.TokenFromContext(ctx)
	claims := token.CustomClaims.(*auth.CustomClaims)
	if !claims.HasScope(auth.ScopeReadTasks) {
		return nil, errors.New("invalid scope")
	}
	return r.TaskRepository.Get(ctx, obj.TaskID)
}

// Tasks is the resolver for the tasks field.
func (r *userResolver) Tasks(ctx context.Context, obj *model.User) ([]*model.Task, error) {
	token := auth.TokenFromContext(ctx)
	claims := token.CustomClaims.(*auth.CustomClaims)
	if !claims.HasScope(auth.ScopeReadTasks) {
		return nil, errors.New("invalid scope")
	}
	return r.TaskRepository.ListByUserID(ctx, obj.ID)
}

// Task returns TaskResolver implementation.
func (r *Resolver) Task() TaskResolver { return &taskResolver{r} }

// Todo returns TodoResolver implementation.
func (r *Resolver) Todo() TodoResolver { return &todoResolver{r} }

// User returns UserResolver implementation.
func (r *Resolver) User() UserResolver { return &userResolver{r} }

type taskResolver struct{ *Resolver }
type todoResolver struct{ *Resolver }
type userResolver struct{ *Resolver }
