package query

import "notion-agenda/src/notion"

type ReadModeRepo interface {
	GetRoadMap(pageID string) (*notion.Roadmap, error)
}
