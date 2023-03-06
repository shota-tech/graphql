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
		user      *model.User
		setup     func(sqlmock.Sqlmock)
		assertErr assert.ErrorAssertionFunc
	}{
		"happy path": {
			user: &model.User{
				ID:   "cg1ltn51nm6u7l352ma0",
				Name: "user1",
			},
			setup: func(mock sqlmock.Sqlmock) {
				query := "INSERT INTO users (id, name) VALUES (?, ?) ON DUPLICATE KEY UPDATE name = VALUES(name);"
				args := []driver.Value{"cg1ltn51nm6u7l352ma0", "user1"}
				mock.ExpectExec(regexp.QuoteMeta(query)).
					WithArgs(args...).
					WillReturnResult(sqlmock.NewResult(0, 1))
			},
			assertErr: assert.NoError,
		},
		"user is nil": {
			user:      nil,
			assertErr: assert.Error,
		},
		"failed to upsert record": {
			user: &model.User{
				ID:   "cg1ltn51nm6u7l352ma0",
				Name: "user1",
			},
			setup: func(mock sqlmock.Sqlmock) {
				query := "INSERT INTO users (id, name) VALUES (?, ?) ON DUPLICATE KEY UPDATE name = VALUES(name);"
				args := []driver.Value{"cg1ltn51nm6u7l352ma0", "user1"}
				mock.ExpectExec(regexp.QuoteMeta(query)).
					WithArgs(args...).
					WillReturnError(assert.AnError)
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
		id        string
		setup     func(sqlmock.Sqlmock)
		want      *model.User
		assertErr assert.ErrorAssertionFunc
	}{
		"happy path": {
			id: "cg1ltn51nm6u7l352ma0",
			setup: func(mock sqlmock.Sqlmock) {
				query := "SELECT id, name FROM users WHERE id = ?;"
				row := sqlmock.NewRows([]string{"id", "name"}).
					AddRow("cg1ltn51nm6u7l352ma0", "user1")
				mock.ExpectQuery(regexp.QuoteMeta(query)).
					WillReturnRows(row)
			},
			want:      &model.User{ID: "cg1ltn51nm6u7l352ma0", Name: "user1"},
			assertErr: assert.NoError,
		},
		"record not found": {
			id: "cg1ltn51nm6u7l352ma0",
			setup: func(mock sqlmock.Sqlmock) {
				query := "SELECT id, name FROM users WHERE id = ?;"
				row := sqlmock.NewRows([]string{"id", "name"})
				mock.ExpectQuery(regexp.QuoteMeta(query)).
					WillReturnRows(row)
			},
			want:      nil,
			assertErr: assert.Error,
		},
		"failed to get record": {
			id: "cg1ltn51nm6u7l352ma0",
			setup: func(mock sqlmock.Sqlmock) {
				query := "SELECT id, name FROM users WHERE id = ?;"
				mock.ExpectQuery(regexp.QuoteMeta(query)).
					WillReturnError(assert.AnError)
			},
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
