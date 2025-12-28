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

func TestUpdateTodoUC(t *testing.T) {
	store := storage.NewDataStorage()
	uc := NewUpdateTodoUC(store)
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

		in := dto.UpdateTodo{
			Todo: dto.Todo{
				ID:          fixedID,
				Title:       "Wash dishes",
				Description: "use a washing machine",
				Completed:   false,
			},
		}
		if _, err := uc.Execute(ctx, in); err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		if updated, _ := store.GetTodo(ctx, fixedID); updated.Title != in.Title {
			t.Errorf("expected title %s, got %s", in.Title, updated.Title)
		}
	})

	t.Run("Error - invalid id", func(t *testing.T) {
		in := dto.UpdateTodo{Todo: dto.Todo{ID: 0}}
		if _, err := uc.Execute(ctx, in); !errors.Is(err, uc_errors.InvalidTodoIDError) {
			t.Errorf("expected InvalidTodoIDError, got %v", err)
		}
	})

	t.Run("Error - empty title", func(t *testing.T) {
		in := dto.UpdateTodo{Todo: dto.Todo{ID: 10, Title: ""}}
		if _, err := uc.Execute(ctx, in); !errors.Is(err, uc_errors.EmptyTitleError) {
			t.Errorf("expected EmptyTitleError, got %v", err)
		}
	})

	t.Run("Error - todo not found", func(t *testing.T) {
		in := dto.UpdateTodo{Todo: dto.Todo{ID: 100, Title: "New title"}}
		if _, err := uc.Execute(ctx, in); !errors.Is(err, uc_errors.TodoNotFoundError) {
			t.Errorf("expected TodoNotFoundError, got %v", err)
		}
	})

	t.Run("Error - canceled context", func(t *testing.T) {
		cancelCtx, cancel := context.WithCancel(context.Background())
		cancel()

		in := dto.UpdateTodo{Todo: dto.Todo{ID: 10, Title: "New title"}}
		if _, err := uc.Execute(cancelCtx, in); !errors.Is(err, uc_errors.GetTodoError) && !errors.Is(err, uc_errors.UpdateTodoError) {
			t.Errorf("expected GetTodoError | UpdateTodoError (canceled context), got %v", err)
		}
	})
}
