package service

import (
	"context"

	"github.com/Foga2H/ya-go-gophermart/internal/logger"
	"github.com/Foga2H/ya-go-gophermart/internal/model"
	"github.com/Foga2H/ya-go-gophermart/internal/repository"
)

type GetUserOrdersService struct {
	repository repository.OrderRepository
	logger     *logger.Logger
}

func NewGetUserOrdersService(orderRepository repository.OrderRepository, logger *logger.Logger) *GetUserOrdersService {
	return &GetUserOrdersService{
		repository: orderRepository,
		logger:     logger,
	}
}

func (o *GetUserOrdersService) GetOrdersByUserID(ctx context.Context, userID string) ([]model.Order, error) {
	orders, err := o.repository.FindByUserID(ctx, userID)
	if err != nil {
		o.logger.Errorf("user orders err: %v", err)
		return nil, err
	}

	return orders, nil
}
