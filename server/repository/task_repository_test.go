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

func TestTaskRepository_Store(t *testing.T) {
	tests := map[string]struct {
		setup     func(sqlmock.Sqlmock)
		task      *model.Task
		assertErr assert.ErrorAssertionFunc
	}{
		"happy path": {
			setup: func(mock sqlmock.Sqlmock) {
				query := "INSERT INTO tasks (id, text, status, user_id) VALUES (?, ?, ?, ?) " +
					"ON DUPLICATE KEY UPDATE text = VALUES(text), status = VALUES(status), user_id = VALUES(user_id);"
				args := []driver.Value{"cg1m0bd1nm6u7kpjp15g", "task1", "TODO", "auth0|123456"}
				mock.ExpectExec(regexp.QuoteMeta(query)).
					WithArgs(args...).
					WillReturnResult(sqlmock.NewResult(0, 1))
			},
			task: &model.Task{
				ID:     "cg1m0bd1nm6u7kpjp15g",
				Text:   "task1",
				Status: model.StatusTodo,
				UserID: "auth0|123456",
			},
			assertErr: assert.NoError,
		},
		"task is nil": {
			setup:     nil,
			task:      nil,
			assertErr: assert.Error,
		},
		"failed to upsert record": {
			setup: func(mock sqlmock.Sqlmock) {
				query := "INSERT INTO tasks (id, text, status, user_id) VALUES (?, ?, ?, ?) " +
					"ON DUPLICATE KEY UPDATE text = VALUES(text), status = VALUES(status), user_id = VALUES(user_id);"
				args := []driver.Value{"cg1m0bd1nm6u7kpjp15g", "task1", "TODO", "auth0|123456"}
				mock.ExpectExec(regexp.QuoteMeta(query)).
					WithArgs(args...).
					WillReturnError(assert.AnError)
			},
			task: &model.Task{
				ID:     "cg1m0bd1nm6u7kpjp15g",
				Text:   "task1",
				Status: model.StatusTodo,
				UserID: "auth0|123456",
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
			sut := repository.NewTaskRepository(db)
			err = sut.Store(context.Background(), tt.task)
			tt.assertErr(t, err)
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestTaskRepository_Get(t *testing.T) {
	tests := map[string]struct {
		setup     func(sqlmock.Sqlmock)
		id        string
		want      *model.Task
		assertErr assert.ErrorAssertionFunc
	}{
		"happy path": {
			setup: func(mock sqlmock.Sqlmock) {
				query := "SELECT id, text, status, user_id FROM tasks WHERE id = ?;"
				row := sqlmock.NewRows([]string{"id", "text", "status", "user_id"}).
					AddRow("cg1m0bd1nm6u7kpjp15g", "task1", "TODO", "auth0|123456")
				mock.ExpectQuery(regexp.QuoteMeta(query)).
					WithArgs("cg1m0bd1nm6u7kpjp15g").
					WillReturnRows(row)
			},
			id: "cg1m0bd1nm6u7kpjp15g",
			want: &model.Task{
				ID:     "cg1m0bd1nm6u7kpjp15g",
				Text:   "task1",
				Status: model.StatusTodo,
				UserID: "auth0|123456",
			},
			assertErr: assert.NoError,
		},
		"record not found": {
			setup: func(mock sqlmock.Sqlmock) {
				query := "SELECT id, text, status, user_id FROM tasks WHERE id = ?;"
				row := sqlmock.NewRows([]string{"id", "text", "status", "user_id"})
				mock.ExpectQuery(regexp.QuoteMeta(query)).
					WithArgs("cg1m0bd1nm6u7kpjp15g").
					WillReturnRows(row)
			},
			id:        "cg1m0bd1nm6u7kpjp15g",
			want:      nil,
			assertErr: assert.Error,
		},
		"failed to get record": {
			setup: func(mock sqlmock.Sqlmock) {
				query := "SELECT id, text, status, user_id FROM tasks WHERE id = ?;"
				mock.ExpectQuery(regexp.QuoteMeta(query)).
					WithArgs("cg1m0bd1nm6u7kpjp15g").
					WillReturnError(assert.AnError)
			},
			id:        "cg1m0bd1nm6u7kpjp15g",
			want:      nil,
			assertErr: assert.Error,
		},
		"failed to scan record": {
			setup: func(mock sqlmock.Sqlmock) {
				query := "SELECT id, text, status, user_id FROM tasks WHERE id = ?;"
				row := sqlmock.NewRows([]string{"id", "text", "status", "user_id"}).
					AddRow("cg1m0bd1nm6u7kpjp15g", nil, "TODO", "auth0|123456")
				mock.ExpectQuery(regexp.QuoteMeta(query)).
					WithArgs("cg1m0bd1nm6u7kpjp15g").
					WillReturnRows(row)
			},
			id:        "cg1m0bd1nm6u7kpjp15g",
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
			sut := repository.NewTaskRepository(db)
			got, err := sut.Get(context.Background(), tt.id)
			assert.Equal(t, tt.want, got)
			tt.assertErr(t, err)
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestTaskRepository_List(t *testing.T) {
	tests := map[string]struct {
		setup     func(sqlmock.Sqlmock)
		ids       []string
		want      []*model.Task
		assertErr assert.ErrorAssertionFunc
	}{
		"happy path": {
			setup: func(mock sqlmock.Sqlmock) {
				query := "SELECT id, text, status, user_id FROM tasks WHERE id IN (?,?);"
				args := []driver.Value{"cg1m0bd1nm6u7kpjp15g", "cg2j6hl1nm6ivqd084m0"}
				rows := sqlmock.NewRows([]string{"id", "text", "status", "user_id"}).
					AddRow("cg1m0bd1nm6u7kpjp15g", "task1", "TODO", "auth0|123456").
					AddRow("cg2j6hl1nm6ivqd084m0", "task2", "TODO", "auth0|567890")
				mock.ExpectQuery(regexp.QuoteMeta(query)).
					WithArgs(args...).
					WillReturnRows(rows)
			},
			ids: []string{"cg1m0bd1nm6u7kpjp15g", "cg2j6hl1nm6ivqd084m0"},
			want: []*model.Task{
				{ID: "cg1m0bd1nm6u7kpjp15g", Text: "task1", Status: model.StatusTodo, UserID: "auth0|123456"},
				{ID: "cg2j6hl1nm6ivqd084m0", Text: "task2", Status: model.StatusTodo, UserID: "auth0|567890"},
			},
			assertErr: assert.NoError,
		},
		"0 records": {
			setup: func(mock sqlmock.Sqlmock) {
				query := "SELECT id, text, status, user_id FROM tasks WHERE id IN (?,?);"
				args := []driver.Value{"cg1m0bd1nm6u7kpjp15g", "cg2j6hl1nm6ivqd084m0"}
				rows := sqlmock.NewRows([]string{"id", "text", "status", "user_id"})
				mock.ExpectQuery(regexp.QuoteMeta(query)).
					WithArgs(args...).
					WillReturnRows(rows)
			},
			ids:       []string{"cg1m0bd1nm6u7kpjp15g", "cg2j6hl1nm6ivqd084m0"},
			want:      []*model.Task{},
			assertErr: assert.NoError,
		},
		"failed to get records": {
			setup: func(mock sqlmock.Sqlmock) {
				query := "SELECT id, text, status, user_id FROM tasks WHERE id IN (?,?);"
				args := []driver.Value{"cg1m0bd1nm6u7kpjp15g", "cg2j6hl1nm6ivqd084m0"}
				mock.ExpectQuery(regexp.QuoteMeta(query)).
					WithArgs(args...).
					WillReturnError(assert.AnError)
			},
			ids:       []string{"cg1m0bd1nm6u7kpjp15g", "cg2j6hl1nm6ivqd084m0"},
			want:      nil,
			assertErr: assert.Error,
		},
		"failed to scan record": {
			setup: func(mock sqlmock.Sqlmock) {
				query := "SELECT id, text, status, user_id FROM tasks WHERE id IN (?,?);"
				args := []driver.Value{"cg1m0bd1nm6u7kpjp15g", "cg2j6hl1nm6ivqd084m0"}
				rows := sqlmock.NewRows([]string{"id", "text", "status", "user_id"}).
					AddRow("cg1m0bd1nm6u7kpjp15g", "task1", "TODO", "auth0|123456").
					AddRow("cg2j6hl1nm6ivqd084m0", nil, "TODO", "auth0|567890")
				mock.ExpectQuery(regexp.QuoteMeta(query)).
					WithArgs(args...).
					WillReturnRows(rows)
			},
			ids:       []string{"cg1m0bd1nm6u7kpjp15g", "cg2j6hl1nm6ivqd084m0"},
			want:      nil,
			assertErr: assert.Error,
		},
		"failed to scan records": {
			setup: func(mock sqlmock.Sqlmock) {
				query := "SELECT id, text, status, user_id FROM tasks WHERE id IN (?,?);"
				args := []driver.Value{"cg1m0bd1nm6u7kpjp15g", "cg2j6hl1nm6ivqd084m0"}
				rows := sqlmock.NewRows([]string{"id", "text", "status", "user_id"}).
					AddRow("cg1m0bd1nm6u7kpjp15g", "task1", "TODO", "auth0|123456").
					AddRow("cg2j6hl1nm6ivqd084m0", "task2", "TODO", "auth0|567890").
					RowError(1, assert.AnError)
				mock.ExpectQuery(regexp.QuoteMeta(query)).
					WithArgs(args...).
					WillReturnRows(rows)
			},
			ids:       []string{"cg1m0bd1nm6u7kpjp15g", "cg2j6hl1nm6ivqd084m0"},
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
			sut := repository.NewTaskRepository(db)
			got, err := sut.List(context.Background(), tt.ids)
			assert.Equal(t, tt.want, got)
			tt.assertErr(t, err)
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestTaskRepository_ListByUserID(t *testing.T) {
	tests := map[string]struct {
		setup     func(sqlmock.Sqlmock)
		userID    string
		want      []*model.Task
		assertErr assert.ErrorAssertionFunc
	}{
		"happy path": {
			setup: func(mock sqlmock.Sqlmock) {
				query := "SELECT id, text, status, user_id FROM tasks WHERE user_id = ?;"
				args := []driver.Value{"auth0|123456"}
				rows := sqlmock.NewRows([]string{"id", "text", "status", "user_id"}).
					AddRow("cg1m0bd1nm6u7kpjp15g", "task1", "TODO", "auth0|123456").
					AddRow("cg2j6hl1nm6ivqd084m0", "task2", "TODO", "auth0|123456")
				mock.ExpectQuery(regexp.QuoteMeta(query)).
					WithArgs(args...).
					WillReturnRows(rows)
			},
			userID: "auth0|123456",
			want: []*model.Task{
				{ID: "cg1m0bd1nm6u7kpjp15g", Text: "task1", Status: model.StatusTodo, UserID: "auth0|123456"},
				{ID: "cg2j6hl1nm6ivqd084m0", Text: "task2", Status: model.StatusTodo, UserID: "auth0|123456"},
			},
			assertErr: assert.NoError,
		},
		"0 records": {
			setup: func(mock sqlmock.Sqlmock) {
				query := "SELECT id, text, status, user_id FROM tasks WHERE user_id = ?;"
				args := []driver.Value{"auth0|123456"}
				rows := sqlmock.NewRows([]string{"id", "text", "status", "user_id"})
				mock.ExpectQuery(regexp.QuoteMeta(query)).
					WithArgs(args...).
					WillReturnRows(rows)
			},
			userID:    "auth0|123456",
			want:      []*model.Task{},
			assertErr: assert.NoError,
		},
		"failed to get records": {
			setup: func(mock sqlmock.Sqlmock) {
				query := "SELECT id, text, status, user_id FROM tasks WHERE user_id = ?;"
				args := []driver.Value{"auth0|123456"}
				mock.ExpectQuery(regexp.QuoteMeta(query)).
					WithArgs(args...).
					WillReturnError(assert.AnError)
			},
			userID:    "auth0|123456",
			want:      nil,
			assertErr: assert.Error,
		},
		"failed to scan record": {
			setup: func(mock sqlmock.Sqlmock) {
				query := "SELECT id, text, status, user_id FROM tasks WHERE user_id = ?;"
				args := []driver.Value{"auth0|123456"}
				rows := sqlmock.NewRows([]string{"id", "text", "status", "user_id"}).
					AddRow("cg1m0bd1nm6u7kpjp15g", "task1", "TODO", "auth0|123456").
					AddRow("cg2j6hl1nm6ivqd084m0", nil, "TODO", "auth0|123456")
				mock.ExpectQuery(regexp.QuoteMeta(query)).
					WithArgs(args...).
					WillReturnRows(rows)
			},
			userID:    "auth0|123456",
			want:      nil,
			assertErr: assert.Error,
		},
		"failed to scan records": {
			setup: func(mock sqlmock.Sqlmock) {
				query := "SELECT id, text, status, user_id FROM tasks WHERE user_id = ?;"
				args := []driver.Value{"auth0|123456"}
				rows := sqlmock.NewRows([]string{"id", "text", "status", "user_id"}).
					AddRow("cg1m0bd1nm6u7kpjp15g", "task1", "TODO", "auth0|123456").
					AddRow("cg2j6hl1nm6ivqd084m0", "task2", "TODO", "auth0|123456").
					RowError(1, assert.AnError)
				mock.ExpectQuery(regexp.QuoteMeta(query)).
					WithArgs(args...).
					WillReturnRows(rows)
			},
			userID:    "auth0|123456",
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
			sut := repository.NewTaskRepository(db)
			got, err := sut.ListByUserID(context.Background(), tt.userID)
			assert.Equal(t, tt.want, got)
			tt.assertErr(t, err)
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestTaskRepository_ListByUserIDs(t *testing.T) {
	tests := map[string]struct {
		setup     func(sqlmock.Sqlmock)
		userIDs   []string
		want      []*model.Task
		assertErr assert.ErrorAssertionFunc
	}{
		"happy path": {
			setup: func(mock sqlmock.Sqlmock) {
				query := "SELECT id, text, status, user_id FROM tasks WHERE user_id IN (?,?);"
				args := []driver.Value{"auth0|123456", "auth0|567890"}
				rows := sqlmock.NewRows([]string{"id", "text", "status", "user_id"}).
					AddRow("cg1m0bd1nm6u7kpjp15g", "task1", "TODO", "auth0|123456").
					AddRow("cg2j6hl1nm6ivqd084m0", "task2", "TODO", "auth0|567890")
				mock.ExpectQuery(regexp.QuoteMeta(query)).
					WithArgs(args...).
					WillReturnRows(rows)
			},
			userIDs: []string{"auth0|123456", "auth0|567890"},
			want: []*model.Task{
				{ID: "cg1m0bd1nm6u7kpjp15g", Text: "task1", Status: model.StatusTodo, UserID: "auth0|123456"},
				{ID: "cg2j6hl1nm6ivqd084m0", Text: "task2", Status: model.StatusTodo, UserID: "auth0|567890"},
			},
			assertErr: assert.NoError,
		},
		"0 records": {
			setup: func(mock sqlmock.Sqlmock) {
				query := "SELECT id, text, status, user_id FROM tasks WHERE user_id IN (?,?);"
				args := []driver.Value{"auth0|123456", "auth0|567890"}
				rows := sqlmock.NewRows([]string{"id", "text", "status", "user_id"})
				mock.ExpectQuery(regexp.QuoteMeta(query)).
					WithArgs(args...).
					WillReturnRows(rows)
			},
			userIDs:   []string{"auth0|123456", "auth0|567890"},
			want:      []*model.Task{},
			assertErr: assert.NoError,
		},
		"failed to get records": {
			setup: func(mock sqlmock.Sqlmock) {
				query := "SELECT id, text, status, user_id FROM tasks WHERE user_id IN (?,?);"
				args := []driver.Value{"auth0|123456", "auth0|567890"}
				mock.ExpectQuery(regexp.QuoteMeta(query)).
					WithArgs(args...).
					WillReturnError(assert.AnError)
			},
			userIDs:   []string{"auth0|123456", "auth0|567890"},
			want:      nil,
			assertErr: assert.Error,
		},
		"failed to scan record": {
			setup: func(mock sqlmock.Sqlmock) {
				query := "SELECT id, text, status, user_id FROM tasks WHERE user_id IN (?,?);"
				args := []driver.Value{"auth0|123456", "auth0|567890"}
				rows := sqlmock.NewRows([]string{"id", "text", "status", "user_id"}).
					AddRow("cg1m0bd1nm6u7kpjp15g", "task1", "TODO", "auth0|123456").
					AddRow("cg2j6hl1nm6ivqd084m0", nil, "TODO", "auth0|567890")
				mock.ExpectQuery(regexp.QuoteMeta(query)).
					WithArgs(args...).
					WillReturnRows(rows)
			},
			userIDs:   []string{"auth0|123456", "auth0|567890"},
			want:      nil,
			assertErr: assert.Error,
		},
		"failed to scan records": {
			setup: func(mock sqlmock.Sqlmock) {
				query := "SELECT id, text, status, user_id FROM tasks WHERE user_id IN (?,?);"
				args := []driver.Value{"auth0|123456", "auth0|567890"}
				rows := sqlmock.NewRows([]string{"id", "text", "status", "user_id"}).
					AddRow("cg1m0bd1nm6u7kpjp15g", "task1", "TODO", "auth0|123456").
					AddRow("cg2j6hl1nm6ivqd084m0", "task2", "TODO", "auth0|567890").
					RowError(1, assert.AnError)
				mock.ExpectQuery(regexp.QuoteMeta(query)).
					WithArgs(args...).
					WillReturnRows(rows)
			},
			userIDs:   []string{"auth0|123456", "auth0|567890"},
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
			sut := repository.NewTaskRepository(db)
			got, err := sut.ListByUserIDs(context.Background(), tt.userIDs)
			assert.Equal(t, tt.want, got)
			tt.assertErr(t, err)
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}
