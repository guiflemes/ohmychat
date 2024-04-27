package schemas

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"gopkg.in/yaml.v3"

	"oh-my-chat/src/actions"
	"oh-my-chat/src/core"
	"oh-my-chat/src/models"
)

type ModelTest struct {
	Engine  string               `yaml:"engine"`
	Intents []models.IntentModel `yaml:"intents"`
}

func ReadYml() {

	data, err := os.ReadFile("src/examples/guided_engine/pokemon.yml")
	if err != nil {
		log.Fatalf("error reading YAML file: %v", err)
	}

	var model ModelTest

	if err := yaml.Unmarshal(data, &model); err != nil {
		log.Fatalf("error Unmarshal YAML file: %v", err)
	}

	tree, err := parseGuidedSchemas(model.Intents)

	if err != nil {
		log.Fatalf("error parsing MessageTree: %v", err)
	}

	fmt.Println(tree.Search("chatoes").RepChildren())
}

func parseGuidedSchemas(intents []models.IntentModel) (*core.MessageTree, error) {
	tree := &core.MessageTree{}

	for _, intent := range intents {
		for _, option := range intent.Options {

			a, err := decodeAction(option.Action)

			if err != nil {
				return nil, err
			}

			if tree.Root() == nil {
				tree.Insert(core.NewMessageNode(intent.Key, "", intent.Name, "", a))
			}

			node := core.NewMessageNode(option.Key, intent.Key, option.Name, option.Content, a)
			tree.Insert(node)
		}
	}

	return tree, nil
}

var Actions = map[models.ModelType]func(rawModel map[string]any) (core.Action, error){
	models.TypeHttpGetModel: func(rawModel map[string]any) (core.Action, error) {

		b, err := json.Marshal(rawModel)
		if err != nil {
			return nil, err
		}

		var model models.HttpGetModel
		err = json.Unmarshal(b, &model)
		if err != nil {
			return nil, err
		}
		//TODO pass models arg

		return actions.NewHttpGetAction(model.Url, "", nil), nil

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
