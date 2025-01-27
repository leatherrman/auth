package service

import (
	"context"

	"github.com/katyafirstova/auth_service/internal/model"
)

type UserService interface {
	Create(ctx context.Context, req model.CreateUser) (string, error)
	Get(ctx context.Context, uuid string) (model.User, error)
	Update(ctx context.Context, uuid string, req model.UpdateUser) error
	Delete(ctx context.Context, uuid string) error
}
