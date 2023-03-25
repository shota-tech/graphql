package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.
// Code generated by github.com/99designs/gqlgen version v0.17.25

import (
	"context"
	"errors"

	jwtMiddleware "github.com/auth0/go-jwt-middleware/v2"
	"github.com/auth0/go-jwt-middleware/v2/validator"
	"github.com/rs/xid"
	"github.com/shota-tech/graphql/server/graph/model"
	"github.com/shota-tech/graphql/server/middleware"
)

// CreateUser is the resolver for the createUser field.
func (r *mutationResolver) CreateUser(ctx context.Context, input model.CreateUserInput) (*model.User, error) {
	token := ctx.Value(jwtMiddleware.ContextKey{}).(*validator.ValidatedClaims)
	claims := token.CustomClaims.(*middleware.CustomClaims)
	if !claims.HasScope(middleware.ScopeWriteUser) {
		return nil, errors.New("invalid scope")
	}
	user := &model.User{
		ID:   token.RegisteredClaims.Subject,
		Name: input.Name,
	}
	if err := r.UserRepository.Store(ctx, user); err != nil {
		return nil, err
	}
	return user, nil
}

// CreateTask is the resolver for the createTask field.
func (r *mutationResolver) CreateTask(ctx context.Context, input model.CreateTaskInput) (*model.Task, error) {
	token := ctx.Value(jwtMiddleware.ContextKey{}).(*validator.ValidatedClaims)
	claims := token.CustomClaims.(*middleware.CustomClaims)
	if !claims.HasScope(middleware.ScopeWriteTasks) {
		return nil, errors.New("invalid scope")
	}
	task := &model.Task{
		ID:     xid.New().String(),
		Text:   input.Text,
		Status: model.StatusTodo,
		UserID: token.RegisteredClaims.Subject,
	}
	if err := r.TaskRepository.Store(ctx, task); err != nil {
		return nil, err
	}
	return task, nil
}

// UpdateTask is the resolver for the updateTask field.
func (r *mutationResolver) UpdateTask(ctx context.Context, input model.UpdateTaskInput) (*model.Task, error) {
	token := ctx.Value(jwtMiddleware.ContextKey{}).(*validator.ValidatedClaims)
	claims := token.CustomClaims.(*middleware.CustomClaims)
	if !claims.HasScope(middleware.ScopeWriteTasks) {
		return nil, errors.New("invalid scope")
	}
	task, err := r.TaskRepository.Get(ctx, input.ID)
	if err != nil {
		return nil, err
	}
	if input.Text != nil {
		task.Text = *input.Text
	}
	if input.Status != nil {
		task.Status = *input.Status
	}
	if err := r.TaskRepository.Store(ctx, task); err != nil {
		return nil, err
	}
	return task, nil
}

// CreateTodo is the resolver for the createTodo field.
func (r *mutationResolver) CreateTodo(ctx context.Context, input model.CreateTodoInput) (*model.Todo, error) {
	token := ctx.Value(jwtMiddleware.ContextKey{}).(*validator.ValidatedClaims)
	claims := token.CustomClaims.(*middleware.CustomClaims)
	if !claims.HasScope(middleware.ScopeWriteTasks) {
		return nil, errors.New("invalid scope")
	}
	todo := &model.Todo{
		ID:     xid.New().String(),
		Text:   input.Text,
		Done:   false,
		TaskID: input.TaskID,
	}
	if err := r.TodoRepository.Store(ctx, todo); err != nil {
		return nil, err
	}
	return todo, nil
}

// UpdateTodo is the resolver for the updateTodo field.
func (r *mutationResolver) UpdateTodo(ctx context.Context, input model.UpdateTodoInput) (*model.Todo, error) {
	token := ctx.Value(jwtMiddleware.ContextKey{}).(*validator.ValidatedClaims)
	claims := token.CustomClaims.(*middleware.CustomClaims)
	if !claims.HasScope(middleware.ScopeWriteTasks) {
		return nil, errors.New("invalid scope")
	}
	todo, err := r.TodoRepository.Get(ctx, input.ID)
	if err != nil {
		return nil, err
	}
	if input.Text != nil {
		todo.Text = *input.Text
	}
	if input.Done != nil {
		todo.Done = *input.Done
	}
	if err := r.TodoRepository.Store(ctx, todo); err != nil {
		return nil, err
	}
	return todo, nil
}

// Mutation returns MutationResolver implementation.
func (r *Resolver) Mutation() MutationResolver { return &mutationResolver{r} }

type mutationResolver struct{ *Resolver }
