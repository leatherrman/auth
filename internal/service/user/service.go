package user

import (
	"github.com/katyafirstova/auth_service/internal/repository"
	"github.com/katyafirstova/auth_service/internal/service"
)

type serv struct {
	userRepository repository.UserRepository
}

func NewService(
	userRepository repository.UserRepository,
) service.UserService {
	return &serv{
		userRepository: userRepository,
	}
}
