package repository

import (
	"context"

	"github.com/Foga2H/ya-go-gophermart/internal/model"
)

//go:generate mockgen -source=user.go -destination=mock/user_repository_mock.go -package=repositorymock
type UserRepository interface {
	Create(ctx context.Context, login string, hashedPassword string) (model.User, error)
	FindByLogin(ctx context.Context, login string) (model.User, error)
}
