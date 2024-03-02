package notion

import "notion-agenda/service"

type Publisher interface {
	Publish(message service.Message)
}

func StudyInspect(publisher Publisher) {
	cmd := &StudyInspectCmd{}
	publisher.Publish(cmd)
}
