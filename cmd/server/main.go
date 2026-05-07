// Command server runs the accident reconstruction API.
package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/baditaflorin/accident-reconstructor/internal/config"
	"github.com/baditaflorin/accident-reconstructor/internal/httpapi"
	"github.com/baditaflorin/accident-reconstructor/internal/jobs"
	"github.com/baditaflorin/accident-reconstructor/internal/reconstruction"
	"github.com/baditaflorin/accident-reconstructor/internal/utils"
)

var (
	version = "0.1.0"
	commit  = "dev"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
	slog.SetDefault(logger)

	cfg, err := config.Load()
	if err != nil {
		utils.HandleErrorOrLogWithMessages(err, "configuration failed", "")
		os.Exit(1)
	}
	if cfg.Version == "" || cfg.Version == "0.1.0" {
		cfg.Version = version
	}
	if cfg.Commit == "" || cfg.Commit == "dev" {
		cfg.Commit = commit
	}

	store := jobs.NewStore()
	processor := reconstruction.NewProcessor(cfg)
	server := &http.Server{
		Addr:              cfg.Addr,
		Handler:           httpapi.NewRouter(cfg, store, processor),
		ReadHeaderTimeout: 15 * time.Second,
	}

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	go func() {
		slog.Info("server_started", "addr", cfg.Addr, "version", cfg.Version, "commit", cfg.Commit)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			utils.HandleErrorOrLogWithMessages(err, "server failed", "")
			os.Exit(1)
		}
	}()

	<-ctx.Done()
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	utils.HandleErrorOrLogWithMessages(server.Shutdown(shutdownCtx), "graceful shutdown failed", "server_stopped")
}
