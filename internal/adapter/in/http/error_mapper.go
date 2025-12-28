package http

import (
	"errors"
	"net/http"
	"todo-api/internal/app/uc_errors"
)

func HttpError(err error) (int, string, error) {
	if w, ok := err.(*uc_errors.WrappedError); ok {
		switch w.Public {
		case uc_errors.CreateTodoError,
			uc_errors.GetTodoError,
			uc_errors.GetTodoListError,
			uc_errors.UpdateTodoError,
			uc_errors.DeleteTodoError:
			return http.StatusInternalServerError, w.Public.Error(), w.Reason
		default:
			return http.StatusInternalServerError, "internal error", w.Reason
		}
	}

	switch {
	case errors.Is(err, uc_errors.TodoNotFoundError):
		return http.StatusNotFound, err.Error(), nil
	case errors.Is(err, uc_errors.TodoAlreadyExistsError),
		errors.Is(err, uc_errors.EmptyTitleError),
		errors.Is(err, uc_errors.InvalidTodoIDError),
		errors.Is(err, uc_errors.InvalidLimitError),
		errors.Is(err, uc_errors.InvalidOffsetError):
		return http.StatusBadRequest, err.Error(), nil
	}

	return http.StatusInternalServerError, "internal error", err
}
