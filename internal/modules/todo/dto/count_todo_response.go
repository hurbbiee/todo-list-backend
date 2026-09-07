package dto

type CountTodoResponse struct {
	All        int64 `json:"all"`
	Pending    int64 `json:"pending"`
	InProgress int64 `json:"inProgress"`
	Completed  int64 `json:"completed"`
}
