package model

import "time"

type Inventory struct {
	ID        string    `json:"id"`
	ProductID string    `json:"productId"`
	Quantity  float64   `json:"quantity"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}
