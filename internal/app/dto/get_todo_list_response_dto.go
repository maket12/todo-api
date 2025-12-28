package dto

type GetTodoListResponse struct {
	Todos []Todo `json:"items"`
}
