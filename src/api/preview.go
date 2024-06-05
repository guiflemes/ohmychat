package api

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"go.uber.org/zap"

	ahttp "oh-my-chat/src/actions/http"
	"oh-my-chat/src/logger"
	"oh-my-chat/src/models"
)

var previewLog = logger.Logger.With(zap.String("context", "api_preview"))

type GetHandlerFactory func(model *models.HttpGetModel) *ahttp.HttpGetAction

type PreviewApi struct {
	getHandlerFactory GetHandlerFactory
}

func NewPreviewApi() *PreviewApi {
	return &PreviewApi{
		getHandlerFactory: ahttp.NewHttpGetAction,
	}
}

func (p *PreviewApi) JsonResponse(w http.ResponseWriter, r *http.Request) {
	action := r.URL.Query().Get("action")

	w.Header().Set("content-type", "application/json")

	if action == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"error": "action cannot be blank"`))
		return
	}

	decoder := json.NewDecoder(r.Body)

	message := models.NewMessage()
	ctx := context.TODO()

	switch action {
	case "get":
		var getModel = &models.HttpGetModel{}
		err := decoder.Decode(getModel)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(`{"error": "Error decoding JSON`))
			previewLog.Error("decoding JSON", zap.Error(err))
			return
		}

		handler := p.getHandlerFactory(getModel)
		err = handler.Handle(ctx, &message)

		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(`{"error": "Error handling message`))
			return
		}
		w.Write([]byte(fmt.Sprintf(`{"preview": "%s"}`, message.Output)))
		return
	case "post":
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"error": "not implemented yet"`))
	default:
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"error": "action must be 'post' or 'get'"`))
	}

}
