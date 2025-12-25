package dto

type UpdateTodo struct {
	ID          int64   `json:"id"`
	Title       string  `json:"title"`
	Description *string `json:"description"`
	Completed   *bool   `json:"completed"`
}
