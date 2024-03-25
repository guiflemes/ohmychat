package notion

import (
	"time"

	"oh-my-chat/src/utils"
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
	now = time.Date(
		time.Now().Year(),
		time.Now().Month(),
		time.Now().Day(),
		0,
		0,
		0,
		0,
		time.UTC,
	)
	expiredDays           int = 15
	needsInitializionDays int = 30
)

type Roadmap struct {
	ID    string
	Steps []StudyStep
}

func (r *Roadmap) StepCount() int {
	return len(r.Steps)
}

func (r *Roadmap) HasPendency() bool {
	return utils.Any(r.Steps, func(s StudyStep) bool {
		if s.Status == Pending {
			return true
		}
		return false
	})
}

func (r *Roadmap) Pendency() []StudyStep {
	pedendy := make([]StudyStep, 0)
	for _, s := range r.Steps {
		if s.Status == Pending {
			pedendy = append(pedendy, s)
		}
	}

	return pedendy
}

func (r *Roadmap) String() string {
	str := ""
	for _, i := range r.Steps {
		str += i.Name + ", "
	}
	return str
}

func (r *Roadmap) NeedsAttention() bool {
	for _, step := range r.Steps {
		if step.NeedsAttention() {
			return true
		}
	}
	return false
}

func (r *Roadmap) Priorities() []StudyStep {
	steps := make([]StudyStep, 0)
	for _, step := range r.Steps {
		if step.NeedsAttention() {
			steps = append(steps, step)
		}
	}
	return steps

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

func (s *StudyStep) NeedsAttention() bool {
	return s.ExpireSoon() && s.Priority == High
}

func (s *StudyStep) NeedsInitialization() bool {
	x := s.CreatedAt.AddDate(0, 0, needsInitializionDays)
	return x.Equal(now) || x.Before(now)
}

func (s *StudyStep) IsBlockAfterInitialization() bool {
	if s.Status != Blocked {
		return false
	}

	if s.StartedAt.IsZero() {
		return false
	}

	x := s.StartedAt.AddDate(0, 0, 15)
	return x.Equal(now) || x.Before(now)
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
