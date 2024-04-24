package models

type ActionModel struct {
	Type   string `json:"type"`
	Object any    `json:"subtype"`
}
