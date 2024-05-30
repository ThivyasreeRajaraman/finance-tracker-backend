package models

import (
	"time"

	"gorm.io/gorm"
)

type RecurringExpense struct {
	gorm.Model
	UserID          uint       `json:"user_id"`
	User            User       `gorm:"foreignkey:UserID;association_foreignkey:UserID"`
	CategoryID      uint       `json:"category_id"`
	Category        Categories `gorm:"foreignkey:CategoryID;association_foreignkey:CategoryID"`
	Amount          uint       `json:"amount"`
	Frequency       string     `json:"frequency"`
	NextExpenseDate time.Time  `json:"next_expense_date" gorm:"type:date"`
}

func MigrateRecurringExpense(db *gorm.DB) error {
	if !db.Migrator().HasIndex(&RecurringExpense{}, "idx_category_id_user_id") {
		if err := db.Exec("CREATE UNIQUE INDEX idx_category_id_user_id ON recurring_expenses (category_id, user_id) WHERE deleted_at IS NULL").Error; err != nil {
			return err
		}
	}
	return nil
}
