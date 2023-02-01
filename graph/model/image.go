package model

type Image struct {
	ID  string `json:"id" gorm:"primaryKey"`
	Url string `json:"url"`
	Alt string `json:"alt"`
}
