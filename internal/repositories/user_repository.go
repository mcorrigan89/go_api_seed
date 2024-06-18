package repositories

import (
	"context"
	"errors"
	"time"

	"corrigan.io/go_api_seed/internal/entities"
	"corrigan.io/go_api_seed/internal/repositories/models"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
	"golang.org/x/crypto/bcrypt"
)

type UserRepository struct {
	utils   ServicesUtils
	DB      *pgxpool.Pool
	queries *models.Queries
}

func NewUserRepository(utils ServicesUtils, db *pgxpool.Pool, queries *models.Queries) *UserRepository {
	return &UserRepository{
		utils:   utils,
		DB:      db,
		queries: queries,
	}
}

func (repo *UserRepository) GetUserByID(ctx context.Context, id uuid.UUID) (*entities.User, error) {
	ctx, cancel := context.WithTimeout(ctx, defaultTimeout)
	defer cancel()
	row, err := repo.queries.GetUserByID(ctx, id)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, entities.ErrUserNotFound
		} else {
			repo.utils.logger.Err(err).Ctx(ctx).Msg("Get user by ID")
			return nil, err
		}
	}

	entity := repo.userModelToEntity(row.User, row.UserAuth)

	return entity, nil
}

func (repo *UserRepository) GetUserByEmail(ctx context.Context, email string) (*entities.User, error) {
	ctx, cancel := context.WithTimeout(ctx, defaultTimeout)
	defer cancel()
	row, err := repo.queries.GetUserByEmail(ctx, email)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, entities.ErrUserNotFound
		} else {
			repo.utils.logger.Err(err).Ctx(ctx).Msg("Get user by email")
			return nil, err
		}
	}

	entity := repo.userModelToEntity(row.User, row.UserAuth)

	return entity, nil
}

func (repo *UserRepository) GetUserBySessionToken(ctx context.Context, token string) (*entities.User, *entities.UserSession, error) {
	ctx, cancel := context.WithTimeout(ctx, defaultTimeout)
	defer cancel()
	row, err := repo.queries.GetUserBySessionToken(ctx, token)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil, entities.ErrUserNotFound
		} else {
			repo.utils.logger.Err(err).Ctx(ctx).Msg("Get user by session token")
			return nil, nil, err
		}
	}

	userEntity := repo.userModelToEntity(row.User, row.UserAuth)
	sessionEntity := entities.NewUserSession(entities.NewUserSessionArgs{
		ID:            row.UserSession.ID,
		UserID:        row.UserSession.UserID,
		Token:         row.UserSession.Token,
		ExpiresAt:     row.UserSession.ExpiresAt.Time,
		ExpiredByUser: row.UserSession.UserExpired,
	})

	return userEntity, sessionEntity, nil
}

type CreateUserArgs struct {
	GivenName  *string
	FamilyName *string
	Email      string
	Password   string
}

func (repo *UserRepository) CreateUserPassword(ctx context.Context, args CreateUserArgs) (*entities.User, error) {

	repo.utils.logger.Info().Ctx(ctx).Interface("args", args).Msg("Creating user")
	ctx, cancel := context.WithTimeout(ctx, defaultTimeout)
	defer cancel()

	tx, err := repo.DB.Begin(ctx)
	if err != nil {
		repo.utils.logger.Err(err).Ctx(ctx).Msg("Begin transaction")
		return nil, err
	}
	defer tx.Rollback(ctx)

	qtx := repo.queries.WithTx(tx)

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(args.Password), 12)
	if err != nil {
		repo.utils.logger.Err(err).Ctx(ctx).Msg("Generate from password")
		return nil, err
	}
	hashedPasswordString := string(hashedPassword)

	userRow, err := qtx.CreateUser(ctx, models.CreateUserParams{
		GivenName:  args.GivenName,
		FamilyName: args.FamilyName,
		Email:      args.Email,
	})
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			if pgErr.Code == "23505" && pgErr.ConstraintName == "users_email_key" {
				repo.utils.logger.Err(err).Ctx(ctx).Msg("Duplicate email")
				return nil, entities.ErrDuplicateEmail
			}
		}
		repo.utils.logger.Err(err).Ctx(ctx).Msg("Create user")
		return nil, err
	}

	userAuthRow, err := qtx.CreateUserAuth(ctx, models.CreateUserAuthParams{
		UserID:   userRow.ID,
		Value:    hashedPasswordString,
		Provider: entities.ProviderPassword,
	})
	if err != nil {
		repo.utils.logger.Err(err).Ctx(ctx).Msg("Create user auth")
		return nil, err
	}

	err = tx.Commit(ctx)
	if err != nil {
		repo.utils.logger.Err(err).Ctx(ctx).Msg("Commit transaction")
		return nil, err
	}

	entity := repo.userModelToEntity(userRow, userAuthRow)

	return entity, nil
}

type CreateUserSessionArgs struct {
	UserID    uuid.UUID
	Token     string
	ExpiresAt time.Time
}

func (repo *UserRepository) CreateUserSession(ctx context.Context, args CreateUserSessionArgs) (*entities.UserSession, error) {
	ctx, cancel := context.WithTimeout(ctx, defaultTimeout)
	defer cancel()

	expires := pgtype.Timestamptz{}
	err := expires.Scan(args.ExpiresAt)
	if err != nil {
		repo.utils.logger.Err(err).Ctx(ctx).Msg("Scan expires at")
		return nil, err
	}

	row, err := repo.queries.CreateUserSession(ctx, models.CreateUserSessionParams{
		UserID:    args.UserID,
		Token:     args.Token,
		ExpiresAt: expires,
	})
	if err != nil {
		repo.utils.logger.Err(err).Ctx(ctx).Msg("Create user session")
		return nil, err
	}

	entity := entities.NewUserSession(entities.NewUserSessionArgs{
		ID:            row.ID,
		UserID:        row.UserID,
		Token:         row.Token,
		ExpiresAt:     row.ExpiresAt.Time,
		ExpiredByUser: row.UserExpired,
	})

	return entity, nil
}

func (repo *UserRepository) ExpireUserSession(ctx context.Context, sessionID uuid.UUID) error {
	ctx, cancel := context.WithTimeout(ctx, defaultTimeout)
	defer cancel()
	err := repo.queries.ExpireUserSession(ctx, sessionID)
	if err != nil {
		repo.utils.logger.Err(err).Ctx(ctx).Msg("Expire user session")
		return err
	}
	return nil
}

func (repo *UserRepository) userModelToEntity(userModel models.User, userAuthModel models.UserAuth) *entities.User {
	entity := entities.NewUserEntity(entities.NewUserEntityArgs{
		ID:         userModel.ID,
		GivenName:  userModel.GivenName,
		FamilyName: userModel.FamilyName,
		Email:      userModel.Email,
		UserAuth: &entities.UserAuth{
			Value:    userAuthModel.Value,
			Provider: userAuthModel.Provider,
		},
	})

	return entity
}
