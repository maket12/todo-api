package http

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"todo-api/internal/app/dto"
	"todo-api/internal/app/usecase"
)

type TodoHandler struct {
	log           *slog.Logger
	createTodoUC  *usecase.CreateTodoUC
	getTodoUC     *usecase.GetTodoUC
	getTodoListUC *usecase.GetTodoListUC
	updateTodoUC  *usecase.UpdateTodoUC
	deleteTodoUC  *usecase.DeleteTodoUC
}

func NewTodoHandler(
	log *slog.Logger,
	createTodoUC *usecase.CreateTodoUC,
	getTodoUC *usecase.GetTodoUC,
	getTodoListUC *usecase.GetTodoListUC,
	updateTodoUC *usecase.UpdateTodoUC,
	deleteTodoUC *usecase.DeleteTodoUC,
) *TodoHandler {
	return &TodoHandler{
		log:           log,
		createTodoUC:  createTodoUC,
		getTodoUC:     getTodoUC,
		getTodoListUC: getTodoListUC,
		updateTodoUC:  updateTodoUC,
		deleteTodoUC:  deleteTodoUC,
	}
}

func (h *TodoHandler) CreateTodo(w http.ResponseWriter, r *http.Request) {
	var input dto.CreateTodo
	if err := json.NewDecoder(r.Body).Decode(&input.Todo); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	response, err := h.createTodoUC.Execute(r.Context(), input)
	if err != nil {
		status, msg, internalErr := HttpError(err)
		h.log.ErrorContext(r.Context(), "failed to create todo",
			slog.Int("status", status),
			slog.String("public_msg", msg),
			slog.Any("cause", internalErr),
		)
		http.Error(w, msg, status)
		return
	}

	h.log.InfoContext(r.Context(), "created todo",
		slog.Int("id", int(response.ID)),
	)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(response)
}

func (h *TodoHandler) GetTodo(w http.ResponseWriter, r *http.Request) {
	var input dto.GetTodo
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	response, err := h.getTodoUC.Execute(r.Context(), input)
	if err != nil {
		status, msg, internalErr := HttpError(err)
		h.log.ErrorContext(r.Context(), "failed to get todo",
			slog.Int("status", status),
			slog.String("public_msg", msg),
			slog.Any("cause", internalErr),
		)
		http.Error(w, msg, status)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(response)
}

func (h *TodoHandler) GetTodoList(w http.ResponseWriter, r *http.Request) {
	var input dto.GetTodoList
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	response, err := h.getTodoListUC.Execute(r.Context(), input)
	if err != nil {
		status, msg, internalErr := HttpError(err)
		h.log.ErrorContext(r.Context(), "failed to get todo list",
			slog.Int("status", status),
			slog.String("public_msg", msg),
			slog.Any("cause", internalErr),
		)
		http.Error(w, msg, status)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(response)
}

func (h *TodoHandler) UpdateTodo(w http.ResponseWriter, r *http.Request) {
	var input dto.UpdateTodo
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	response, err := h.updateTodoUC.Execute(r.Context(), input)
	if err != nil {
		status, msg, internalErr := HttpError(err)
		h.log.ErrorContext(r.Context(), "failed to update todo",
			slog.Int("status", status),
			slog.String("public_msg", msg),
			slog.Any("cause", internalErr),
		)
		http.Error(w, msg, status)
		return
	}

	h.log.InfoContext(r.Context(), "updated todo",
		slog.Int("id", int(response.ID)),
	)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(response)
}

func (h *TodoHandler) DeleteTodo(w http.ResponseWriter, r *http.Request) {
	var input dto.DeleteTodo
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	response, err := h.deleteTodoUC.Execute(r.Context(), input)
	if err != nil {
		status, msg, internalErr := HttpError(err)
		h.log.ErrorContext(r.Context(), "failed to delete todo",
			slog.Int("status", status),
			slog.String("public_msg", msg),
			slog.Any("cause", internalErr),
		)
		http.Error(w, msg, status)
		return
	}

	h.log.InfoContext(r.Context(), "deleted todo",
		slog.Int("id", int(response.ID)),
	)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(response)
}
