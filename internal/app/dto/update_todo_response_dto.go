package dto

type UpdateTodoResponse struct {
	ID      int64 `json:"id"`
	Updated bool  `json:"updated"`
}
