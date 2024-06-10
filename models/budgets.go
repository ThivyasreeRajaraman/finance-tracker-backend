package models

import "gorm.io/gorm"

type Budgets struct {
	gorm.Model
	UserID     uint       `json:"user_id"`
	User       User       `gorm:"foreignkey:UserID;association_foreignkey:UserID"`
	CategoryID uint       `json:"category_id"`
	Category   Categories `gorm:"foreignkey:CategoryID;association_foreignkey:CategoryID"`
	Amount     uint       `json:"amount"`
	Threshold  uint       `json:"threshold"`
	Currency   string     `json:"currency"`
}

func MigrateBudgets(db *gorm.DB) error {
	if !db.Migrator().HasIndex(&Budgets{}, "idx_caytegory_user_id") {
		if err := db.Exec("CREATE UNIQUE INDEX idx_caytegory_user_id ON budgets (category_id, user_id) WHERE deleted_at IS NULL").Error; err != nil {
			return err
		}
	}
	return nil
}
