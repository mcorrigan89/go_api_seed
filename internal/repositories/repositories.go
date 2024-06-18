package repositories

import (
	"errors"
	"sync"
	"time"

	"corrigan.io/go_api_seed/internal/repositories/models"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog"
)

const defaultTimeout = 10 * time.Second

var (
	ErrNotFound = errors.New("not found")
)

func StringToText(input string) (pgtype.Text, error) {
	text := pgtype.Text{}
	err := text.Scan(input)
	if err != nil {
		return text, err
	}
	return text, nil
}

type ServicesUtils struct {
	logger *zerolog.Logger
	wg     *sync.WaitGroup
}

type Repositories struct {
	utils          ServicesUtils
	UserRepository *UserRepository
}

func NewRepositories(db *pgxpool.Pool, logger *zerolog.Logger, wg *sync.WaitGroup) Repositories {
	queries := models.New(db)
	utils := ServicesUtils{
		logger: logger,
		wg:     wg,
	}

	userRepo := NewUserRepository(utils, db, queries)

	return Repositories{
		utils:          utils,
		UserRepository: userRepo,
	}
}
