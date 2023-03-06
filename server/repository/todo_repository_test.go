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

func TestTodoRepository_Store(t *testing.T) {
	tests := map[string]struct {
		todo      *model.Todo
		setup     func(sqlmock.Sqlmock)
		assertErr assert.ErrorAssertionFunc
	}{
		"happy path": {
			todo: &model.Todo{
				ID:     "cg1m0bd1nm6u7kpjp15g",
				Text:   "todo1",
				Done:   false,
				UserID: "cg1ltn51nm6u7l352ma0",
			},
			setup: func(mock sqlmock.Sqlmock) {
				query := "INSERT INTO todos (id, text, done, user_id) VALUES (?, ?, ?, ?) " +
					"ON DUPLICATE KEY UPDATE text = VALUES(text), done = VALUES(done), user_id = VALUES(user_id);"
				args := []driver.Value{"cg1m0bd1nm6u7kpjp15g", "todo1", false, "cg1ltn51nm6u7l352ma0"}
				mock.ExpectExec(regexp.QuoteMeta(query)).
					WithArgs(args...).
					WillReturnResult(sqlmock.NewResult(0, 1))
			},
			assertErr: assert.NoError,
		},
		"todo is nil": {
			todo:      nil,
			assertErr: assert.Error,
		},
		"failed to upsert record": {
			todo: &model.Todo{
				ID:     "cg1m0bd1nm6u7kpjp15g",
				Text:   "todo1",
				Done:   false,
				UserID: "cg1ltn51nm6u7l352ma0",
			},
			setup: func(mock sqlmock.Sqlmock) {
				query := "INSERT INTO todos (id, text, done, user_id) VALUES (?, ?, ?, ?) " +
					"ON DUPLICATE KEY UPDATE text = VALUES(text), done = VALUES(done), user_id = VALUES(user_id);"
				args := []driver.Value{"cg1m0bd1nm6u7kpjp15g", "todo1", false, "cg1ltn51nm6u7l352ma0"}
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
			sut := repository.NewTodoRepository(db)
			err = sut.Store(context.Background(), tt.todo)
			tt.assertErr(t, err)
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestTodoRepository_List(t *testing.T) {
	tests := map[string]struct {
		setup     func(sqlmock.Sqlmock)
		want      []*model.Todo
		assertErr assert.ErrorAssertionFunc
	}{
		"happy path": {
			setup: func(mock sqlmock.Sqlmock) {
				query := "SELECT id, text, done, user_id FROM todos;"
				rows := sqlmock.NewRows([]string{"id", "text", "done", "user_id"}).
					AddRow("cg1m0bd1nm6u7kpjp15g", "todo1", false, "cg1ltn51nm6u7l352ma0").
					AddRow("cg2j6hl1nm6ivqd084m0", "todo2", false, "cg1ltn51nm6u7l352ma0")
				mock.ExpectQuery(regexp.QuoteMeta(query)).WillReturnRows(rows)
			},
			want: []*model.Todo{
				{ID: "cg1m0bd1nm6u7kpjp15g", Text: "todo1", Done: false, UserID: "cg1ltn51nm6u7l352ma0"},
				{ID: "cg2j6hl1nm6ivqd084m0", Text: "todo2", Done: false, UserID: "cg1ltn51nm6u7l352ma0"},
			},
			assertErr: assert.NoError,
		},
		"0 records": {
			setup: func(mock sqlmock.Sqlmock) {
				query := "SELECT id, text, done, user_id FROM todos;"
				rows := sqlmock.NewRows([]string{"id", "text", "done", "user_id"})
				mock.ExpectQuery(regexp.QuoteMeta(query)).WillReturnRows(rows)
			},
			want:      []*model.Todo{},
			assertErr: assert.NoError,
		},
		"failed to get records": {
			setup: func(mock sqlmock.Sqlmock) {
				query := "SELECT id, text, done, user_id FROM todos;"
				mock.ExpectQuery(regexp.QuoteMeta(query)).
					WillReturnError(assert.AnError)
			},
			want:      nil,
			assertErr: assert.Error,
		},
		"failed to scan record": {
			setup: func(mock sqlmock.Sqlmock) {
				query := "SELECT id, text, done, user_id FROM todos;"
				rows := sqlmock.NewRows([]string{"id", "text", "done", "user_id"}).
					AddRow("cg1m0bd1nm6u7kpjp15g", "todo1", false, "cg1ltn51nm6u7l352ma0").
					AddRow("cg2j6hl1nm6ivqd084m0", nil, false, "cg1ltn51nm6u7l352ma0")
				mock.ExpectQuery(regexp.QuoteMeta(query)).WillReturnRows(rows)
			},
			want:      nil,
			assertErr: assert.Error,
		},
		"failed to scan records": {
			setup: func(mock sqlmock.Sqlmock) {
				query := "SELECT id, text, done, user_id FROM todos;"
				rows := sqlmock.NewRows([]string{"id", "text", "done", "user_id"}).
					AddRow("cg1m0bd1nm6u7kpjp15g", "todo1", false, "cg1ltn51nm6u7l352ma0").
					AddRow("cg2j6hl1nm6ivqd084m0", "todo2", false, "cg1ltn51nm6u7l352ma0").
					RowError(1, assert.AnError)
				mock.ExpectQuery(regexp.QuoteMeta(query)).WillReturnRows(rows)
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
			sut := repository.NewTodoRepository(db)
			got, err := sut.List(context.Background())
			assert.Equal(t, tt.want, got)
			tt.assertErr(t, err)
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}
