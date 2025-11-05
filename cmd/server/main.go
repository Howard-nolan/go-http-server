package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/joeynolan/go-http-server/internal/http/handlers"
	"github.com/joeynolan/go-http-server/internal/platform/config"
	ilog "github.com/joeynolan/go-http-server/internal/platform/log"
)

func main() {
	// config + logger
	cfg := config.Load()
	logger := ilog.New()

	mux := http.NewServeMux()
	handlers.Register(mux)

	srv := &http.Server{
		Addr:              fmt.Sprintf(":%d", cfg.Port),
		Handler:           mux,
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
