package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	apphttp "github.com/joeynolan/go-http-server/internal/http"
	handlers "github.com/joeynolan/go-http-server/internal/http/handlers"
	"github.com/joeynolan/go-http-server/internal/platform/config"
	ilog "github.com/joeynolan/go-http-server/internal/platform/log"

	db "github.com/joeynolan/go-http-server/internal/db"
)

func main() {
	sqlDB, err := db.OpenAndMigrate("./data/app.db")
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to open database: %v\n", err)
		os.Exit(1)
	}
	defer sqlDB.Close()

	// config + logger
	cfg := config.Load()
	logger := ilog.New()

	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	h := handlers.NewHandler(sqlDB)

	apphttp.Register(r, h)

	srv := &http.Server{
		Addr:              fmt.Sprintf(":%d", cfg.Port),
		Handler:           r,
		ReadHeaderTimeout: 5 * time.Second,
	}

	logger.Infof("starting server on port %d", cfg.Port)

	// Start server
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Errorf("listen: %v", err)
			os.Exit(1)
		}
	}()

	// Graceful shutdown on SIGINT/SIGTERM
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	<-stop

	logger.Infof("shutting down server...")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		logger.Errorf("server shutdown failed: %v", err)
	}
	logger.Infof("bye 👋")
}
