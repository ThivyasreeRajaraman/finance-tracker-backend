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
	)
	if err != nil {
		log.Fatalf("Error migrating database: %v", err)
	}

	if err := models.MigrateCategories(DB); err != nil {
		log.Fatalf("Error migrating categories table: %v", err)
	}
}
