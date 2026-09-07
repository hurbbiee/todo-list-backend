package dto

import "time"

type UpdateTodoRequest struct {
	Title       string    `json:"title" validate:"required"`
	Description string    `json:"description" validate:"required"`
	Status      string    `json:"status" validate:"required,oneof=pending in_progress completed"`
	Priority    string    `json:"priority" validate:"required,oneof=low normal high"`
	DueDate     time.Time `json:"dueDate" validate:"required"`
	UpdatedVia  string    `json:"updatedVia" validate:"required"`
}

type UpdateTodoStatus struct {
	Status     string `json:"status" validate:"required,oneof=pending in_progress completed"`
	UpdatedVia string `json:"updatedVia" validate:"required"`
}
