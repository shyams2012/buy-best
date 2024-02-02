package lib

import (
	"github.com/shyams2012/buy-best/graph/model"
	"gorm.io/gorm"
)

func SavePayment(db *gorm.DB, charge *model.Charge) (err error) {

	if err = db.Create(charge).Error; err != nil {
		return err
	}
	return nil

}
