package serializers

import (
	"time"

	"corrigan.io/go_api_seed/graph/dto"
	"corrigan.io/go_api_seed/internal/entities"
)

func UserSessionSerializer(entity *entities.UserSession) dto.UserSession {

	response := dto.UserSession{
		Token:     entity.Token,
		ExpiresAt: entity.ExpiresAt.Format(time.RFC1123Z),
	}

	return response
}

func UserSessionErrorSerializer(err error) (dto.UserSessionResult, error) {
	switch err {
	case entities.ErrUserNotFound:
		return &dto.UserNotFound{
			Message: err.Error(),
			Code:    "USER_NOT_FOUND",
		}, nil
	case entities.ErrInvalidCredentials:
		return &dto.InvalidCredentials{
			Message: err.Error(),
			Code:    "INVALID_CREDENTIALS",
		}, nil
	default:
		return nil, err
	}
}
