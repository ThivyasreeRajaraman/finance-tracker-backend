package dbhelper

import (
	"reflect"
	"strings"

	"github.com/Thivyasree-Rajaraman/finance-tracker/initializers"
)

func GenericCreate[T any](model *T) error {
	return initializers.DB.Create(model).Error
}

func GenericUpdate[T any](model *T) error {
	return initializers.DB.Save(model).Error
}

func GenericDelete[T any](model *T) error {
	return initializers.DB.Delete(&model).Error
}

func FetchDataWithPagination[T any](model *T, page, limit int, conditions map[string]interface{}, orderBy string) ([]T, int, error) {
	offset := (page - 1) * limit

	var data []T
	var totalCount int64
	db := initializers.DB.Model(model).Offset(offset).Limit(limit)
	for key, value := range conditions {
		db = db.Where(key, value)
	}

	if orderBy != "" {
		db = db.Order(orderBy)
	}

	modelType := reflect.TypeOf(model).Elem()
	for i := 0; i < modelType.NumField(); i++ {
		field := modelType.Field(i)
		if field.Type.Kind() == reflect.Struct {
			gormTag := field.Tag.Get("gorm")
			if gormTag != "" {
				if strings.Contains(gormTag, "foreignkey") {
					db = db.Preload(field.Name)
				}
			}
		}
	}

	if err := db.Find(&data).Error; err != nil {
		return nil, 0, err
	}
	countDB := initializers.DB.Model(model)
	for key, value := range conditions {
		countDB = countDB.Where(key, value)
	}
	if err := countDB.Count(&totalCount).Error; err != nil {
		return nil, 0, err
	}
	return data, int(totalCount), nil
}
