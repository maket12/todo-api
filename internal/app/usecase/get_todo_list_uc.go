package usecase

import (
	"context"
	"todo-api/internal/app/dto"
	"todo-api/internal/app/mappers"
	"todo-api/internal/app/uc_errors"
	"todo-api/internal/domain/port"
)

type GetTodoListUC struct {
	Storage port.DataStorage
}

func NewGetTodoListUC(storage port.DataStorage) *GetTodoListUC {
	return &GetTodoListUC{Storage: storage}
}

func (uc *GetTodoListUC) Execute(ctx context.Context, in dto.GetTodoList) (dto.GetTodoListResponse, error) {
	if in.Limit < 0 {
		return dto.GetTodoListResponse{}, uc_errors.InvalidLimitError
	}
	if in.Offset < 0 {
		return dto.GetTodoListResponse{}, uc_errors.InvalidOffsetError
	}

	todos, err := uc.Storage.GetTodoList(ctx, in.Limit, in.Offset)
	if err != nil {
		return dto.GetTodoListResponse{}, uc_errors.Wrap(uc_errors.GetTodoListError, err)
	}

	return mappers.MapDomainTodoListToTodoListDTO(todos), nil
}
