package main

import (
	"log"

	"github.com/joho/godotenv"

	"oh-my-chat/src/telegram"
)

func main() {

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	telegram.NewEngine().Chating(30)
}
