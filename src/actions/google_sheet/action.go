package googlesheet

import (
	"context"
	"oh-my-chat/src/models"
)

type GoogleSheetAction struct {
}

func NewGoogleSheetAction(model *models.GoogleSheetModel) *GoogleSheetAction {
	return &GoogleSheetAction{}
}

func (a *GoogleSheetAction) Handle(ctx context.Context, message *models.Message) error { return nil }
