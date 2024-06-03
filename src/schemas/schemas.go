package schemas

import (
	"encoding/json"
	"fmt"

	"oh-my-chat/src/actions/http"
	"oh-my-chat/src/core"
	"oh-my-chat/src/models"
)

func ParseGuidedSchemas(intents []models.IntentModel) (*core.MessageTree, error) {
	tree := &core.MessageTree{}

	for _, intent := range intents {
		for _, option := range intent.Options {

			a, err := decodeAction(option.Action)

			if err != nil {
				return nil, err
			}

			if tree.Root() == nil {
				tree.Insert(core.NewMessageNode(intent.Key, "", intent.Name, intent.Name, a))
			}

			node := core.NewMessageNode(option.Key, intent.Key, option.Name, option.Content, a)
			tree.Insert(node)
		}
	}

	return tree, nil
}

type decodeRawAction func(rawModel map[string]any) (core.Action, error)

var Actions = map[models.ModelType]decodeRawAction{
	models.TypeHttpGetModel: func(rawModel map[string]any) (core.Action, error) {

		model := &models.HttpGetModel{}
		if err := parseRawModel(model, rawModel); err != nil {
			return nil, err
		}

		fmt.Println(model.JsonResponseConfig.Summarize.MaxInner)

		return http.NewHttpGetAction(model), nil

	},
}

func decodeAction(model *models.ActionModel) (core.Action, error) {

	if model == nil {
		return nil, nil

	}

	builder, ok := Actions[models.ModelType(model.Type)]
	if !ok {
		return nil, fmt.Errorf("unsupported model type: %s", model.Type)
	}

	switch rawModel := model.Object.(type) {
	case map[string]any:
		action, err := builder(rawModel)
		if err != nil {
			return nil, err
		}
		return action, nil

	default:
		return nil, fmt.Errorf("unsupported model format %T", model.Object)
	}

}

func parseRawModel(model models.Model, rawModel map[string]any) error {
	b, err := json.Marshal(rawModel)

	if err != nil {
		return err
	}

	if err := json.Unmarshal(b, &model); err != nil {
		return err
	}

	return nil
}
