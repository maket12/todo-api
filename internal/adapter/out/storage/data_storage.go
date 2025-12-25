package storage

import (
	"context"
	"sort"
	"sync"
	"sync/atomic"
	"todo-api/internal/app/uc_errors"
	"todo-api/internal/domain/entity"
)

type DataStorage struct {
	data   sync.Map
	prevID int64
}

func NewDataStorage() *DataStorage {
	return &DataStorage{}
}

func (s *DataStorage) CreateTodo(ctx context.Context, todo *entity.Todo) error {
	if err := ctx.Err(); err != nil {
		return err
	}

	if todo.ID == 0 {
		todo.ID = atomic.AddInt64(&s.prevID, 1)
	} else {
		if _, exists := s.data.Load(todo.ID); exists {
			return uc_errors.TodoAlreadyExistsError
		}
	}

	s.data.Store(todo.ID, *todo)
	return nil
}

func (s *DataStorage) GetTodo(ctx context.Context, id int64) (*entity.Todo, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}

	raw, ok := s.data.Load(id)
	if !ok {
		return nil, uc_errors.TodoNotFoundError
	}

	todo := raw.(entity.Todo)

	return &todo, nil
}

func (s *DataStorage) GetTodoList(ctx context.Context, limit, offset int) ([]*entity.Todo, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}

	var (
		todos    []entity.Todo
		rangeErr error
	)

	s.data.Range(func(key, value any) bool {
		select {
		case <-ctx.Done():
			rangeErr = ctx.Err()
			return false
		default:
			todos = append(todos, value.(entity.Todo))
			return true
		}
	})

	if rangeErr != nil {
		return nil, rangeErr
	}

	sort.Slice(todos, func(i, j int) bool {
		return todos[i].ID < todos[j].ID
	})

	start := offset
	if start > len(todos) {
		return []*entity.Todo{}, nil
	}

	end := offset + limit
	if end > len(todos) || limit == 0 {
		end = len(todos)
	}

	result := make([]*entity.Todo, 0, end-start)
	for i := start; i < end; i++ {
		result = append(result, &todos[i])
	}

	return result, nil
}

func (s *DataStorage) UpdateTodo(ctx context.Context, todo *entity.Todo) error {
	if err := ctx.Err(); err != nil {
		return err
	}

	if _, ok := s.data.Load(todo.ID); !ok {
		return uc_errors.TodoNotFoundError
	}

	s.data.Store(todo.ID, *todo)
	return nil
}

func (s *DataStorage) DeleteTodo(ctx context.Context, id int64) error {
	if err := ctx.Err(); err != nil {
		return err
	}

	if _, ok := s.data.LoadAndDelete(id); !ok {
		return uc_errors.TodoNotFoundError
	}

	return nil
}
