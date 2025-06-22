package main

import (
	"lexia/internal/logger"
	"lexia/internal/modules"
	"lexia/internal/shared"
	"os"
	"os/signal"
	"syscall"
	_ "time/tzdata"
)

func main() {
	shared.LoadEnv()

	db, err := shared.InitializeDatabase()
	if err != nil {
		panic(err)
	}

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

	if err := server.Run(); err != nil {
		logger.Fatal("Failed to start HTTP server: ", err)
		return
	}

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt, syscall.SIGTERM)

	<-signalChan
	logger.Info("Received shutdown signal, shutting down server...")
}
