package usecase

import (
	"context"
	"errors"
	"todo-api/internal/app/dto"
	"todo-api/internal/app/mappers"
	"todo-api/internal/app/uc_errors"
	"todo-api/internal/domain/port"
)

type CreateTodoUC struct {
	Storage port.DataStorage
}

func NewCreateTodoUC(storage port.DataStorage) *CreateTodoUC {
	return &CreateTodoUC{Storage: storage}
}

func (uc *CreateTodoUC) Execute(ctx context.Context, in dto.CreateTodo) (dto.CreateTodoResponse, error) {
	if in.Title == "" {
		return dto.CreateTodoResponse{ID: in.ID}, uc_errors.EmptyTitleError
	}

	mappedIn := mappers.MapTodoDTOToDomainTodo(in.Todo)
	if err := uc.Storage.CreateTodo(ctx, mappedIn); err != nil {
		if !errors.Is(err, uc_errors.TodoAlreadyExistsError) {
			return dto.CreateTodoResponse{ID: mappedIn.ID}, uc_errors.Wrap(uc_errors.CreateTodoError, err)
		}
		return dto.CreateTodoResponse{ID: mappedIn.ID}, err
	}

	return dto.CreateTodoResponse{ID: mappedIn.ID}, nil
}
