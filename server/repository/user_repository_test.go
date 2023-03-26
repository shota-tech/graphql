package repository_test

import (
	"context"
	"database/sql/driver"
	"regexp"
	"testing"
	"time"

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
				query := "INSERT INTO `users` (`id`,`name`,`created_at`,`updated_at`) VALUES (?,?,?,?) " +
					"ON DUPLICATE KEY UPDATE `name` = VALUES(`name`),`created_at` = VALUES(`created_at`),`updated_at` = VALUES(`updated_at`)"
				args := []driver.Value{"auth0|123456", "user1", sqlmock.AnyArg(), sqlmock.AnyArg()}
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
				query := "INSERT INTO `users` (`id`,`name`,`created_at`,`updated_at`) VALUES (?,?,?,?) " +
					"ON DUPLICATE KEY UPDATE `name` = VALUES(`name`),`created_at` = VALUES(`created_at`),`updated_at` = VALUES(`updated_at`)"
				args := []driver.Value{"auth0|123456", "user1", sqlmock.AnyArg(), sqlmock.AnyArg()}
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
				query := "SELECT `users`.* FROM `users` WHERE (`users`.`id` = ?) LIMIT 1;"
				row := sqlmock.NewRows([]string{"id", "name", "created_at", "updated_at"}).
					AddRow("auth0|123456", "user1", time.Now(), time.Now())
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
				query := "SELECT `users`.* FROM `users` WHERE (`users`.`id` = ?) LIMIT 1;"
				row := sqlmock.NewRows([]string{"id", "name", "created_at", "updated_at"})
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
				query := "SELECT `users`.* FROM `users` WHERE (`users`.`id` = ?) LIMIT 1;"
				mock.ExpectQuery(regexp.QuoteMeta(query)).
					WithArgs("auth0|123456").
					WillReturnError(assert.AnError)
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

func TestUserRepository_List(t *testing.T) {
	tests := map[string]struct {
		setup     func(sqlmock.Sqlmock)
		ids       []string
		want      []*model.User
		assertErr assert.ErrorAssertionFunc
	}{
		"happy path": {
			setup: func(mock sqlmock.Sqlmock) {
				query := "SELECT `users`.* FROM `users` WHERE (`users`.`id` IN (?,?));"
				args := []driver.Value{"auth0|123456", "auth0|567890"}
				rows := sqlmock.NewRows([]string{"id", "name", "created_at", "updated_at"}).
					AddRow("auth0|123456", "user1", time.Now(), time.Now()).
					AddRow("auth0|567890", "user2", time.Now(), time.Now())
				mock.ExpectQuery(regexp.QuoteMeta(query)).
					WithArgs(args...).
					WillReturnRows(rows)
			},
			ids: []string{"auth0|123456", "auth0|567890"},
			want: []*model.User{
				{ID: "auth0|123456", Name: "user1"},
				{ID: "auth0|567890", Name: "user2"},
			},
			assertErr: assert.NoError,
		},
		"0 records": {
			setup: func(mock sqlmock.Sqlmock) {
				query := "SELECT `users`.* FROM `users` WHERE (`users`.`id` IN (?,?));"
				args := []driver.Value{"auth0|123456", "auth0|567890"}
				rows := sqlmock.NewRows([]string{"id", "name", "created_at", "updated_at"})
				mock.ExpectQuery(regexp.QuoteMeta(query)).
					WithArgs(args...).
					WillReturnRows(rows)
			},
			ids:       []string{"auth0|123456", "auth0|567890"},
			want:      []*model.User{},
			assertErr: assert.NoError,
		},
		"failed to get records": {
			setup: func(mock sqlmock.Sqlmock) {
				query := "SELECT `users`.* FROM `users` WHERE (`users`.`id` IN (?,?));"
				args := []driver.Value{"auth0|123456", "auth0|567890"}
				mock.ExpectQuery(regexp.QuoteMeta(query)).
					WithArgs(args...).
					WillReturnError(assert.AnError)
			},
			ids:       []string{"auth0|123456", "auth0|567890"},
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
			got, err := sut.List(context.Background(), tt.ids)
			assert.Equal(t, tt.want, got)
			tt.assertErr(t, err)
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}
