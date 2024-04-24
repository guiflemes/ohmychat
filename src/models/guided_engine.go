package models

type ActionModel struct {
	Type   string `json:"type"`
	Object any    `json:"subtype"`
}

type OptionModel struct {
	Name   string      `json:"name"`
	Intent string      `json:"intent"`
	Action ActionModel `json:"action,omitempty"`
}

type IntentModel struct {
	Name    string        `json:"intent"`
	Options []OptionModel `json:"options"`
}
