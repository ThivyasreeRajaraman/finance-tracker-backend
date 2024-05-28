package initializers

import (
	"log"

	"github.com/joho/godotenv"
)

func LoadEnvVariables(isTest bool) {
	var envPath string
	if isTest {
		envPath = "../../.env"
	} else {
		envPath = "./.env"
	}

	err := godotenv.Load(envPath)
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}
}
