package dbhelper

import (
	"github.com/Thivyasree-Rajaraman/finance-tracker/initializers"
)

func GenericCreate[T any](model *T) error {
	return initializers.DB.Create(model).Error
}

func GenericUpdate[T any](model *T) error {
	return initializers.DB.Save(model).Error
}
