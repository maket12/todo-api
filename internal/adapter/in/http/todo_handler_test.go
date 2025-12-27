package http

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
	"todo-api/internal/adapter/out/storage"
	"todo-api/internal/app/dto"
	"todo-api/internal/app/usecase"
)

func TestTH_Create(t *testing.T) {
	store := storage.NewDataStorage()
	uc := usecase.NewCreateTodoUC(store)
	testLogger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	handler := NewTodoHandler(
		testLogger,
		uc,
		nil,
		nil,
		nil,
		nil,
	)

	router := NewRouter(handler)
	mux := router.InitRoutes()

	t.Run("Success", func(t *testing.T) {
		reqBody := `{
						"id": 0,
						"title": "Clean room №1",
						"description": "use water",
						"completed": true
					}`
		request := httptest.NewRequest("POST", "/todos", strings.NewReader(reqBody))
		recorder := httptest.NewRecorder()
		mux.ServeHTTP(recorder, request)

		if recorder.Code != http.StatusCreated {
			t.Errorf("expected status 201, got %v", recorder.Code)
		}

		var response dto.CreateTodoResponse
		_ = json.NewDecoder(recorder.Body).Decode(&response)
		if response.ID == 0 {
			t.Error("expected non-zero ID in response")
		}
	})

}
