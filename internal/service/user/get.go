package user

import (
	"context"

	"github.com/katyafirstova/auth_service/internal/model"
)

func (s *serv) Get(ctx context.Context, uuid string) (model.User, error) {
	user, err := s.userRepository.Get(ctx, uuid)
	if err != nil {
		return model.User{}, err
	}

	return user, nil
}
