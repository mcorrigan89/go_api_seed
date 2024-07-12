package main

import (
	"os"
	"sync"

	"corrigan.io/go_api_seed/internal/config"
	"corrigan.io/go_api_seed/internal/repositories"
	"corrigan.io/go_api_seed/internal/services"
	"github.com/rs/zerolog"
)

type application struct {
	config   config.Config
	wg       sync.WaitGroup
	logger   *zerolog.Logger
	services *services.Services
}

func main() {

	cfg := config.Config{}

	config.LoadConfig(&cfg)

	logger := getLogger(&cfg)

	db, err := openDBPool(cfg, &logger)
	if err != nil {
		logger.Err(err)
		os.Exit(1)
	}
	defer db.Close()

	wg := sync.WaitGroup{}

	repositories := repositories.NewRepositories(db, &logger, &wg)
	services := services.NewServices(&repositories, &cfg, &logger, &wg)

	app := &application{
		config:   cfg,
		logger:   &logger,
		services: &services,
	}

	err = app.serve()
	if err != nil {
		logger.Err(err)
		os.Exit(1)
	}
}
