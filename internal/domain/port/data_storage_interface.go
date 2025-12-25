package port

import (
	"context"
	"todo-api/internal/domain/entity"
)

type DataStorage interface {
	CreateTodo(ctx context.Context, todo *entity.Todo) error
	GetTodo(ctx context.Context, id int64) (*entity.Todo, error)
	GetTodoList(ctx context.Context, limit, offset int) ([]*entity.Todo, error)
	UpdateTodo(ctx context.Context, todo *entity.Todo) error
	DeleteTodo(ctx context.Context, id int64) error
}
