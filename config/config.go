package config

import (
	"fmt"

	"github.com/joho/godotenv"
)

func InitEnvironmentFile() {
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Error loading .env file")
	}
}
