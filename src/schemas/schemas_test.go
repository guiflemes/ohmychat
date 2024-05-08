package schemas

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"oh-my-chat/src/actions"
	"oh-my-chat/src/core"
	"oh-my-chat/src/models"
)

var (
	httpGetModel *models.ActionModel = &models.ActionModel{
		Type: "http_get",
		Object: map[string]any{
			"url": "https://pokeapi.co/api/v2/pokemon/pikachu",
			"headers": map[string]string{
				"authorization": "",
				"content_type":  "application/json",
			},
			"response_field": "abilities[1].ability.name",
		},
	}
)

func TestDecodeAction(t *testing.T) {
	assert := assert.New(t)

	type testCase struct {
		desc               string
		model              *models.ActionModel
		expectedActionType core.Action
		expectedError      error
	}

	for _, c := range []testCase{
		{
			desc:               "decode HttpGetModel",
			model:              httpGetModel,
			expectedActionType: &actions.HttpGetAction{},
			expectedError:      nil,
		},
	} {
		t.Run(c.desc, func(t *testing.T) {})

		action, err := decodeAction(c.model)
		assert.Equal(c.expectedError, err)
		assert.IsType(action, c.expectedActionType)
	}
}
