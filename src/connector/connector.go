package connector

import "oh-my-chat/src/models"

type Connector interface {
	Acquire(input chan<- models.Message)
	Dispatch(message models.Message)
}
