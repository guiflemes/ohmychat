package action

import (
	"sync"

	"oh-my-chat/src/core"
)

type memoryQueue struct {
	actions []core.ActionReplyPair
	l       sync.Mutex
}

func NewMemoryQueue() *memoryQueue {
	return &memoryQueue{
		actions: make([]core.ActionReplyPair, 0),
	}
}

func (q *memoryQueue) Enqueue(actionPair core.ActionReplyPair) {
	q.l.Lock()
	defer q.l.Unlock()
	q.actions = append(q.actions, actionPair)
}

func (q *memoryQueue) Dequeue() (core.ActionReplyPair, bool) {
	if len(q.actions) == 0 {
		return core.ActionReplyPair{}, false
	}
	q.l.Lock()
	defer q.l.Unlock()

	actionPair := q.actions[0]
	q.actions = q.actions[1:]
	return actionPair, true
}
