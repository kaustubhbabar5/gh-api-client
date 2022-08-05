package main

import (
	"log"

	"github.com/kaustubhbabar5/gh-api-client/pkg/config"
	"github.com/kaustubhbabar5/gh-api-client/pkg/logger"
	"github.com/kaustubhbabar5/gh-api-client/server"
	"go.uber.org/zap"
)

// @title           GitHub Users API
// @version         1.0

// @contact.name   Kaustubh Babar
// @contact.email  kaustubhbabar5@gmail.com
func main() {
	logger, err := logger.NewProduction()
	if err != nil {
		log.Fatalf("failed to create logger %s", err.Error())
	}
	defer func() {
		e := logger.Sync()
		if e != nil {
			log.Fatalf("failed to flush log buffer :%s", err.Error())
		}
	}()

	config, err := config.Load(".", logger)
	if err != nil {
		logger.Error("failed to load config", zap.Error(err))
	}

	httpServer, err := server.New(config, logger)
	if err != nil {
		logger.Error("failed to create http server", zap.Error(err))
	}

	err = httpServer.ListenAndServe()
	if err != nil {
		logger.Error("failed to start http server", zap.Error(err))
	}

	// TODO: graceful shutdown
}
