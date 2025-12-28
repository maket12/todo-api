package uc_errors

import "errors"

var (
	InvalidTodoIDError     = errors.New("todo id must be positive digit")
	EmptyTitleError        = errors.New("empty todo title")
	InvalidLimitError      = errors.New("limit must be a positive digit or 0")
	InvalidOffsetError     = errors.New("offset must be a positive digit or 0")
	TodoNotFoundError      = errors.New("todo with this id is not found")
	TodoAlreadyExistsError = errors.New("todo with this id already exists")
	CreateTodoError        = errors.New("failed to create todo")
	GetTodoError           = errors.New("failed to get todo")
	GetTodoListError       = errors.New("failed to get todo list")
	UpdateTodoError        = errors.New("failed to update todo")
	DeleteTodoError        = errors.New("failed to delete todo")
)
