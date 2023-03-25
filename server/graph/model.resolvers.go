package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.
// Code generated by github.com/99designs/gqlgen version v0.17.25

import (
	"context"
	"errors"

	jwtMiddleware "github.com/auth0/go-jwt-middleware/v2"
	"github.com/auth0/go-jwt-middleware/v2/validator"
	"github.com/shota-tech/graphql/server/graph/model"
	"github.com/shota-tech/graphql/server/middleware"
)

// User is the resolver for the user field.
func (r *taskResolver) User(ctx context.Context, obj *model.Task) (*model.User, error) {
	token := ctx.Value(jwtMiddleware.ContextKey{}).(*validator.ValidatedClaims)
	claims := token.CustomClaims.(*middleware.CustomClaims)
	if !claims.HasScope(middleware.ScopeReadUser) {
		return nil, errors.New("invalid scope")
	}
	return r.UserRepository.Get(ctx, obj.UserID)
}

// Todos is the resolver for the todos field.
func (r *taskResolver) Todos(ctx context.Context, obj *model.Task) ([]*model.Todo, error) {
	token := ctx.Value(jwtMiddleware.ContextKey{}).(*validator.ValidatedClaims)
	claims := token.CustomClaims.(*middleware.CustomClaims)
	if !claims.HasScope(middleware.ScopeReadTasks) {
		return nil, errors.New("invalid scope")
	}
	return r.TodoRepository.ListByTaskID(ctx, obj.ID)
}

// Task is the resolver for the task field.
func (r *todoResolver) Task(ctx context.Context, obj *model.Todo) (*model.Task, error) {
	token := ctx.Value(jwtMiddleware.ContextKey{}).(*validator.ValidatedClaims)
	claims := token.CustomClaims.(*middleware.CustomClaims)
	if !claims.HasScope(middleware.ScopeReadTasks) {
		return nil, errors.New("invalid scope")
	}
	return r.TaskRepository.Get(ctx, obj.TaskID)
}

// Tasks is the resolver for the tasks field.
func (r *userResolver) Tasks(ctx context.Context, obj *model.User) ([]*model.Task, error) {
	token := ctx.Value(jwtMiddleware.ContextKey{}).(*validator.ValidatedClaims)
	claims := token.CustomClaims.(*middleware.CustomClaims)
	if !claims.HasScope(middleware.ScopeReadTasks) {
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
