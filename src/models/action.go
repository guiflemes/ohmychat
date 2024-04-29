package models

type Model interface {
	GetType() string
}

type ActionModel struct {
	Type   string `yaml:"type"`
	Object any    `yaml:"object"`
}
