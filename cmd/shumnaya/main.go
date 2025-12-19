package main

import (
	"log"

	"shumnaya/internal/config"

	"github.com/joho/godotenv"
)

func main() {

	err := godotenv.Load()
	if err != nil {
		log.Println("Warning: .env file not found, using system environment variables")
	}

	config.ConnectDB()

}
