package dto

type UpsertDiscordConnectionRequest struct {
	WebhookURL          string `json:"webhookUrl" validate:"required,url"`
	IsEnabled           *bool  `json:"isEnabled" validate:"required"`
	NotifyTodoCreated   *bool  `json:"notifyTodoCreated" validate:"required"`
	NotifyTodoCompleted *bool  `json:"notifyTodoCompleted" validate:"required"`
	NotifyBeforeDue     *bool  `json:"notifyBeforeDue" validate:"required"`
	RemindBeforeMinutes int    `json:"remindBeforeMinutes" validate:"required,min=1,max=10080"`
}

type UpdateDiscordEnabledRequest struct {
	IsEnabled *bool `json:"isEnabled" validate:"required"`
}

type DiscordConnectionResponse struct {
	Connected           bool `json:"connected"`
	IsEnabled           bool `json:"isEnabled"`
	NotifyTodoCreated   bool `json:"notifyTodoCreated"`
	NotifyTodoCompleted bool `json:"notifyTodoCompleted"`
	NotifyBeforeDue     bool `json:"notifyBeforeDue"`
	RemindBeforeMinutes int  `json:"remindBeforeMinutes"`
}
