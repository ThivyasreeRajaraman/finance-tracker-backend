package models

import "gorm.io/gorm"

type Categories struct {
	gorm.Model
	Name   string `json:"name"`
	Type   string `json:"type"`
	UserID *uint  `json:"user_id"`
	User   User   `gorm:"foreignkey:UserID;association_foreignkey:UserID"`
}

func MigrateCategories(db *gorm.DB) error {
	if !db.Migrator().HasIndex(&Categories{}, "idx_name_null_user_id") {
		if err := db.Exec("CREATE UNIQUE INDEX idx_name_null_user_id ON categories (name) WHERE user_id IS NULL AND deleted_at IS NULL").Error; err != nil {
			return err
		}
	}

	if !db.Migrator().HasIndex(&Categories{}, "idx_name_user_id") {
		if err := db.Exec("CREATE UNIQUE INDEX idx_name_user_id ON categories (name, user_id) WHERE user_id IS NOT NULL AND deleted_at IS NULL").Error; err != nil {
			return err
		}
	}

	return nil
}
