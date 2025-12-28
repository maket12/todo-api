package main

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
	"todo-api/cmd/todo/config"
	adapterhttp "todo-api/internal/adapter/in/http"
	adapterstore "todo-api/internal/adapter/out/storage"
	"todo-api/internal/app/usecase"
)

func parseLogLevel(level string) slog.Level {
	switch level {
	case "DEBUG":
		return slog.LevelDebug
	case "INFO":
		return slog.LevelInfo
	case "WARN":
		return slog.LevelWarn
	case "ERROR":
		return slog.LevelError
	default:
		return slog.LevelInfo
	}
}

func main() {
	// ======================
	// 1. Load config
	// ======================
	cfg := config.Load()

	// ======================
	// 2. Setup logger
	// ======================
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: parseLogLevel(cfg.LogLevel),
	}))

	// ======================
	// 3. Storage
	// ======================
	storage := adapterstore.NewDataStorage()

	// ======================
	// 4. Usecases
	// ======================
	createTodoUC := usecase.NewCreateTodoUC(storage)
	getTodoUC := usecase.NewGetTodoUC(storage)
	updateTodoUC := usecase.NewUpdateTodoUC(storage)
	deleteTodoUC := usecase.NewDeleteTodoUC(storage)
	getTodoListUC := usecase.NewGetTodoListUC(storage)

	// ======================
	// 5. Handlers (REST)
	// ======================
	todoHandler := adapterhttp.NewTodoHandler(
		logger, createTodoUC, getTodoUC,
		updateTodoUC, deleteTodoUC, getTodoListUC,
	)

	// ======================
	// 6. Router
	// ======================
	router := adapterhttp.NewRouter(todoHandler).InitRoutes()

	// ======================
	// 7. Run HTTP server
	// ======================
	srv := &http.Server{
		Addr:    cfg.HTTPAddress,
		Handler: router,
	}

	go func() {
		logger.Info("starting server", slog.String("address", cfg.HTTPAddress))
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			logger.Error("server failed to start", slog.Any("err", err))
			os.Exit(1)
		}
	}()

	// ======================
	// 8. Graceful Shutdown
	// ======================
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	<-quit
	logger.Info("shutdown signal received")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		logger.Error("server forced to shutdown", slog.Any("err", err))
	}

	logger.Info("server exited properly")
}
