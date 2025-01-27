package user

import (
	"github.com/katyafirstova/auth_service/internal/repository"
)

type serv struct {
	userRepository repository.UserRepository
}

func NewService(
	userRepository repository.UserRepository,
) *serv {
	return &serv{
		userRepository: userRepository,
	}
}
