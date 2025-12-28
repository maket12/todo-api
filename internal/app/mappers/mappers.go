package mappers

import (
	"todo-api/internal/app/dto"
	"todo-api/internal/domain/entity"
)

func MapTodoDTOToDomainTodo(input dto.Todo) *entity.Todo {
	return &entity.Todo{
		ID:          input.ID,
		Title:       input.Title,
		Description: input.Description,
		Completed:   input.Completed,
	}
}

func MapDomainTodoToTodoDTO(input *entity.Todo) dto.Todo {
	return dto.Todo{
		ID:          input.ID,
		Title:       input.Title,
		Description: input.Description,
		Completed:   input.Completed,
	}
}

func MapDomainTodoListToTodoListDTO(input []*entity.Todo) dto.GetTodoListResponse {
	todos := make([]dto.Todo, len(input))
	for i := 0; i < len(input); i++ {
		todos[i] = dto.Todo{
			ID:          input[i].ID,
			Title:       input[i].Title,
			Description: input[i].Description,
			Completed:   input[i].Completed,
		}
	}
	return dto.GetTodoListResponse{Todos: todos}
}
