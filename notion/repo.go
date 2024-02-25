package notion

import (
	"context"
	"notion-agenda/settings"
	"strconv"
	"time"

	"github.com/jomei/notionapi"
)

func SketchRepo(pageID string) (*Roadmap, error) {
	client := notionapi.NewClient(notionapi.Token(settings.GETENV("NOTION_API")))
	dbQuery, err := client.Database.Query(context.Background(), notionapi.DatabaseID(pageID), nil)

	if err != nil {
		return nil, err
	}

	pages := dbQuery.Results

	roadmap := FromSlice[notionapi.Page](pageID, pages, func(item notionapi.Page) StudyStep {
		builder := NewStepBuilder()
		for desc, prop := range item.Properties {
			builder.SetProperty(desc, prop)
		}
		return builder.GetStepStudy()
	})

	return roadmap, nil
}

type stepBuilder struct {
	step StudyStep
}

func NewStepBuilder() *stepBuilder {
	builder := &stepBuilder{
		step: StudyStep{},
	}

	return builder
}

func (s *stepBuilder) GetStepStudy() StudyStep {
	return s.step
}

func (s *stepBuilder) builders() map[string]func(prop notionapi.Property) {
	//TODO map name on yml config file
	return map[string]func(prop notionapi.Property){
		"Name":        s.setName,
		"Category":    s.setCategory,
		"Notes":       s.setNotes,
		"Deadline":    s.setDeadline,
		"Links":       s.setLink,
		"Pontos":      s.setPoints,
		"Iniciado em": s.setStartAt,
		"Priority":    s.setPriority,
		"Type":        s.setType,
		"Criado em":   s.setCreatedAt,
		"Status":      s.setStatus,
	}
}

func (s *stepBuilder) SetProperty(name string, prop notionapi.Property) {
	builders := s.builders()

	fn, ok := builders[name]
	if ok {
		fn(prop)
	}
}

func (s *stepBuilder) setName(prop notionapi.Property) {
	text := prop.(*notionapi.TitleProperty).Title

	if len(text) == 0 {
		return
	}

	s.step.Name = text[0].PlainText
}

func (s *stepBuilder) setCategory(prop notionapi.Property) {
	option := prop.(*notionapi.SelectProperty).Select
	s.step.Category = option.Name
}

func (s *stepBuilder) setLink(prop notionapi.Property) {
	option := prop.(*notionapi.URLProperty)
	s.step.Link = &option.URL
}

func (s *stepBuilder) setStatus(prop notionapi.Property) {
	option := prop.(*notionapi.SelectProperty).Select
	s.step.Status = Status(option.Name)
}

func (s *stepBuilder) setNotes(prop notionapi.Property) {
	text := prop.(*notionapi.RichTextProperty).RichText

	if len(text) == 0 {
		return
	}

	s.step.Notes = &text[0].PlainText
}

func (s *stepBuilder) setPoints(prop notionapi.Property) {
	option := prop.(*notionapi.SelectProperty).Select
	o, _ := strconv.Atoi(option.Name)
	s.step.Points = o
}

func (s *stepBuilder) setPriority(prop notionapi.Property) {
	option := prop.(*notionapi.SelectProperty).Select
	s.step.Priority = Priority(option.Name)
}

func (s *stepBuilder) setType(prop notionapi.Property) {
	options := prop.(*notionapi.MultiSelectProperty).MultiSelect

	for _, option := range options {
		s.step.Type = append(s.step.Type, option.Name)
	}
}

func (s *stepBuilder) setStartAt(prop notionapi.Property) {
	date := prop.(*notionapi.DateProperty).Date
	if date != nil {
		s.step.StartedAt = time.Time(*date.Start)
	}
}

func (s *stepBuilder) setCreatedAt(prop notionapi.Property) {
	date := prop.(*notionapi.DateProperty).Date
	if date != nil {
		s.step.CreatedAt = time.Time(*date.Start)
	}
}

func (s *stepBuilder) setDeadline(prop notionapi.Property) {
	date := prop.(*notionapi.DateProperty).Date
	if date != nil {
		s.step.Deadline = time.Time(*date.Start)
	}
}
