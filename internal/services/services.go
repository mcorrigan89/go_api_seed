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
	config *config.Config
}

type Services struct {
	utils        ServicesUtils
	UserService  *UserService
	OAuthService *OAuthService
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
		config: cfg,
	}

	userService := NewUserService(utils, repositories.UserRepository)
	oAuthService := NewOAuthService(utils, userService, repositories.UserRepository)

	return Services{
		utils:        utils,
		UserService:  userService,
		OAuthService: oAuthService,
	}
}
