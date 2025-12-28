package storage

import (
	"context"
	"errors"
	"testing"
	"todo-api/internal/app/uc_errors"
	"todo-api/internal/domain/entity"
)

func TestStorage_CreateTodo(t *testing.T) {
	s := NewDataStorage()
	ctx := context.Background()

	t.Run("Success", func(t *testing.T) {
		todo := entity.Todo{
			ID:          0,
			Title:       "Get a coffee",
			Description: "Get an ice-latte in Starbucks",
			Completed:   false,
		}
		if err := s.CreateTodo(ctx, &todo); err != nil {
			t.Fatalf("expected no error, got %v", err)
			return
		}
		if todo.ID == 0 {
			t.Error("expected auto-gen id, got 0")
		}
	})

	t.Run("Context cancelled", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())

		todo := entity.Todo{
			ID:          0,
			Title:       "Get a coffee",
			Description: "Get an ice-latte in Starbucks",
			Completed:   false,
		}

		cancel()
		if err := s.CreateTodo(ctx, &todo); !errors.Is(err, context.Canceled) {
			t.Errorf("expected context.Canceled, got %v", err)
		}
	})

	t.Run("Task already exists", func(t *testing.T) {
		todo := entity.Todo{
			ID:          10,
			Title:       "Cook noodles",
			Description: "",
			Completed:   true,
		}

		_ = s.CreateTodo(ctx, &todo)

		if err := s.CreateTodo(ctx, &todo); !errors.Is(err, uc_errors.TodoAlreadyExistsError) {
			t.Errorf("expected ErrTodoAlreadyExists, got %v", err)
		}
	})
}

func TestStorage_GetTodo(t *testing.T) {
	s := NewDataStorage()
	ctx := context.Background()

	t.Run("Success", func(t *testing.T) {
		todo := entity.Todo{
			ID:          0,
			Title:       "Get a coffee",
			Description: "Get an ice-latte in Starbucks",
			Completed:   false,
		}

		_ = s.CreateTodo(ctx, &todo)

		got, err := s.GetTodo(ctx, todo.ID)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
			return
		}

		if todo.Title != got.Title || todo.Description != got.Description {
			t.Errorf("expected %v, got %v", todo, got)
		}
	})

	t.Run("Task not found", func(t *testing.T) {
		var fakeID = int64(1500)

		got, err := s.GetTodo(ctx, fakeID)
		if !errors.Is(err, uc_errors.TodoNotFoundError) {
			t.Errorf("expected ErrTodoNotFound, got %v", err)
			return
		}
		if got != nil {
			t.Errorf("expected nothing, got %v", got)
		}
	})
}

func TestStorage_GetList(t *testing.T) {
	s := NewDataStorage()
	t.Run("Success", func(t *testing.T) {
		ctx := context.Background()

		todo1 := entity.Todo{Title: "Get something"}
		todo2 := entity.Todo{Title: "Go somewhere"}
		todo3 := entity.Todo{Title: "Do somehow"}
		_ = s.CreateTodo(ctx, &todo1)
		_ = s.CreateTodo(ctx, &todo2)
		_ = s.CreateTodo(ctx, &todo3)

		var limit, offset = 2, 0
		list, err := s.GetTodoList(ctx, limit, offset)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		if len(list) != 2 {
			t.Errorf("expected 2 items, got %d", len(list))
		}

		if list[0].Title != todo1.Title || list[1].Title != todo2.Title {
			t.Errorf("expected %v, %v, but got %v, %v", todo1, todo2, list[0], list[1])
		}
	})

	t.Run("Empty list", func(t *testing.T) {
		emptyStorage := NewDataStorage()
		list, err := emptyStorage.GetTodoList(context.Background(), 10, 0)
		if err != nil {
			t.Errorf("expected no error, got %v", err)
		}
		if len(list) != 0 {
			t.Errorf("expected 0 items, got %d", len(list))
		}
	})

	t.Run("Context cancellation during Range", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		for i := 0; i < 100; i++ {
			_ = s.CreateTodo(ctx, &entity.Todo{Title: "Get something"})
		}

		cancel()

		_, err := s.GetTodoList(ctx, 100, 0)

		if !errors.Is(err, context.Canceled) {
			t.Errorf("expected context.Canceled error, got %v", err)
		}
	})
}

func TestStorage_UpdateTodo(t *testing.T) {
	s := NewDataStorage()
	ctx := context.Background()

	t.Run("Success", func(t *testing.T) {
		todo := entity.Todo{
			ID:          0,
			Title:       "Get a coffee",
			Description: "Get an ice-latte in Starbucks",
			Completed:   false,
		}

		_ = s.CreateTodo(ctx, &todo)

		todo.Title = "Get a cake"

		if err := s.UpdateTodo(ctx, &todo); err != nil {
			t.Errorf("expected no error, got %v", err)
		}
	})

	t.Run("Task not found", func(t *testing.T) {
		var fakeTodo = entity.Todo{
			ID:        100,
			Title:     "Get a pizza",
			Completed: true,
		}

		if err := s.UpdateTodo(ctx, &fakeTodo); !errors.Is(err, uc_errors.TodoNotFoundError) {
			t.Errorf("expected ErrTodoNotFound, got %v", err)
		}
	})
}

func TestStorage_DeleteTodo(t *testing.T) {
	s := NewDataStorage()
	ctx := context.Background()

	t.Run("Success", func(t *testing.T) {
		todo := entity.Todo{
			ID:          0,
			Title:       "Get a coffee",
			Description: "Get an ice-latte in Starbucks",
			Completed:   false,
		}

		_ = s.CreateTodo(ctx, &todo)

		if err := s.DeleteTodo(ctx, todo.ID); err != nil {
			t.Errorf("expected no error, got %v", err)
		}
	})

	t.Run("Task not found", func(t *testing.T) {
		var fakeID = int64(1000)

		if err := s.DeleteTodo(ctx, fakeID); !errors.Is(err, uc_errors.TodoNotFoundError) {
			t.Errorf("expected ErrTodoNotFound, got %v", err)
		}
	})
}
