package auth

import (
	"context"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	jwtMiddleware "github.com/auth0/go-jwt-middleware/v2"
	"github.com/auth0/go-jwt-middleware/v2/jwks"
	"github.com/auth0/go-jwt-middleware/v2/validator"
)

const (
	ScopeReadTasks  = "read:tasks"
	ScopeWriteTasks = "write:tasks"
	ScopeReadUser   = "read:user"
	ScopeWriteUser  = "write:user"
)

type CustomClaims struct {
	Scope string `json:"scope"`
}

func (c CustomClaims) Validate(ctx context.Context) error {
	return nil
}

func (c CustomClaims) HasScope(expectedScope string) bool {
	result := strings.Split(c.Scope, " ")
	for i := range result {
		if result[i] == expectedScope {
			return true
		}
	}
	return false
}

func EnsureValidToken() func(next http.Handler) http.Handler {
	issuerURL, err := url.Parse("https://" + os.Getenv("AUTH0_DOMAIN") + "/")
	if err != nil {
		log.Fatalf("failed to parse issuer url: %v", err)
	}
	provider := jwks.NewCachingProvider(issuerURL, 5*time.Minute)

	jwtValidator, err := validator.New(
		provider.KeyFunc,
		validator.RS256,
		issuerURL.String(),
		[]string{os.Getenv("AUTH0_AUDIENCE")},
		validator.WithCustomClaims(
			func() validator.CustomClaims {
				return &CustomClaims{}
			},
		),
		validator.WithAllowedClockSkew(time.Minute),
	)
	if err != nil {
		log.Fatal("failed to setup jwt validator")
	}

	errorHandler := func(w http.ResponseWriter, r *http.Request, err error) {
		log.Printf("encounterd error while validating JWT: %v", err)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte(`{"message":"failed to validate JWT."}`))
	}

	middleware := jwtMiddleware.New(
		jwtValidator.ValidateToken,
		jwtMiddleware.WithErrorHandler(errorHandler),
	)

	return func(next http.Handler) http.Handler {
		return middleware.CheckJWT(next)
	}
}

func TokenFromContext(ctx context.Context) *validator.ValidatedClaims {
	return ctx.Value(jwtMiddleware.ContextKey{}).(*validator.ValidatedClaims)
}
