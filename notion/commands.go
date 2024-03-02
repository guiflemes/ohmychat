package notion

import (
	"notion-agenda/service"

	"github.com/google/uuid"
)

type StudyInspectCmd struct {
	RoadmapID string
}

func (c *StudyInspectCmd) Meta() service.MessageMeta {
	return service.MessageMeta{
		Id:    uuid.New(),
		Topic: "notion_inspect_study_road_map",
	}
}
