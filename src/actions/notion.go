package actions

import (
	"context"
	"fmt"

	"oh-my-chat/settings"
	"oh-my-chat/src/notion"
	"oh-my-chat/src/notion/app"
	"oh-my-chat/src/notion/app/query"
	"oh-my-chat/src/utils"
)

type NotionCredentions struct {
	PageID string
}

type UserRepoConfig interface {
	GetNotionCredentions(userID string) NotionCredentions
}

type MemoryRepo struct{}

func (m *MemoryRepo) GetNotionCredentions(userID string) NotionCredentions {
	return NotionCredentions{
		PageID: settings.GETENV("PAGE_ID"),
	}
}

type notionActions struct {
	PedencyGetter *pedencyGetter
}

func NewNotionActions() *notionActions {
	repo := &notion.SketchRepo{}
	application := app.NotionApplication{
		Queries: app.Queries{
			PendencyTasks: *query.NewPendencyTaskHandler(repo),
		},
	}

	return &notionActions{
		PedencyGetter: &pedencyGetter{
			app:    application,
			config: &MemoryRepo{},
		},
	}
}

type pedencyGetter struct {
	app    app.NotionApplication
	config UserRepoConfig
}

func (p *pedencyGetter) Handler(ctx context.Context, userID string) string {
	credentions := p.config.GetNotionCredentions(userID)
	tasks, err := p.app.Queries.PendencyTasks.Handler(
		ctx,
		query.Pendendy{PageID: credentions.PageID},
	)
	if err != nil {
		return "Some error has ocurred"
	}

	if tasks == nil {
		return "No pendencies has found"
	}

	builder := utils.NewBulletListBuilder[notion.StudyStep]()
	return builder.Build(tasks, func(item notion.StudyStep) string {
		return fmt.Sprintf("[%s] - %s", item.Category, item.Name)
	})
}
