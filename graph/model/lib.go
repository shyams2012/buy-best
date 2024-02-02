package model

import (
	"fmt"
	"log"
	"reflect"
	"strings"

	"gorm.io/gorm"
)

func DeleteObject(db *gorm.DB, object interface{}, id string) (bool, error) {
	name := strings.ToLower(reflect.ValueOf(object).Elem().Type().Name())
	if tx := db.First(object, "id = ?", id); tx.Error != nil {
		log.Print(tx.Error)
		return false, fmt.Errorf("%s not found, id='%s'", name, id)
	}
	if tx := db.Delete(object); tx.Error != nil {
		log.Print(tx.Error)
		return false, fmt.Errorf("error deleting %s, id='%s'", name, id)
	}

	return true, nil
}


