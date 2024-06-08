package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"

	"oh-my-chat/src/models"
)

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
      {
        "subKey1": "sub value 1",
        "subKey2": false
      }
    ]
  }
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
			Desc:     "all ok",
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
