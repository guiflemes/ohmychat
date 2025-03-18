package rule_engine

import "strings"

type CachedOptions struct {
	intent  string
	options map[string]Option
}

func NewCachedOptions(intent *Intent) CachedOptions {
	c := CachedOptions{
		intent: intent.Key,
	}
	options := make(map[string]Option)
	for _, option := range intent.Options.items {
		c.options[strings.ToLower(option.Name)] = option
	}
	c.options = options
	return c
}

func (e *CachedOptions) GetOption(input string) (*Option, bool) {
	input = strings.ToLower(input)
	if opt, ok := e.options[input]; ok {
		return &opt, true
	}
	return nil, false
}
