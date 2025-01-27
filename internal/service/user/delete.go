package user

import (
	"context"
)

func (s *serv) Delete(ctx context.Context, uuid string) error {
	err := s.userRepository.Delete(ctx, uuid)
	if err != nil {
		return err
	}

	return nil
}
