package notion

import (
	"notion-agenda/service"

	"github.com/google/uuid"
)

type PendencyEvent struct {
	Pendency []StudyStep
}

func (e *PendencyEvent) Meta() service.MessageMeta {
	return service.MessageMeta{
		Id:    uuid.New(),
		Topic: "notion_study_pendency",
	}
}
