package usecase

import (
	"context"
	"errors"
	"todo-api/internal/app/dto"
	"todo-api/internal/app/uc_errors"
	"todo-api/internal/domain/port"
)

type DeleteTodoUC struct {
	Storage port.DataStorage
}

func NewDeleteTodoUC(storage port.DataStorage) *DeleteTodoUC {
	return &DeleteTodoUC{Storage: storage}
}

func (uc *DeleteTodoUC) Execute(ctx context.Context, in dto.DeleteTodo) (dto.DeleteTodoResponse, error) {
	if in.ID <= 0 {
		return dto.DeleteTodoResponse{ID: in.ID}, uc_errors.InvalidTodoIDError
	}

	if err := uc.Storage.DeleteTodo(ctx, in.ID); err != nil {
		if !errors.Is(err, uc_errors.TodoNotFoundError) {
			return dto.DeleteTodoResponse{ID: in.ID}, uc_errors.Wrap(uc_errors.DeleteTodoError, err)
		}
		return dto.DeleteTodoResponse{ID: in.ID}, err
	}

	return dto.DeleteTodoResponse{
		ID:      in.ID,
		Deleted: true,
	}, nil
}
