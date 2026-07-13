package main

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"project-manager-go/internal/auth"
	"project-manager-go/internal/config"
	"project-manager-go/internal/httpapi"
	"project-manager-go/internal/service"
	"project-manager-go/internal/store/postgres"
)

func main() {
	if err := run(); err != nil {
		slog.Error("server stopped", "error", err)
		os.Exit(1)
	}
}

func run() error {
	cfg, err := config.Load()
	if err != nil {
		return err
	}
	if cfg.DatabaseURL == "" {
		return errors.New("DATABASE_URL is required")
	}

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	db, err := postgres.Open(ctx, cfg.DatabaseURL)
	if err != nil {
		return err
	}
	defer db.Close()

	tokenManager := auth.TokenManager{
		AccessSecret:  []byte(cfg.JWTSecret),
		RefreshSecret: []byte(cfg.JWTRefreshSecret),
		AccessTTL:     cfg.JWTExpiresIn,
		RefreshTTL:    cfg.JWTRefreshExpiresIn,
	}
	svc := service.New(postgres.New(db), tokenManager)
	router := httpapi.NewRouter(cfg, svc)

	server := &http.Server{
		Addr:              fmt.Sprintf(":%d", cfg.Port),
		Handler:           router,
		ReadHeaderTimeout: 5 * time.Second,
	}

	go func() {
		slog.Info("server started", "port", cfg.Port, "env", cfg.NodeEnv)
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			slog.Error("listen failed", "error", err)
			stop()
		}
	}()

	<-ctx.Done()
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	return server.Shutdown(shutdownCtx)
}
