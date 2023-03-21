package repository_test

import (
	"context"
	"database/sql/driver"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/shota-tech/graphql/server/graph/model"
	"github.com/shota-tech/graphql/server/repository"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUserRepository_Store(t *testing.T) {
	tests := map[string]struct {
		setup     func(sqlmock.Sqlmock)
		user      *model.User
		assertErr assert.ErrorAssertionFunc
	}{
		"happy path": {
			setup: func(mock sqlmock.Sqlmock) {
				query := "INSERT INTO users (id, name) VALUES (?, ?) ON DUPLICATE KEY UPDATE name = VALUES(name);"
				args := []driver.Value{"auth0|123456", "user1"}
				mock.ExpectExec(regexp.QuoteMeta(query)).
					WithArgs(args...).
					WillReturnResult(sqlmock.NewResult(0, 1))
			},
			user: &model.User{
				ID:   "auth0|123456",
				Name: "user1",
			},
			assertErr: assert.NoError,
		},
		"user is nil": {
			user:      nil,
			setup:     nil,
			assertErr: assert.Error,
		},
		"failed to upsert record": {
			setup: func(mock sqlmock.Sqlmock) {
				query := "INSERT INTO users (id, name) VALUES (?, ?) ON DUPLICATE KEY UPDATE name = VALUES(name);"
				args := []driver.Value{"auth0|123456", "user1"}
				mock.ExpectExec(regexp.QuoteMeta(query)).
					WithArgs(args...).
					WillReturnError(assert.AnError)
			},
			user: &model.User{
				ID:   "auth0|123456",
				Name: "user1",
			},
			assertErr: assert.Error,
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			// setup sqlmock
			db, mock, err := sqlmock.New()
			require.NoError(t, err)
			defer db.Close()
			if tt.setup != nil {
				tt.setup(mock)
			}
			// test
			sut := repository.NewUserRepository(db)
			err = sut.Store(context.Background(), tt.user)
			tt.assertErr(t, err)
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestUserRepository_Get(t *testing.T) {
	tests := map[string]struct {
		setup     func(sqlmock.Sqlmock)
		id        string
		want      *model.User
		assertErr assert.ErrorAssertionFunc
	}{
		"happy path": {
			setup: func(mock sqlmock.Sqlmock) {
				query := "SELECT id, name FROM users WHERE id = ?;"
				row := sqlmock.NewRows([]string{"id", "name"}).
					AddRow("auth0|123456", "user1")
				mock.ExpectQuery(regexp.QuoteMeta(query)).
					WithArgs("auth0|123456").
					WillReturnRows(row)
			},
			id:        "auth0|123456",
			want:      &model.User{ID: "auth0|123456", Name: "user1"},
			assertErr: assert.NoError,
		},
		"record not found": {
			setup: func(mock sqlmock.Sqlmock) {
				query := "SELECT id, name FROM users WHERE id = ?;"
				row := sqlmock.NewRows([]string{"id", "name"})
				mock.ExpectQuery(regexp.QuoteMeta(query)).
					WithArgs("auth0|123456").
					WillReturnRows(row)
			},
			id:        "auth0|123456",
			want:      nil,
			assertErr: assert.Error,
		},
		"failed to get record": {
			setup: func(mock sqlmock.Sqlmock) {
				query := "SELECT id, name FROM users WHERE id = ?;"
				mock.ExpectQuery(regexp.QuoteMeta(query)).
					WithArgs("auth0|123456").
					WillReturnError(assert.AnError)
			},
			id:        "auth0|123456",
			want:      nil,
			assertErr: assert.Error,
		},
		"failed to scan record": {
			setup: func(mock sqlmock.Sqlmock) {
				query := "SELECT id, name FROM users WHERE id = ?;"
				row := sqlmock.NewRows([]string{"id", "name"}).
					AddRow("auth0|123456", nil)
				mock.ExpectQuery(regexp.QuoteMeta(query)).
					WithArgs("auth0|123456").
					WillReturnRows(row)
			},
			id:        "auth0|123456",
			want:      nil,
			assertErr: assert.Error,
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			// setup sqlmock
			db, mock, err := sqlmock.New()
			require.NoError(t, err)
			defer db.Close()
			if tt.setup != nil {
				tt.setup(mock)
			}
			// test
			sut := repository.NewUserRepository(db)
			got, err := sut.Get(context.Background(), tt.id)
			assert.Equal(t, tt.want, got)
			tt.assertErr(t, err)
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}
