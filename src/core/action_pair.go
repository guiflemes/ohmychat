package core

import (
	"oh-my-chat/src/models"
)

type ActionReplyPair struct {
	ReplyTo chan<- models.Message
	Action  Action
	Input   models.Message
}
