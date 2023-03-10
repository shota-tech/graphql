package model

type Task struct {
	ID     string `json:"id"`
	Text   string `json:"text"`
	Status Status `json:"status"`
	UserID string `json:"userId"`
}
