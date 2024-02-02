package model

import "time"

type CustomerAmount struct {
	ID         string    `json:"id"`
	CustomerID string    `json:"customerId"`
	Amount     float64       `json:"amount"`
	CreatedAt  time.Time `json:"createdAt"`
	ModifiedAt time.Time `json:"modifiedAt"`
}
