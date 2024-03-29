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

func TestTodoRepository_Store(t *testing.T) {
	tests := map[string]struct {
		setup     func(sqlmock.Sqlmock)
		todo      *model.Todo
		assertErr assert.ErrorAssertionFunc
	}{
		"happy path": {
			setup: func(mock sqlmock.Sqlmock) {
				query := "INSERT INTO `todos` (`id`,`text`,`done`,`task_id`,`created_at`,`updated_at`) VALUES (?,?,?,?,?,?) " +
					"ON DUPLICATE KEY UPDATE `text` = VALUES(`text`),`done` = VALUES(`done`),`task_id` = VALUES(`task_id`),`created_at` = VALUES(`created_at`),`updated_at` = VALUES(`updated_at`)"
				args := []driver.Value{"cgf90odvqc7hkkh47tg0", "todo1", false, "cg1m0bd1nm6u7kpjp15g", sqlmock.AnyArg(), sqlmock.AnyArg()}
				mock.ExpectExec(regexp.QuoteMeta(query)).
					WithArgs(args...).
					WillReturnResult(sqlmock.NewResult(0, 1))
			},
			todo: &model.Todo{
				ID:     "cgf90odvqc7hkkh47tg0",
				Text:   "todo1",
				Done:   false,
				TaskID: "cg1m0bd1nm6u7kpjp15g",
			},
			assertErr: assert.NoError,
		},
		"todo is nil": {
			setup:     nil,
			todo:      nil,
			assertErr: assert.Error,
		},
		"failed to upsert record": {
			setup: func(mock sqlmock.Sqlmock) {
				query := "INSERT INTO `todos` (`id`,`text`,`done`,`task_id`,`created_at`,`updated_at`) VALUES (?,?,?,?,?,?) " +
					"ON DUPLICATE KEY UPDATE `text` = VALUES(`text`),`done` = VALUES(`done`),`task_id` = VALUES(`task_id`),`created_at` = VALUES(`created_at`),`updated_at` = VALUES(`updated_at`)"
				args := []driver.Value{"cgf90odvqc7hkkh47tg0", "todo1", false, "cg1m0bd1nm6u7kpjp15g", sqlmock.AnyArg(), sqlmock.AnyArg()}
				mock.ExpectExec(regexp.QuoteMeta(query)).
					WithArgs(args...).
					WillReturnError(assert.AnError)
			},
			todo: &model.Todo{
				ID:     "cgf90odvqc7hkkh47tg0",
				Text:   "todo1",
				Done:   false,
				TaskID: "cg1m0bd1nm6u7kpjp15g",
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

func TestTodoRepository_Get(t *testing.T) {
	tests := map[string]struct {
		setup     func(sqlmock.Sqlmock)
		id        string
		want      *model.Todo
		assertErr assert.ErrorAssertionFunc
	}{
		"happy path": {
			setup: func(mock sqlmock.Sqlmock) {
				query := "SELECT `todos`.* FROM `todos` WHERE (`todos`.`id` = ?) LIMIT 1;"
				row := sqlmock.NewRows([]string{"id", "text", "done", "task_id", "created_at", "updated_at"}).
					AddRow("cgf90odvqc7hkkh47tg0", "todo1", false, "cg1m0bd1nm6u7kpjp15g", time.Now(), time.Now())
				mock.ExpectQuery(regexp.QuoteMeta(query)).
					WithArgs("cgf90odvqc7hkkh47tg0").
					WillReturnRows(row)
			},
			id: "cgf90odvqc7hkkh47tg0",
			want: &model.Todo{
				ID:     "cgf90odvqc7hkkh47tg0",
				Text:   "todo1",
				Done:   false,
				TaskID: "cg1m0bd1nm6u7kpjp15g",
			},
			assertErr: assert.NoError,
		},
		"record not found": {
			setup: func(mock sqlmock.Sqlmock) {
				query := "SELECT `todos`.* FROM `todos` WHERE (`todos`.`id` = ?) LIMIT 1;"
				row := sqlmock.NewRows([]string{"id", "text", "done", "task_id", "created_at", "updated_at"})
				mock.ExpectQuery(regexp.QuoteMeta(query)).
					WithArgs("cgf90odvqc7hkkh47tg0").
					WillReturnRows(row)
			},
			id:        "cgf90odvqc7hkkh47tg0",
			want:      nil,
			assertErr: assert.Error,
		},
		"failed to get record": {
			setup: func(mock sqlmock.Sqlmock) {
				query := "SELECT `todos`.* FROM `todos` WHERE (`todos`.`id` = ?) LIMIT 1;"
				mock.ExpectQuery(regexp.QuoteMeta(query)).
					WithArgs("cgf90odvqc7hkkh47tg0").
					WillReturnError(assert.AnError)
			},
			id:        "cgf90odvqc7hkkh47tg0",
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
			got, err := sut.Get(context.Background(), tt.id)
			assert.Equal(t, tt.want, got)
			tt.assertErr(t, err)
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestTodoRepository_List(t *testing.T) {
	tests := map[string]struct {
		setup     func(sqlmock.Sqlmock)
		ids       []string
		want      []*model.Todo
		assertErr assert.ErrorAssertionFunc
	}{
		"happy path": {
			setup: func(mock sqlmock.Sqlmock) {
				query := "SELECT `todos`.* FROM `todos` WHERE (`todos`.`id` IN (?,?));"
				args := []driver.Value{"cgf90odvqc7hkkh47tg0", "cgf95atvqc7hriet4at0"}
				rows := sqlmock.NewRows([]string{"id", "text", "done", "task_id", "created_at", "updated_at"}).
					AddRow("cgf90odvqc7hkkh47tg0", "todo1", false, "cg1m0bd1nm6u7kpjp15g", time.Now(), time.Now()).
					AddRow("cgf95atvqc7hriet4at0", "todo2", true, "cg2j6hl1nm6ivqd084m0", time.Now(), time.Now())
				mock.ExpectQuery(regexp.QuoteMeta(query)).
					WithArgs(args...).
					WillReturnRows(rows)
			},
			ids: []string{"cgf90odvqc7hkkh47tg0", "cgf95atvqc7hriet4at0"},
			want: []*model.Todo{
				{ID: "cgf90odvqc7hkkh47tg0", Text: "todo1", Done: false, TaskID: "cg1m0bd1nm6u7kpjp15g"},
				{ID: "cgf95atvqc7hriet4at0", Text: "todo2", Done: true, TaskID: "cg2j6hl1nm6ivqd084m0"},
			},
			assertErr: assert.NoError,
		},
		"0 records": {
			setup: func(mock sqlmock.Sqlmock) {
				query := "SELECT `todos`.* FROM `todos` WHERE (`todos`.`id` IN (?,?));"
				args := []driver.Value{"cgf90odvqc7hkkh47tg0", "cgf95atvqc7hriet4at0"}
				rows := sqlmock.NewRows([]string{"id", "text", "done", "task_id", "created_at", "updated_at"})
				mock.ExpectQuery(regexp.QuoteMeta(query)).
					WithArgs(args...).
					WillReturnRows(rows)
			},
			ids:       []string{"cgf90odvqc7hkkh47tg0", "cgf95atvqc7hriet4at0"},
			want:      []*model.Todo{},
			assertErr: assert.NoError,
		},
		"failed to get records": {
			setup: func(mock sqlmock.Sqlmock) {
				query := "SELECT `todos`.* FROM `todos` WHERE (`todos`.`id` IN (?,?));"
				args := []driver.Value{"cgf90odvqc7hkkh47tg0", "cgf95atvqc7hriet4at0"}
				mock.ExpectQuery(regexp.QuoteMeta(query)).
					WithArgs(args...).
					WillReturnError(assert.AnError)
			},
			ids:       []string{"cgf90odvqc7hkkh47tg0", "cgf95atvqc7hriet4at0"},
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
			got, err := sut.List(context.Background(), tt.ids)
			assert.Equal(t, tt.want, got)
			tt.assertErr(t, err)
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestTodoRepository_ListByTaskIDs(t *testing.T) {
	tests := map[string]struct {
		setup     func(sqlmock.Sqlmock)
		taskIDs   []string
		want      []*model.Todo
		assertErr assert.ErrorAssertionFunc
	}{
		"happy path": {
			setup: func(mock sqlmock.Sqlmock) {
				query := "SELECT `todos`.* FROM `todos` WHERE (`todos`.`task_id` IN (?,?));"
				args := []driver.Value{"cg1m0bd1nm6u7kpjp15g", "cg2j6hl1nm6ivqd084m0"}
				rows := sqlmock.NewRows([]string{"id", "text", "done", "task_id", "created_at", "updated_at"}).
					AddRow("cgf90odvqc7hkkh47tg0", "todo1", false, "cg1m0bd1nm6u7kpjp15g", time.Now(), time.Now()).
					AddRow("cgf95atvqc7hriet4at0", "todo2", true, "cg2j6hl1nm6ivqd084m0", time.Now(), time.Now())
				mock.ExpectQuery(regexp.QuoteMeta(query)).
					WithArgs(args...).
					WillReturnRows(rows)
			},
			taskIDs: []string{"cg1m0bd1nm6u7kpjp15g", "cg2j6hl1nm6ivqd084m0"},
			want: []*model.Todo{
				{ID: "cgf90odvqc7hkkh47tg0", Text: "todo1", Done: false, TaskID: "cg1m0bd1nm6u7kpjp15g"},
				{ID: "cgf95atvqc7hriet4at0", Text: "todo2", Done: true, TaskID: "cg2j6hl1nm6ivqd084m0"},
			},
			assertErr: assert.NoError,
		},
		"0 records": {
			setup: func(mock sqlmock.Sqlmock) {
				query := "SELECT `todos`.* FROM `todos` WHERE (`todos`.`task_id` IN (?,?));"
				args := []driver.Value{"cg1m0bd1nm6u7kpjp15g", "cg2j6hl1nm6ivqd084m0"}
				rows := sqlmock.NewRows([]string{"id", "text", "done", "task_id", "created_at", "updated_at"})
				mock.ExpectQuery(regexp.QuoteMeta(query)).
					WithArgs(args...).
					WillReturnRows(rows)
			},
			taskIDs:   []string{"cg1m0bd1nm6u7kpjp15g", "cg2j6hl1nm6ivqd084m0"},
			want:      []*model.Todo{},
			assertErr: assert.NoError,
		},
		"failed to get records": {
			setup: func(mock sqlmock.Sqlmock) {
				query := "SELECT `todos`.* FROM `todos` WHERE (`todos`.`task_id` IN (?,?));"
				args := []driver.Value{"cg1m0bd1nm6u7kpjp15g", "cg2j6hl1nm6ivqd084m0"}
				mock.ExpectQuery(regexp.QuoteMeta(query)).
					WithArgs(args...).
					WillReturnError(assert.AnError)
			},
			taskIDs:   []string{"cg1m0bd1nm6u7kpjp15g", "cg2j6hl1nm6ivqd084m0"},
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
			got, err := sut.ListByTaskIDs(context.Background(), tt.taskIDs)
			assert.Equal(t, tt.want, got)
			tt.assertErr(t, err)
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}
