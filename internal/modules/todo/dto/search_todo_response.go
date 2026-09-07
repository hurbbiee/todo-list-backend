package dto

import "time"

type SearchTodoResponse struct {
	Id          int64     `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Status      string    `json:"status"`
	Priority    string    `json:"priority"`
	DueDate     time.Time `json:"dueDate"`
	CompletedAt *time.Time `json:"completedAt"`
	CreatedVia  string    `json:"createdVia"`
}
