package notion

import (
	"fmt"
	"notion-agenda/utils"
)

type RoadMapGetter interface {
	GetRoadMap(pageID string) (*Roadmap, error)
}

type StudyEvent struct {
	name string
	step StudyStep
}

type ProcessHandler struct {
	notionRepo RoadMapGetter
	processing []func(roadMap *Roadmap)
	events     []StudyStep
}

func NewProcessHandler() *ProcessHandler {
	events := make([]StudyEvent, 0)

	return &ProcessHandler{
		notionRepo: nil,
		processing: []func(roadMap *Roadmap){
			func(roadMap *Roadmap) {
				if roadMap.HasPendency() {
					e := utils.Map(roadMap.Pendency(), func(s StudyStep) StudyEvent { return StudyEvent{} })
					events = append(events, e...)
				}
			},
		},
	}
}

func (h *ProcessHandler) Process() {
	roadMap, err := h.notionRepo.GetRoadMap("")

	if err != nil {
		fmt.Println("some error", err)
		return
	}

	for _, proc := range h.processing {
		proc(roadMap)
	}

}
