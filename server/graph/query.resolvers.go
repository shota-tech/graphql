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

// FetchTasks is the resolver for the fetchTasks field.
func (r *queryResolver) FetchTasks(ctx context.Context) ([]*model.Task, error) {
	token := ctx.Value(jwtMiddleware.ContextKey{}).(*validator.ValidatedClaims)
	claims := token.CustomClaims.(middleware.CustomClaims)
	if claims.HasScope(middleware.ScopeReadTasks) {
		return nil, errors.New("invalid scope")
	}
	return r.TaskRepository.ListByUserID(ctx, token.RegisteredClaims.Subject)
}

// FetchUser is the resolver for the fetchUser field.
func (r *queryResolver) FetchUser(ctx context.Context) (*model.User, error) {
	token := ctx.Value(jwtMiddleware.ContextKey{}).(*validator.ValidatedClaims)
	return r.UserRepository.Get(ctx, token.RegisteredClaims.Subject)
}

// Query returns QueryResolver implementation.
func (r *Resolver) Query() QueryResolver { return &queryResolver{r} }

type queryResolver struct{ *Resolver }
