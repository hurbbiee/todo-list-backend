package dto

import (
	"time"
)

type CreateTodoRequest struct {
	Title       string    `json:"title" validate:"required"`
	Description string    `json:"description" validate:"required"`
	Status      string    `json:"status" validate:"required,oneof=pending in_progress completed"`
	Priority    string    `json:"priority" validate:"required,oneof=low normal high"`
	DueDate     time.Time `json:"dueDate" validate:"required"`
	CreatedVia  string    `json:"createdVia" validate:"required"`
}
