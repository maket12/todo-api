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

func TestGetTodoUC(t *testing.T) {
	store := storage.NewDataStorage()
	uc := usecase.NewGetTodoUC(store)
	ctx := context.Background()

	t.Run("Success", func(t *testing.T) {
		fixedID := int64(10)
		testTodo := &entity.Todo{
			ID:          fixedID,
			Title:       "Wash clothes",
			Description: "",
			Completed:   true,
		}
		_ = store.CreateTodo(ctx, testTodo)

		result, err := uc.Execute(ctx, dto.GetTodo{ID: fixedID})
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if result.Title != testTodo.Title {
			t.Errorf("expected %v, got %v", testTodo.Title, result.Title)
		}
	})

	t.Run("Error - invalid ID", func(t *testing.T) {
		in := dto.GetTodo{ID: 0}
		if _, err := uc.Execute(ctx, in); !errors.Is(err, uc_errors.InvalidTodoIDError) {
			t.Errorf("expected InvalidTodoIDError, got %v", err)
		}
	})

	t.Run("Error - todo not found", func(t *testing.T) {
		in := dto.GetTodo{ID: 100}
		if _, err := uc.Execute(ctx, in); !errors.Is(err, uc_errors.TodoNotFoundError) {
			t.Errorf("expected TodoNotFoundError, got %v", err)
		}
	})

	t.Run("Error - canceled context", func(t *testing.T) {
		cancelCtx, cancel := context.WithCancel(context.Background())
		cancel()

		in := dto.GetTodo{ID: 100}
		if _, err := uc.Execute(cancelCtx, in); !errors.Is(err, uc_errors.GetTodoError) {
			t.Errorf("expected GetTodoError (canceled context), got %v", err)
		}
	})
}
