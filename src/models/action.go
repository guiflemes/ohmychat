package models

type ActionModel struct {
	Type   string `yaml:"type"`
	Object any    `yaml:"object"`
}
