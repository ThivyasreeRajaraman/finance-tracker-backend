package models

import (
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
	NextExpenseDate string     `json:"next_expense_date"`
	Currency        string     `json:"currency"`
	Active          bool       `json:"active"`
}

func MigrateRecurringExpense(db *gorm.DB) error {
	if !db.Migrator().HasIndex(&RecurringExpense{}, "idx_category_id_user_id") {
		if err := db.Exec("CREATE UNIQUE INDEX idx_category_id_user_id ON recurring_expenses (category_id, user_id) WHERE deleted_at IS NULL").Error; err != nil {
			return err
		}
	}
	return nil
}

func (r *RecurringExpense) BeforeCreate(tx *gorm.DB) (err error) {
	r.Active = true
	return nil
}
