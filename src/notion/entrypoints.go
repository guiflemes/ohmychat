package notion

import "notion-agenda/service"

type Publisher interface {
	Publish(message service.Message)
}

func StudyInspect(publisher Publisher) {
	cmd := &StudyInspectCmd{RoadmapID: "037c048f1a9d4a1e88b54de89d0c58c0"}
	publisher.Publish(cmd)
}
