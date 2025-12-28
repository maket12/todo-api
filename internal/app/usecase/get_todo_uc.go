package usecase

import (
	"context"
	"errors"
	"todo-api/internal/app/dto"
	"todo-api/internal/app/mappers"
	"todo-api/internal/app/uc_errors"
	"todo-api/internal/domain/port"
)

type GetTodoUC struct {
	Storage port.DataStorage
}

func NewGetTodoUC(storage port.DataStorage) *GetTodoUC {
	return &GetTodoUC{Storage: storage}
}

func (uc *GetTodoUC) Execute(ctx context.Context, in dto.GetTodo) (dto.GetTodoResponse, error) {
	if in.ID <= 0 {
		return dto.GetTodoResponse{Todo: dto.Todo{ID: in.ID}}, uc_errors.InvalidTodoIDError
	}

	todo, err := uc.Storage.GetTodo(ctx, in.ID)
	if err != nil {
		if !errors.Is(err, uc_errors.TodoNotFoundError) {
			return dto.GetTodoResponse{Todo: dto.Todo{ID: in.ID}}, uc_errors.Wrap(uc_errors.GetTodoError, err)
		}
		return dto.GetTodoResponse{Todo: dto.Todo{ID: in.ID}}, err
	}

	return dto.GetTodoResponse{
		Todo: mappers.MapDomainTodoToTodoDTO(todo),
	}, nil
}
