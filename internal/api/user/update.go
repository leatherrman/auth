package user

import (
	"context"

	"google.golang.org/protobuf/types/known/emptypb"

	"github.com/katyafirstova/auth_service/internal/converter"
	"github.com/katyafirstova/auth_service/pkg/user_v1"
)

func (i *Implementation) Update(ctx context.Context, req *user_v1.UpdateRequest) (*emptypb.Empty, error) {
	err := i.userService.Update(ctx, req.Uuid, converter.UpdateUserToServiceFromApi(req))
	if err != nil {
		return nil, err
	}

	return nil, nil
}
