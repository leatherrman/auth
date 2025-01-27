package converter

import (
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/katyafirstova/auth_service/internal/model"
	desc "github.com/katyafirstova/auth_service/pkg/user_v1"
)

func CreateUserToServiceFromApi(user *desc.CreateRequest) model.CreateUser {
	return model.CreateUser{
		Name:            user.Name,
		Email:           user.Email,
		Password:        user.Password,
		PasswordConfirm: user.PasswordConfirm,
		Role:            model.Role(user.Role),
	}
}

func UpdateUserToServiceFromApi(user *desc.UpdateRequest) model.UpdateUser {
	var res model.UpdateUser
	if user.Name != nil {
		res.Name = &user.Name.Value
	}

	if user.Email != nil {
		res.Email = &user.Email.Value
	}

	if user.Role != desc.Role_UNKNOWN {
		res.Role = model.Role(user.Role)

	}

	return res
}

func GetUserFromServiceToApi(user model.User) *desc.GetResponse {
	var updatedAt *timestamppb.Timestamp
	if user.UpdatedAt != nil {
		updatedAt = timestamppb.New(*user.UpdatedAt)
	}

	return &desc.GetResponse{
		Uuid:      user.Uuid,
		Name:      user.Name,
		Email:     user.Email,
		Role:      desc.Role(user.Role),
		CreatedAt: timestamppb.New(user.CreatedAt),
		UpdatedAt: updatedAt,
	}
}
