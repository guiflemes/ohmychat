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
	Blocked  Status = "Blocked"
	Process  Status = "Progress"
	Pending  Status = "Pending"
	Finished Status = "Finished"
)

var (
	now         = time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day(), 0, 0, 0, 0, time.UTC)
	expiredDays = 15
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

type StudyStep struct {
	ID        string
	Name      string
	Category  string
	Link      string
	Status    Status
	Notes     string
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

func (s *StudyStep) IsExpired() bool {
	if s.Deadline.IsZero() {
		return false
	}

	return s.Deadline.Equal(now) || s.Deadline.Before((now))
}

func (s *StudyStep) ExpireSoon() bool {
	if s.Status == Finished {
		return false
	}

	if s.IsExpired() {
		return false
	}

	if s.Deadline.IsZero() {
		return false
	}

	x := now.AddDate(0, 0, expiredDays)

	return x.Equal(s.Deadline) || x.After(s.Deadline)
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
