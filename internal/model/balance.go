package model

import "time"

type Balance struct {
	Id        string    `json:"id" db:"id"`
	UserId    string    `json:"user_id" db:"user_id"`
	Current   float32   `json:"current" db:"current"`
	Withdrawn float32   `json:"withdrawn" db:"withdrawn"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}
