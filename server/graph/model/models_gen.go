// Code generated by github.com/99designs/gqlgen, DO NOT EDIT.

package model

import (
	"fmt"
	"io"
	"strconv"
)

type CreateTaskInput struct {
	Text string `json:"text"`
}

type CreateTodoInput struct {
	Text   string `json:"text"`
	TaskID string `json:"taskID"`
}

type CreateUserInput struct {
	Name string `json:"name"`
}

type UpdateTaskInput struct {
	ID     string  `json:"id"`
	Text   *string `json:"text"`
	Status *Status `json:"status"`
}

type UpdateTodoInput struct {
	ID   string  `json:"id"`
	Text *string `json:"text"`
	Done *bool   `json:"done"`
}

type Status string

const (
	StatusTodo       Status = "TODO"
	StatusInProgress Status = "IN_PROGRESS"
	StatusDone       Status = "DONE"
)

var AllStatus = []Status{
	StatusTodo,
	StatusInProgress,
	StatusDone,
}

func (e Status) IsValid() bool {
	switch e {
	case StatusTodo, StatusInProgress, StatusDone:
		return true
	}
	return false
}

func (e Status) String() string {
	return string(e)
}

func (e *Status) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*e = Status(str)
	if !e.IsValid() {
		return fmt.Errorf("%s is not a valid Status", str)
	}
	return nil
}

func (e Status) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(e.String()))
}
