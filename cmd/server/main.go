package main

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/adapters/humachi"
	"github.com/danielgtaylor/huma/v2/humacli"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/serroba/features/internal/flags"
	"github.com/serroba/features/internal/handler"
)

type Options struct {
	Port int `default:"8080" help:"Port to listen on"`
}

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	cli := humacli.New(func(hooks humacli.Hooks, options *Options) {
		router := chi.NewRouter()
		router.Use(middleware.Logger)
		router.Use(middleware.Recoverer)

		api := humachi.New(router, huma.DefaultConfig("Feature Flags API", "1.0.0"))

		repo := flags.NewMemoryRepository()
		service := flags.NewService(repo)
		h := handler.New(service)
		h.Register(api)

		var server *http.Server

		hooks.OnStart(func() {
			server = &http.Server{
				Addr:              ":" + strconv.Itoa(options.Port),
				Handler:           router,
				ReadHeaderTimeout: 10 * time.Second,
			}

			logger.Info("server starting", slog.Int("port", options.Port))

			if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
				logger.Error("server failed", slog.Any("error", err))
				os.Exit(1)
			}
		})

		hooks.OnStop(func() {
			logger.Info("shutting down")

			ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
			defer cancel()

			if server != nil {
				if err := server.Shutdown(ctx); err != nil {
					logger.Error("server shutdown error", slog.Any("error", err))
				}
			}

			logger.Info("shutdown complete")
		})
	})

	cli.Run()
}
