package main

import (
	"log"

	"github.com/joho/godotenv"

	"oh-my-chat/src/core"
)

func main() {

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

}

func Run() {
	actionQueue := core.NewGoActionQueue()
  actionQueue.Consume()
  
  guidedEngine = &core.GuidedResponseEngine{}

  core.NewProcessor(workflowGetter core.WorkflowGetter, engines core.Engines)
}
