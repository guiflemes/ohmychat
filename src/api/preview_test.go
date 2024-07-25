package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"

	"oh-my-chat/src/logger"
	"oh-my-chat/src/models"
)

func TestMain(m *testing.M) {
	logger.InitLog("disable")

	m.Run()
}

func fakeHandlerPreview(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusAccepted)

	jsonData := []byte(`{
  "mainKey": "value",
  "nestedObject": {
    "key1": "nested value 1",
    "key2": 123.45,
    "list": [
      "item1",
      "item2",
    ],
  },
  "caca": [ 
    {
      "subKey1": "value1",
      "subKey2": false
    },
    {
      "subKey1": "value2",
      "subKey2": false
    },
    {
      "subKey1": "value3",
      "subKey2": false
    },
  ],
}`)
	w.Write(jsonData)
}

func TestJsonPreview(t *testing.T) {
	assert := assert.New(t)

	mockServer := httptest.NewServer(http.HandlerFunc(fakeHandlerPreview))
	defer mockServer.Close()

	type testCase struct {
		Desc           string
		Method         string
		QueryStr       string
		Model          models.HttpGetModel
		ExpectedStatus int
		ExpectedResult PreviewData
	}

	for _, _case := range []testCase{
		{
			Desc:     "get nested key",
			Method:   http.MethodPost,
			QueryStr: "action=get",
			Model: models.HttpGetModel{
				Url:     mockServer.URL,
				Headers: models.Headers{ContentType: "application/json"},
				TimeOut: 60,
				JsonResponseConfig: models.JsonResponseConfig{
					Summarize: models.Summarize{
						Separator: "colon",
						MaxInner:  10,
						SummarizeFields: []models.SummarizeField{
							{Name: "Teste1", Path: "nestedObject.key1"},
							{Name: "Teste2", Path: "mainKey"},
						},
					},
				},
			},
			ExpectedStatus: http.StatusOK,
			ExpectedResult: PreviewData{Fields: []string{"Teste1: nested value 1", "Teste2: value"}},
		},
		{
			Desc:     "get inner list values",
			Method:   http.MethodPost,
			QueryStr: "action=get",
			Model: models.HttpGetModel{
				Url:     mockServer.URL,
				Headers: models.Headers{ContentType: "application/json"},
				TimeOut: 60,
				JsonResponseConfig: models.JsonResponseConfig{
					Summarize: models.Summarize{
						Separator: "colon",
						MaxInner:  10,
						SummarizeFields: []models.SummarizeField{
							{Name: "Teste1", Path: "nestedObject.list"},
						},
					},
				},
			},
			ExpectedStatus: http.StatusOK,
			ExpectedResult: PreviewData{Fields: []string{"Teste1: item1, item2"}},
		},
		{
			Desc:     "omitted, it tries to access a struct into array without a specific key",
			Method:   http.MethodPost,
			QueryStr: "action=get",
			Model: models.HttpGetModel{
				Url:     mockServer.URL,
				Headers: models.Headers{ContentType: "application/json"},
				TimeOut: 60,
				JsonResponseConfig: models.JsonResponseConfig{
					Summarize: models.Summarize{
						Separator: "colon",
						MaxInner:  10,
						SummarizeFields: []models.SummarizeField{
							{Path: "caca"},
						},
					},
				},
			},
			ExpectedStatus: http.StatusOK,
			ExpectedResult: PreviewData{Fields: []string{"omitted"}},
		},
		{
			Desc:     "get value from key into array of structs",
			Method:   http.MethodPost,
			QueryStr: "action=get",
			Model: models.HttpGetModel{
				Url:     mockServer.URL,
				Headers: models.Headers{ContentType: "application/json"},
				TimeOut: 60,
				JsonResponseConfig: models.JsonResponseConfig{
					Summarize: models.Summarize{
						Separator: "colon",
						MaxInner:  2,
						SummarizeFields: []models.SummarizeField{
							{Path: "caca.#.subKey1"},
						},
					},
				},
			},
			ExpectedStatus: http.StatusOK,
			ExpectedResult: PreviewData{Fields: []string{"value1, value2, ..."}},
		},
	} {
		t.Run(_case.Desc, func(t *testing.T) {

			body, _ := json.Marshal(_case.Model)

			recorder := httptest.NewRecorder()
			req, err := http.NewRequest(
				_case.Method,
				fmt.Sprintf("/api/preview/json-response?%s", _case.QueryStr),
				bytes.NewReader(body),
			)

			assert.NoError(err)

			api := NewOhMyChatApi()

			api.PreviewApi.JsonResponse(recorder, req)

			resp := recorder.Result()

			assert.Equal(_case.ExpectedStatus, resp.StatusCode)

			exptectResult, _ := json.Marshal(_case.ExpectedResult)
			assert.Equal(string(exptectResult), recorder.Body.String())
		})
	}

}
