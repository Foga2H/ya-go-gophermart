package model

import "time"

type OrderStatus string

const (
	OrderStatusNew        OrderStatus = "NEW"
	OrderStatusProcessing OrderStatus = "PROCESSING"
	OrderStatusInvalid    OrderStatus = "INVALID"
	OrderStatusProcessed  OrderStatus = "PROCESSED"
)

type Order struct {
	Id         string      `json:"id" db:"id"`
	UserId     string      `json:"user_id" db:"user_id"`
	Number     string      `json:"number" db:"number"`
	Sum        float32     `json:"sum" db:"sum"`
	Status     OrderStatus `json:"status" db:"status"`
	UploadedAt time.Time   `json:"uploaded_at" db:"uploaded_at"`
	UpdatedAt  time.Time   `json:"updated_at" db:"updated_at"`
}
