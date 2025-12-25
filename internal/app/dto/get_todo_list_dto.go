package dto

type GetTodoList struct {
	Limit  int `json:"limit"`
	Offset int `json:"offset"`
}
