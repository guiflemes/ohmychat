package queue

import "oh-my-chat/src/models"

type QueueMessage struct {
	Message    models.Message `json:"message"`
	ActionName string         `json:"action_name"`
	ModelName  string         `json:"model_name"`
	ModelData  any
}
