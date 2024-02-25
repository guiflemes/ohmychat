package main

import (
	"fmt"
	"log"
	"notion-agenda/notion"
	"notion-agenda/settings"

	"github.com/joho/godotenv"
)

func main() {

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	r, err := notion.SketchRepo(settings.GETENV("PAGE_ID"))
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(r)
}
