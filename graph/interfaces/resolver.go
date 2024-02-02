package interfaces

//go:generate go run github.com/99designs/gqlgen

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

import (
	"fmt"
	"log"
	"reflect"
	"strings"
	"time"

	"github.com/iancoleman/strcase"
	"github.com/shyams2012/buy-best/graph/model"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type Resolver struct {
	database *gorm.DB
}

func NewResolver(dsn string, config *gorm.Config, initialize bool) *Resolver {
	resolver := &Resolver{}
	resolver._Init(dsn, config, initialize)

	return resolver
}

func (r *Resolver) _Init(dsn string, config *gorm.Config, initialize bool) {
	db, err := gorm.Open(mysql.Open(dsn), config)
	if err != nil {
		panic(err)
	}

	if initialize {

		err = db.AutoMigrate(
			&model.User{},
			&model.Product{},
			&model.Transaction{},
			&model.CustomerCart{},
			&model.CustomerAmount{},
			&model.Image{},
			&model.Charge{},
			&model.Inventory{},
		)

		if err != nil {
			panic(err)
		}

		var user model.User
		adminUserId := "admin"
		tx := db.Where("id = ?", adminUserId).First(&user)
		if tx.Error != nil {
			if tx.Error == gorm.ErrRecordNotFound {
				at := time.Now()
				adminUser := model.User{
					ID:             adminUserId,
					CreatedAt:      at,
					ModifiedAt:     at,
					Username:       "admin",
					Role:           "ADMIN",
					RefreshCounter: 0,
					RefreshedTill:  0,
					IsActive:       true,
				}
				adminUser.SetPassword("admin")
				if tx := db.Create(adminUser); tx.Error != nil {
					log.Println("[WARN]", tx.Error)
				}
			} else {
				log.Println("[WARN]", tx.Error)
			}

		}
	}

	r.database = db
}

func (r *Resolver) DB() *gorm.DB {
	return r.database
}

func (r *Resolver) DBWithFilter(cond interface{}) *gorm.DB {
	dbWithFilter := r.DB().Limit(100)

	if cond == nil {
		return dbWithFilter
	}

	filter := reflect.ValueOf(cond)
	if filter.IsValid() && filter.Kind() == reflect.Ptr {
		filter = filter.Elem()
	}
	if !filter.IsValid() {
		return dbWithFilter
	}

	filterType := filter.Type()
	for i := 0; i < filter.NumField(); i++ {
		fieldValue := filter.Field(i).Elem()
		if fieldValue.IsValid() && fieldValue.Kind() == reflect.Ptr {
			fieldValue = fieldValue.Elem()
		}

		if !fieldValue.IsValid() {
			continue
		}

		modelName := strings.TrimSuffix(filterType.Name(), "Filter")
		tableName := fmt.Sprintf("%ss", strcase.ToSnake(modelName))

		fieldName := strings.TrimSpace(strcase.ToSnake(filterType.Field(i).Name))
		if fieldName == "query" {
			if fieldValueStr, ok := fieldValue.Interface().(string); ok {
				fieldValueStr = strings.TrimSpace(fieldValueStr)
				if fieldValueStr != "" {
					searchValue := fmt.Sprintf("%%%v%%", fieldValueStr)
					whereCond := fmt.Sprintf("%s::text ILIKE ?", tableName)

					dbWithFilter = dbWithFilter.Where(whereCond, searchValue)
				}
			}
		} else {
			where := fmt.Sprintf("%s.%s = ?", tableName, fieldName)
			dbWithFilter = dbWithFilter.Where(where, fieldValue.Interface())
		}
	}

	return dbWithFilter
}
