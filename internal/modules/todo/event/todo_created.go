package event

import "time"

const TodoCreatedRoutingKey = "todo.created"

type TodoCreated struct {
	EventType   string    `json:"eventType"`
	TodoID      int64     `json:"todoId"`
	UserID      int64     `json:"userId"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Status      string    `json:"status"`
	Priority    string    `json:"priority"`
	DueAt       time.Time `json:"dueAt"`
	CreatedVia  string    `json:"createdVia"`
	OccurredAt  time.Time `json:"occurredAt"`
}