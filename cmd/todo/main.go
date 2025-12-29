package main

import (
	"context"
	"errors"
	"log/slog"
	"net"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"todo-api/cmd/todo/config"
	adapterhttp "todo-api/internal/adapter/in/http"
	adapterstore "todo-api/internal/adapter/out/storage"
	"todo-api/internal/app/usecase"
)

const (
	shutdownTimeout = 10 * time.Second
)

func parseLogLevel(s string) slog.Level {
	switch strings.ToUpper(strings.TrimSpace(s)) {
	case "DEBUG":
		return slog.LevelDebug
	case "INFO":
		return slog.LevelInfo
	case "WARN", "WARNING":
		return slog.LevelWarn
	case "ERROR":
		return slog.LevelError
	default:
		return slog.LevelInfo
	}
}

func newLogger(level string) *slog.Logger {
	lv := parseLogLevel(level)
	return slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: lv,
	}))
}

func buildRouter(logger *slog.Logger) http.Handler {
	storage := adapterstore.NewDataStorage()

	createTodoUC := usecase.NewCreateTodoUC(storage)
	getTodoUC := usecase.NewGetTodoUC(storage)
	updateTodoUC := usecase.NewUpdateTodoUC(storage)
	deleteTodoUC := usecase.NewDeleteTodoUC(storage)
	getTodoListUC := usecase.NewGetTodoListUC(storage)

	todoHandler := adapterhttp.NewTodoHandler(
		logger,
		createTodoUC,
		getTodoUC,
		updateTodoUC,
		deleteTodoUC,
		getTodoListUC,
	)

	return adapterhttp.NewRouter(todoHandler).InitRoutes()
}

func run(ctx context.Context, cfg config.Config) error {
	logger := newLogger(cfg.LogLevel)
	router := buildRouter(logger)

	srv := &http.Server{
		Addr:    cfg.HTTPAddress,
		Handler: router,

		BaseContext: func(_ net.Listener) context.Context { // any чтобы не тянуть net.Listener в импорты
			return ctx
		},
	}

	errCh := make(chan error, 1)

	go func() {
		logger.Info("starting server", slog.String("address", cfg.HTTPAddress))
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			errCh <- err
			return
		}
		errCh <- nil
	}()

	select {
	case <-ctx.Done():
		logger.Info("shutdown signal received")
	case err := <-errCh:
		if err != nil {
			logger.Error("server failed", slog.Any("err", err))
			return err
		}
		logger.Info("server stopped")
		return nil
	}

	shutdownCtx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		logger.Error("graceful shutdown failed", slog.Any("err", err))
		_ = srv.Close() // fallback
		return err
	}

	logger.Info("server exited properly")
	return nil
}

func main() {
	cfg := config.Load()

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	if err := run(ctx, *cfg); err != nil {
		os.Exit(1)
	}
}
