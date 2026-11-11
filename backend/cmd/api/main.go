package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/xsolla-school/office-explorer/backend/internal/adapter/handler"
	"github.com/xsolla-school/office-explorer/backend/internal/adapter/storage/postgres"
	"github.com/xsolla-school/office-explorer/backend/internal/config"
	"github.com/xsolla-school/office-explorer/backend/internal/usecase"
)

func main() {
	cfg := config.Load()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	pool, err := postgres.NewPool(ctx, cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("failed to connect to PostgreSQL: %v", err)
	}
	defer pool.Close()

	zoneRepo := postgres.NewZoneRepo(pool)
	questRepo := postgres.NewQuestRepo(pool)
	progressRepo := postgres.NewProgressRepo(pool)

	explorerUC := usecase.NewExplorerUseCase(zoneRepo, progressRepo)
	questUC := usecase.NewQuestUseCase(zoneRepo, questRepo, progressRepo)

	explorerHandler := handler.NewExplorerHandler(explorerUC)
	questHandler := handler.NewQuestHandler(questUC)
	router := handler.NewRouter(explorerHandler, questHandler, cfg.AllowedOrigin, []byte(cfg.LMSJWTSecret))
	if cfg.LMSJWTSecret == "" {
		log.Println("LMS sign-in is off (LMS_JWT_SECRET not set): every visitor is a guest")
	} else {
		log.Println("LMS sign-in is on: an LMS access token in Authorization identifies the student")
	}

	server := &http.Server{
		Addr:         ":" + cfg.HTTPPort,
		Handler:      router,
		ReadTimeout:  cfg.ReadTimeout,
		WriteTimeout: cfg.WriteTimeout,
	}

	go func() {
		log.Printf("Xsolla Office Explorer API listening on http://localhost:%s", cfg.HTTPPort)
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("server error: %v", err)
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
	<-stop

	log.Println("shutting down server gracefully...")
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), cfg.ShutdownTimeout)
	defer shutdownCancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		log.Printf("graceful shutdown failed: %v", err)
	}
	log.Println("server stopped")
}
