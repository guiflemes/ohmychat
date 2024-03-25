package query

import "oh-my-chat/src/notion"

type ReadModeRepo interface {
	GetRoadMap(pageID string) (*notion.Roadmap, error)
}
