package app

import "oh-my-chat/src/notion/app/query"

type NotionApplication struct {
	Queries Queries
}

type Queries struct {
	PendencyTasks query.PendendyTasksHandler
}
