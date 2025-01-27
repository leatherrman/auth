package user

import (
	"context"

	"google.golang.org/protobuf/types/known/emptypb"

	"github.com/katyafirstova/auth_service/pkg/user_v1"
)

func (i *Implementation) Delete(ctx context.Context, req *user_v1.DeleteRequest) (*emptypb.Empty, error) {
	err := i.userService.Delete(ctx, req.Uuid)
	if err != nil {
		return nil, err
	}

	return nil, nil
}
