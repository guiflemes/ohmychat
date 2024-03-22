package query

import (
	"context"
	"log"

	"notion-agenda/src/notion"
)

type Pendendy struct {
	PageID string
}

type PendendyTasksHandler struct {
	readMode ReadModeRepo
}

func NewPendencyTaskHandler(readModeRepo ReadModeRepo) *PendendyTasksHandler {
	return &PendendyTasksHandler{readModeRepo}
}

func (p *PendendyTasksHandler) Handler(
	ctx context.Context,
	pendency Pendendy,
) ([]notion.StudyStep, error) {
	roadMap, err := p.readMode.GetRoadMap(pendency.PageID)
	if err != nil {
		log.Println("error", err)
		return nil, err
	}

	if !roadMap.HasPendency() {
		return nil, nil
	}

	return roadMap.Pendency(), nil
}
