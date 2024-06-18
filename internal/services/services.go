package services

import (
	"sync"

	"corrigan.io/go_api_seed/internal/config"
	"corrigan.io/go_api_seed/internal/repositories"
	"github.com/rs/zerolog"
)

type ServicesUtils struct {
	logger *zerolog.Logger
	wg     *sync.WaitGroup
}

type Services struct {
	utils       ServicesUtils
	UserService *UserService
}

func (utils *ServicesUtils) background(fn func()) {
	utils.wg.Add(1)

	go func() {
		defer utils.wg.Done()

		defer func() {
			if err := recover(); err != nil {
				utils.logger.Error().Msg("panic in background function")
			}
		}()

		fn()
	}()
}

func NewServices(repositories *repositories.Repositories, cfg *config.Config, logger *zerolog.Logger, wg *sync.WaitGroup) Services {
	utils := ServicesUtils{
		logger: logger,
		wg:     wg,
	}

	userService := NewUserService(utils, repositories.UserRepository)

	return Services{
		utils:       utils,
		UserService: userService,
	}
}
