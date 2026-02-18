package db

import (
	"context"
	"database/sql"

	"github.com/Foga2H/ya-go-gophermart/internal/logger"
	"github.com/Foga2H/ya-go-gophermart/internal/model"
	"github.com/google/uuid"
)

type OrderRepository struct {
	db     *sql.DB
	logger *logger.Logger
}

func NewOrderRepository(db *sql.DB, logger *logger.Logger) *OrderRepository {
	return &OrderRepository{db: db, logger: logger}
}

func (o *OrderRepository) Create(ctx context.Context, userID string, orderID string) (string, error) {
	var orderUuid string
	err := o.db.QueryRowContext(
		ctx,
		"INSERT INTO orders (id, user_id, number, sum, status) VALUES ($1, $2, $3, $4, $5) RETURNING id",
		uuid.NewString(),
		userID,
		orderID,
		0,
		model.OrderStatusNew,
	).Scan(&orderUuid)

	if err != nil {
		return "", err
	}

	return orderUuid, nil
}

func (o *OrderRepository) FindByOrderID(ctx context.Context, orderID string) (model.Order, error) {
	var order model.Order
	err := o.db.QueryRowContext(
		ctx,
		"SELECT id, user_id, number, sum, status, uploaded_at, updated_at FROM orders WHERE number = $1",
		orderID,
	).Scan(&order.Id, &order.UserId, &order.Number, &order.Sum, &order.Status, &order.UploadedAt, &order.UpdatedAt)

	if err != nil {
		return model.Order{}, err
	}

	return order, nil
}

func (o *OrderRepository) FindByUserID(ctx context.Context, userID string) ([]model.Order, error) {
	var orders []model.Order
	rows, err := o.db.QueryContext(
		ctx,
		"SELECT id, user_id, number, sum, status, uploaded_at, updated_at FROM orders WHERE user_id = $1 ORDER BY uploaded_at DESC",
		userID,
	)

	if err != nil {
		o.logger.Errorf("Failed to query for orders by user_id: %s", err.Error())
		return []model.Order{}, err
	}

	defer rows.Close()

	for rows.Next() {
		var order model.Order
		err := rows.Scan(&order.Id, &order.UserId, &order.Number, &order.Sum, &order.Status, &order.UploadedAt, &order.UpdatedAt)
		if err != nil {
			o.logger.Errorf("Error scanning row: %s", err)
			return []model.Order{}, err
		}
		orders = append(orders, order)
	}

	return orders, nil
}
