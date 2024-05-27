package action

import "oh-my-chat/src/core"

type StorageActionMessage struct{}

func (s *StorageActionMessage) Pop() (*core.ActionReplyPair, bool) {
	return nil, false
}

func (s *StorageActionMessage) Put(core.ActionReplyPair) {}
