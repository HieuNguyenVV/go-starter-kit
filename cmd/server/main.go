package main

import (
	"fmt"
	"go-starter-kit/internal/log"
	"go-starter-kit/internal/pkg/database"
	"go-starter-kit/internal/server"
	"go-starter-kit/internal/server/config"
)

func main() {
	conf, err := config.NewConfig()
	if err != nil {
		panic("init config failed: " + err.Error())
	}
	logger, err := log.NewLogger(conf)
	if err != nil {
		panic("init logger failed: " + err.Error())
	}

	logger.Infof("Git test")
	postgres, err := database.NewPostgres(conf, logger)
	if err != nil {
		logger.Fatal("database:NewPostgres: init failed: %s", err)
	}

	httpClient := server.NewHTTPServer(logger, conf)

	srv := server.NewServer(conf, logger, httpClient, postgres)

	go func() {
		defer func() {
			if r := recover(); r != nil {
				logger.Info(fmt.Sprintf("Recovered from panic: %v", r))
			}
		}()
	}()
	srv.Run()
}
