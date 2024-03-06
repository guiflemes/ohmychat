package main

import (
	"fmt"
	"log"
	"notion-agenda/src/notion"
	"notion-agenda/src/service"
	"notion-agenda/src/telegram"

	"github.com/joho/godotenv"
)

func main() {

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	telegram.Tele()

}

type myHandler struct{}

func (h *myHandler) Handle(message service.Message) error {
	fmt.Println("handler event")
	return nil
}

func sketchRepo() {

	bus := service.NewBus()
	bus.SetHandler("notion_inspect_study_road_map", notion.NewStudyInspectHandler(&notion.SketchRepo{}, bus))
	bus.SetHandler("notion_study_pendency", &myHandler{})

	bus.Consume()

	notion.StudyInspect(bus)

	select {}

}
