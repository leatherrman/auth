package user

import (
	"context"
	"errors"

	"github.com/katyafirstova/auth_service/internal/model"
)

func (s *serv) Create(ctx context.Context, req model.CreateUser) (string, error) {
	if req.Password != req.PasswordConfirm {
		return "", errors.New("passwords are not equal")
	}

	uuid, err := s.userRepository.Create(ctx, req)
	if err != nil {
		return "", err
	}

	return uuid, nil
}
