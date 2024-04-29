package models

type OptionModel struct {
	Name    string       `yaml:"name"`
	Key     string       `yaml:"key"`
	Content string       `yaml:"content"`
	Action  *ActionModel `yaml:"action,omitempty"`
}

type IntentModel struct {
	Name    string        `yaml:"name"`
	Key     string        `yaml:"key"`
	Options []OptionModel `yaml:"options"`
}

type WorkFlowGuided struct {
	Engine  string        `yaml:"engine"`
	Intents []IntentModel `yaml:"intents"`
}
