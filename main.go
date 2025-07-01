package main

import (
	"context"
	"lexia/internal/logger"
	"lexia/internal/modules"
	"lexia/internal/shared"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
	_ "time/tzdata"
)

func main() {
	shared.LoadEnv()

	envVars, err := shared.ParseEnv()
	if err != nil {
		panic(err)
	}

	db, err := shared.InitializeDatabase()
	if err != nil {
		panic(err)
	}

	defer db.Close()

	resouceConfig := &shared.ResourceConfig{
		DB: db,
	}

	apiCfg := shared.ApiConfig{
		ResourceConfig: resouceConfig,
	}

	server, err := modules.CreateWebserver(&apiCfg)
	if err != nil {
		panic(err)
	}

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt, syscall.SIGTERM)

	httpServer := &http.Server{
		Addr:    ":" + envVars.Port,
		Handler: server,
	}

	go func() {
		logger.Info("Starting HTTP server on :" + envVars.Port)
		if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatal("Failed to start HTTP server: ", err)
		}
	}()

	<-signalChan
	logger.Info("Received shutdown signal, shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := httpServer.Shutdown(ctx); err != nil {
		logger.Fatal("Server forced to shutdown: ", err)
	}

	logger.Info("Server exited gracefully")
}
