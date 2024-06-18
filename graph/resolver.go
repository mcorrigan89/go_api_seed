package graph

import (
	"corrigan.io/go_api_seed/internal/services"
	"github.com/rs/zerolog"
)

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct {
	services *services.Services
	logger   *zerolog.Logger
}

func NewResolver(services *services.Services, logger *zerolog.Logger) *Resolver {
	return &Resolver{
		services: services,
		logger:   logger,
	}
}
