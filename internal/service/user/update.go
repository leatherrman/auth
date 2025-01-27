package user

import (
	"context"

	"github.com/katyafirstova/auth_service/internal/model"
)

func (s *serv) Update(ctx context.Context, uuid string, req model.UpdateUser) error {
	err := s.userRepository.Update(ctx, uuid, req)
	if err != nil {
		return err
	}

	return nil
}
