package model

import "time"

type Transaction struct {
	ID                string          `json:"id"`
	UserId            string          `json:"userId"`
	CustomerProductID string          `json:"customerProductId" gorm:"size:191"`
	Price             float64         `json:"price"`
	Name              string          `json:"name"`
	Type              TransactionType `json:"type"`
	CreatedAt         time.Time       `json:"createdAt"`
}
