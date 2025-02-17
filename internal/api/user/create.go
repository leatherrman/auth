package user

import (
	"context"

	"github.com/katyafirstova/auth_service/internal/converter"
	"github.com/katyafirstova/auth_service/pkg/user_v1"
)

func (i *Implementation) Create(ctx context.Context, req *user_v1.CreateRequest) (*user_v1.CreateResponse, error) {
	uuid, err := i.userService.Create(ctx, converter.CreateUserToServiceFromAPI(req))
	if err != nil {
		return nil, err
	}

	return &user_v1.CreateResponse{Uuid: uuid}, nil
}
