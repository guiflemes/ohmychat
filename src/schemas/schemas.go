package schemas

import (
	"encoding/json"
	"errors"
	"fmt"

	"oh-my-chat/src/actions"
	"oh-my-chat/src/core"
	"oh-my-chat/src/models"
)

type Schemas map[string]Schema

type SchemaType string

type Schema interface {
	GetID() string
	GetType() SchemaType
}

func (p *Schemas) UnmarshalJSON(data []byte) error {
	var raw map[string]any
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}
	props, err := parseSchemas(raw)

	if err != nil {
		return err
	}

	*p = props
	return nil
}

// type Action interface {Handle(ctx context.Context, message *Message)}
//           engine
//  guided            rule
//  options           intent
//  action            action

func parseGuidedSchemas(intents []Intent) (map[string]Schema, error) {
	result := make(map[string]Schema)

	for _, intent := range intents {
		for _, option := range intent.Options {
			if &option.Action == nil {
				continue
			}

			a, err := decodeAction(option.Action)

			if err != nil {
				return nil, err
			}

		}
	}

	return result, nil
}

var Actions = map[models.ModelType]func(action models.ActionModel) (core.Action, bool){
	models.TypeHttpGetModel: func(action models.ActionModel) (core.Action, bool) {
		model, ok := action.Object.(models.HttpGetModel)
		if !ok {
			return nil, ok
		}
		//TODO pass models arg
		return actions.NewHttpGetAction(model.Url, "", nil), ok

	},
}

func decodeAction(model models.ActionModel) (core.Action, error) {
	builder, ok := Actions[models.ModelType(model.Type)]

	if !ok {
		return nil, fmt.Errorf("unsupported model type: %s", model.Type)
	}

	action, ok := builder(model)

	if !ok {
		return nil, fmt.Errorf("unsupported type swicth: %s", model.Type)
	}

	return action, nil
}
