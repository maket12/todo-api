package dto

type DeleteTodoResponse struct {
	ID      int64 `json:"id"`
	Deleted bool  `json:"deleted"`
}
