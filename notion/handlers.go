package notion

import (
	"fmt"
)

type RoadMapGetter interface {
	GetRoadMap(pageID string) (*Roadmap, error)
}

func pendency(roadMap *Roadmap, publisher Publisher) {
	if roadMap.HasPendency() {
		publisher.Publish(&PendencyEvent{Pendency: roadMap.Pendency()})
	}
}

type studyInspectHandler struct {
	notionRepo  RoadMapGetter
	inspections []func(roadMap *Roadmap, publisher Publisher)
	publisher   Publisher
}

func NewStudyInspectHandler(notionRepo RoadMapGetter, publisher Publisher) *studyInspectHandler {
	return &studyInspectHandler{
		notionRepo: notionRepo,
		publisher:  publisher,
		inspections: []func(roadMap *Roadmap, publisher Publisher){
			pendency,
		},
	}
}

func (h *studyInspectHandler) Handle(message StudyInspectCmd) {
	roadMap, err := h.notionRepo.GetRoadMap(message.RoadmapID)

	if err != nil {
		fmt.Println("some error", err)
		return
	}

	for _, insp := range h.inspections {
		insp(roadMap, h.publisher)
	}

}
