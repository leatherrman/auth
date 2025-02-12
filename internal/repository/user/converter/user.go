package converter

import (
	"github.com/katyafirstova/auth_service/internal/model"
	modelRepo "github.com/katyafirstova/auth_service/internal/repository/user/model"
)

func ToUserFromRepo(user modelRepo.User) model.User {
	return model.User{
		UUID:      user.UUID,
		Name:      user.Name,
		Email:     user.Email,
		Role:      model.Role(user.Role),
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}
}
