package usecase

import (
	"context"
	"errors"
	"testing"
	"todo-api/internal/adapter/out/storage"
	"todo-api/internal/app/dto"
	"todo-api/internal/app/uc_errors"
	"todo-api/internal/domain/entity"
)

func TestDeleteTodoUC(t *testing.T) {
	store := storage.NewDataStorage()
	uc := NewDeleteTodoUC(store)
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

		in := dto.DeleteTodo{ID: fixedID}
		result, err := uc.Execute(ctx, in)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if result.ID != fixedID {
			t.Errorf("expected deletion of %v, actual %v", fixedID, result.ID)
		}

		if _, err := store.GetTodo(ctx, fixedID); !errors.Is(err, uc_errors.TodoNotFoundError) {
			t.Errorf("deleted todo %v found in storage", fixedID)
		}
	})

	t.Run("Error - invalid ID", func(t *testing.T) {
		in := dto.DeleteTodo{ID: 0}
		if _, err := uc.Execute(ctx, in); !errors.Is(err, uc_errors.InvalidTodoIDError) {
			t.Errorf("expected InvalidTodoIDError, got %v", err)
		}
	})

	t.Run("Error - todo not found", func(t *testing.T) {
		in := dto.DeleteTodo{ID: 100}
		if _, err := uc.Execute(ctx, in); !errors.Is(err, uc_errors.TodoNotFoundError) {
			t.Errorf("expected TodoNotFoundError, got %v", err)
		}
	})

	t.Run("Error - canceled context", func(t *testing.T) {
		cancelCtx, cancel := context.WithCancel(context.Background())
		cancel()

		in := dto.DeleteTodo{ID: 100}
		if _, err := uc.Execute(cancelCtx, in); !errors.Is(err, uc_errors.DeleteTodoError) {
			t.Errorf("expected DeleteTodoError (canceled context), got %v", err)
		}
	})
}
