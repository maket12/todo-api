package usecase_test

import (
	"context"
	"errors"
	"testing"
	"todo-api/internal/adapter/out/storage"
	"todo-api/internal/app/dto"
	"todo-api/internal/app/uc_errors"
	"todo-api/internal/app/usecase"
	"todo-api/internal/domain/entity"
)

func TestGetTodoListUC(t *testing.T) {
	store := storage.NewDataStorage()
	uc := usecase.NewGetTodoListUC(store)
	ctx := context.Background()

	t.Run("Success", func(t *testing.T) {
		todo1 := entity.Todo{Title: "Get something"}
		todo2 := entity.Todo{Title: "Go somewhere"}
		todo3 := entity.Todo{Title: "Do somehow"}
		_ = store.CreateTodo(ctx, &todo1)
		_ = store.CreateTodo(ctx, &todo2)
		_ = store.CreateTodo(ctx, &todo3)

		testLimit, testOffset := 2, 1

		result, err := uc.Execute(ctx, dto.GetTodoList{
			Limit:  testLimit,
			Offset: testOffset,
		})
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if len(result.Todos) != 2 {
			t.Errorf("expected 2 elements, got %v", result.Todos)
		}
		if result.Todos[0].Title != todo2.Title {
			t.Errorf("expected %v, got %v", todo2.Title, result.Todos[0].Title)
		}
	})

	t.Run("Error - invalid limit", func(t *testing.T) {
		in := dto.GetTodoList{Limit: -1}
		if _, err := uc.Execute(ctx, in); !errors.Is(err, uc_errors.InvalidLimitError) {
			t.Errorf("expected InvalidLimitError, got %v", err)
		}
	})

	t.Run("Error - invalid offset", func(t *testing.T) {
		in := dto.GetTodoList{Offset: -1}
		if _, err := uc.Execute(ctx, in); !errors.Is(err, uc_errors.InvalidOffsetError) {
			t.Errorf("expected InvalidOffsetError, got %v", err)
		}
	})

	t.Run("Error - canceled context", func(t *testing.T) {
		cancelCtx, cancel := context.WithCancel(context.Background())
		cancel()

		in := dto.GetTodoList{}
		if _, err := uc.Execute(cancelCtx, in); !errors.Is(err, uc_errors.GetTodoListError) {
			t.Errorf("expected GetTodoListError (canceled context), got %v", err)
		}
	})
}
