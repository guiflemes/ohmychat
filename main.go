package main

import (
	"log"

	"github.com/joho/godotenv"

	"notion-agenda/src/notion"
	"notion-agenda/src/service"
	"notion-agenda/src/telegram"
)

func main() {

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	Run()
}

func Run() {

	bus := service.NewBus()
	bus.SetHandler(
		"notion_inspect_study_road_map",
		notion.NewStudyInspectHandler(&notion.SketchRepo{}, bus),
	)
	bus.SetHandler("notion_study_pendency", telegram.NewTelegramPendencyHandler())

	bus.Consume()

	notion.StudyInspect(bus)

	select {}

}
