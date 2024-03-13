package notion

import (
	"github.com/google/uuid"

	"notion-agenda/src/service"
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

type PriorityEvent struct {
	Priorities []StudyStep
}

func (e *PriorityEvent) Meta() service.MessageMeta {
	return service.MessageMeta{
		Id:    uuid.New(),
		Topic: "notion_study_attention",
	}
}
