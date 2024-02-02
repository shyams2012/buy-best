package model

import "time"

type CustomerCart struct {
	ID         string    `json:"id"`
	CustomerID string    `json:"customerId"`
	ProductID  string    `json:"productId"`
	CreatedAt  time.Time `json:"createdAt"`
	ModifiedAt time.Time `json:"modifiedAt"`
	Amount     float64   `json:"amount"`
}
