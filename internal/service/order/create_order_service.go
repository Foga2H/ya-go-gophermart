package service

import (
	"context"
	"database/sql"
	"errors"

	"github.com/Foga2H/ya-go-gophermart/internal/logger"
	"github.com/Foga2H/ya-go-gophermart/internal/model"
	"github.com/Foga2H/ya-go-gophermart/internal/repository"
)

var ErrOrderNotFound = errors.New("order not found")

type CreateOrderService struct {
	repository repository.OrderRepository
	logger     *logger.Logger
}

func NewCreateOrderService(orderRepository repository.OrderRepository, logger *logger.Logger) *CreateOrderService {
	return &CreateOrderService{
		repository: orderRepository,
		logger:     logger,
	}
}

func (s *CreateOrderService) Create(ctx context.Context, userID, orderID string) (string, error) {
	orderUUID, err := s.repository.Create(ctx, userID, orderID)
	if err != nil {
		s.logger.Errorf("create order: %v", err)
		return "", err
	}

	return orderUUID, nil
}

func (s *CreateOrderService) FindByOrderID(ctx context.Context, orderID string) (model.Order, error) {
	order, err := s.repository.FindByOrderID(ctx, orderID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return model.Order{}, ErrOrderNotFound
		}

		s.logger.Errorf("find by order id failed: %v", err)
		return model.Order{}, err
	}

	return order, nil
}
