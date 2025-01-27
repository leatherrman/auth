package user

import (
	"context"

	"github.com/katyafirstova/auth_service/internal/converter"
	"github.com/katyafirstova/auth_service/pkg/user_v1"
)

func (i *Implementation) Get(ctx context.Context, req *user_v1.GetRequest) (*user_v1.GetResponse, error) {
	user, err := i.userService.Get(ctx, req.Uuid)
	if err != nil {
		return nil, err
	}

	return converter.GetUserFromServiceToApi(user), nil
}
