package rule_engine

import (
	"strings"
)

type Action struct {
	Name string
}

type Option struct {
	Content  string `yaml:"content"`
	Name     string `yaml:"name"`
	Key      string `yaml:"key"`
	Response string `yaml:"response"`
}

type Intent struct {
	Name     string   `yaml:"name"`
	Key      string   `yaml:"key"`
	Patterns []string `yaml:"patterns"`
	Options  Options  `yaml:"options"`
	Response string   `yaml:"response"`
}

type Options struct {
	items []Option
}

func (i *Intent) HasOptions() bool {
	return len(i.Options.items) > 0
}

type Intents map[string]*Intent

func (i Intents) GetIntent(pattern string) (*Intent, bool) {
	pattern = strings.ToLower(pattern)
	if intent, ok := i[pattern]; ok {
		return intent, true
	}
	return nil, false
}
