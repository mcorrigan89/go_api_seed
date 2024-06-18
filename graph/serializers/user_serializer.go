package serializers

import (
	"corrigan.io/go_api_seed/graph/dto"
	"corrigan.io/go_api_seed/internal/entities"
)

func UserSerializer(entity *entities.User) dto.User {

	response := dto.User{
		ID:         entity.ID.String(),
		GivenName:  entity.GivenName,
		FamilyName: entity.FamilyName,
		Email:      entity.Email,
		AvatarURL:  entity.AvatarUrl,
	}

	return response
}

func UserErrorSerializer(err error) (dto.UserResult, error) {
	switch err {
	case entities.ErrUserNotFound:
		return &dto.UserNotFound{
			Message: err.Error(),
			Code:    "USER_NOT_FOUND",
		}, nil
	default:
		return nil, err
	}
}

func CreateUserErrorSerializer(err error) (dto.CreateUserResult, error) {
	switch err {
	case entities.ErrDuplicateEmail:
		return &dto.EmailUnavailable{
			Message: err.Error(),
			Code:    "DUPLICATE_EMAIL",
		}, nil
	default:
		return nil, err
	}
}
