package usecase

import (
	"context"
	"errors"
	"todo-api/internal/app/dto"
	"todo-api/internal/app/uc_errors"
	"todo-api/internal/domain/port"
)

type UpdateTodoUC struct {
	Storage port.DataStorage
}

func NewUpdateTodoUC(storage port.DataStorage) *UpdateTodoUC {
	return &UpdateTodoUC{Storage: storage}
}

func (uc *UpdateTodoUC) Execute(ctx context.Context, in dto.UpdateTodo) (dto.UpdateTodoResponse, error) {
	if in.ID <= 0 {
		return dto.UpdateTodoResponse{ID: in.ID}, uc_errors.InvalidTodoIDError
	}
	if in.Title == "" {
		return dto.UpdateTodoResponse{ID: in.ID}, uc_errors.EmptyTitleError
	}

	currentTodo, err := uc.Storage.GetTodo(ctx, in.ID)
	if err != nil {
		if !errors.Is(err, uc_errors.TodoNotFoundError) {
			return dto.UpdateTodoResponse{ID: in.ID}, uc_errors.Wrap(uc_errors.GetTodoError, err)
		}
		return dto.UpdateTodoResponse{ID: in.ID}, err
	}

	currentTodo.Title = in.Title
	if in.Description != nil {
		currentTodo.Description = *in.Description
	}
	if in.Completed != nil {
		currentTodo.Completed = *in.Completed
	}

	err = uc.Storage.UpdateTodo(ctx, currentTodo)
	if err != nil {
		if !errors.Is(err, uc_errors.TodoNotFoundError) {
			return dto.UpdateTodoResponse{ID: currentTodo.ID}, uc_errors.Wrap(uc_errors.UpdateTodoError, err)
		}
		return dto.UpdateTodoResponse{ID: currentTodo.ID}, err
	}

	return dto.UpdateTodoResponse{
		ID:      currentTodo.ID,
		Updated: true,
	}, nil
}
