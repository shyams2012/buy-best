package model

type Product struct {
	ID          string  `json:"id"`
	Name        string  `json:"name"`
	Model       string  `json:"model"`
	Price       float64 `json:"price"`
	Description string  `json:"description"`
}
