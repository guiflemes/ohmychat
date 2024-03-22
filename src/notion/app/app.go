package app

import "notion-agenda/src/notion/app/query"

type NotionApplication struct {
	Queries Queries
}

type Queries struct {
	PendencyTasks query.PendendyTasksHandler
}
