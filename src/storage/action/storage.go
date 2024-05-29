package action

import "oh-my-chat/src/core"

type StorageActionMessage interface {
	Dequeue() (core.ActionReplyPair, bool)
	Enqueue(core.ActionReplyPair)
}
