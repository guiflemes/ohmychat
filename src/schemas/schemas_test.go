package schemas

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"oh-my-chat/src/actions/http"
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
	unsupportedModelType *models.ActionModel = &models.ActionModel{
		Type: "unsupportedModelType",
		Object: map[string]any{
			"error": "error",
		},
	}

	unsupportedModelFmt *models.ActionModel = &models.ActionModel{
		Type: "unsupportedModeFmt",
		Object: map[string]any{
			"error": "error",
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
			expectedActionType: &http.HttpGetAction{},
			expectedError:      nil,
		},
	} {
		t.Run(c.desc, func(t *testing.T) {})

		action, err := decodeAction(c.model)
		assert.Equal(c.expectedError, err)
		assert.IsType(action, c.expectedActionType)
	}
}

func TestDecodeActionError(t *testing.T) {
	assert := assert.New(t)

	action, err := decodeAction(unsupportedModelType)
	assert.ErrorContains(err, "unsupported model type")
	assert.Nil(action)

}
