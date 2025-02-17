package user

import (
	"github.com/katyafirstova/auth_service/internal/service"
	"github.com/katyafirstova/auth_service/pkg/user_v1"
)

type Implementation struct {
	user_v1.UnimplementedUserV1Server
	userService service.UserService
}

func NewImplementation(userService service.UserService) *Implementation {
	return &Implementation{userService: userService}

}
