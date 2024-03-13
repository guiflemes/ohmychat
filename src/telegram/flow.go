package telegram

import "sync"

type Card struct {
	ID          string `json:"id"`
	Type        string `json:"type"`
	Initial     bool   `json:"initial"`
	ExpectedMsg string `json:"expected_msg"`
	Template    any    `json:"template"`
}

type Relationship struct {
	SourceCardID      string `json:"source_card_id"`
	TargetCardID      string `json:"target_card_id"`
	RelationshipType  string `json:"relationship_type"`
	AdditionalDetails string `json:"additional_details"`
}

type Flow struct {
	Name          string          `json:"name"`
	Key           string          `json:"key"`
	Cards         map[string]Card `json:"cards"`
	Relationships []Relationship  `json:"relationships"`
	lock          sync.RWMutex
}

func (f *Flow) Lock() {
	f.lock.Lock()
}

func (f *Flow) Unlock() {
	f.lock.Unlock()
}
