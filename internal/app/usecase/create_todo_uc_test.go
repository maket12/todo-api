package usecase_test

import (
	"context"
	"errors"
	"testing"
	"todo-api/internal/adapter/out/storage"
	"todo-api/internal/app/dto"
	"todo-api/internal/app/uc_errors"
	"todo-api/internal/app/usecase"
)

func TestCreateTodoUC(t *testing.T) {
	store := storage.NewDataStorage()
	uc := usecase.NewCreateTodoUC(store)
	ctx := context.Background()

	t.Run("Success", func(t *testing.T) {
		in := dto.CreateTodo{
			Todo: dto.Todo{
				ID:          0,
				Title:       "Learn math",
				Description: "Lineal algebra",
				Completed:   false,
			},
		}

		result, err := uc.Execute(ctx, in)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if result.ID == 0 {
			t.Error("expected auto-generated ID, got 0")
		}

		saved, getErr := store.GetTodo(ctx, result.ID)
		if getErr != nil {
			t.Fatalf("could not find created todo in storage: %v", getErr)
		}
		if saved.Title != in.Title {
			t.Errorf("expected title %s, got %s", in.Title, saved.Title)
		}
	})

	t.Run("Error - empty title", func(t *testing.T) {
		in := dto.CreateTodo{Todo: dto.Todo{Title: ""}}
		if _, err := uc.Execute(ctx, in); !errors.Is(err, uc_errors.EmptyTitleError) {
			t.Errorf("expected EmptyTitleError, got %v", err)
		}
	})

	t.Run("Error - duplicate id", func(t *testing.T) {
		testID := int64(200)
		in := dto.CreateTodo{
			Todo: dto.Todo{
				ID:          testID,
				Title:       "Test task",
				Description: "Nothing else",
				Completed:   true,
			},
		}
		_, _ = uc.Execute(ctx, in)

		_, err := uc.Execute(ctx, in)
		if !errors.Is(err, uc_errors.TodoAlreadyExistsError) {
			t.Errorf("expected TodoAlreadyExistsError, got %v", err)
		}
	})

	t.Run("Error - canceled context", func(t *testing.T) {
		cancelCtx, cancel := context.WithCancel(context.Background())
		cancel()

		in := dto.CreateTodo{Todo: dto.Todo{Title: "Cancel me!"}}
		if _, err := uc.Execute(cancelCtx, in); !errors.Is(err, uc_errors.CreateTodoError) {
			t.Errorf("expected CreateTodoError (canceled context), got %v", err)
		}
	})
}
