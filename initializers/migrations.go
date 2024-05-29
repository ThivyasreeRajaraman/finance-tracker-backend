package initializers

import (
	"log"

	"github.com/Thivyasree-Rajaraman/finance-tracker/models"
)

func SyncDatabase() {
	err := DB.AutoMigrate(
		&models.User{},
		&models.Categories{},
		&models.Budgets{},
		&models.Transaction{},
		&models.RecurringExpense{},
	)
	if err != nil {
		log.Fatalf("Error migrating database: %v", err)
	}

	if err := models.MigrateCategories(DB); err != nil {
		log.Fatalf("Error migrating categories table: %v", err)
	}
	if err := models.MigrateRecurringExpense(DB); err != nil {
		log.Fatalf("Error migrating recurring expense table: %v", err)
	}
	if err := models.MigrateBudgets(DB); err != nil {
		log.Fatalf("Error migrating budgets table: %v", err)
	}
}
