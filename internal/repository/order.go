package repository

import (
	"context"

	"github.com/Foga2H/ya-go-gophermart/internal/model"
)

type OrderRepository interface {
	Create(ctx context.Context, userID string, orderID string) (string, error)
	FindByOrderID(ctx context.Context, orderID string) (model.Order, error)
	FindByUserID(ctx context.Context, userID string) ([]model.Order, error)
}
