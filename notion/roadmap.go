package notion

import (
	"time"
)

type Priority string
type Status string

const (
	High   Priority = "High"
	Medium Priority = "Medium"
	Low    Priority = "Low"
)

const (
	Blocked Status = "Blocked"
	Process Status = "Progress"
)

type Roadmap struct {
	ID    string
	Steps []StudyStep
}

func (r *Roadmap) StepCount() int {
	return len(r.Steps)
}

func (r *Roadmap) String() string {
	str := ""
	for _, i := range r.Steps {
		str += i.Name + ", "
	}
	return str
}

func FromSlice[T any](ID string, items []T, stepBuilder func(item T) StudyStep) *Roadmap {
	steps := make([]StudyStep, 0, len(items))
	for _, item := range items {
		step := stepBuilder(item)
		steps = append(steps, step)
	}
	return &Roadmap{
		ID:    ID,
		Steps: steps,
	}
}

type StudyStep struct {
	ID        string
	Name      string
	Category  string
	Link      *string
	Status    Status
	Notes     *string
	Deadline  time.Time
	Priority  Priority
	CreatedAt time.Time
	Type      []string
	StartedAt time.Time
	Points    int
}

func (s *StudyStep) String() string {
	return s.Name
}
