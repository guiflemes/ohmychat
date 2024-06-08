package api

import "net/http"

func NewHttpMux(ohMyChat *OhMyChatApi) http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("POST api/preview/json-response", ohMyChat.PreviewApi.JsonResponse)

	return mux
}
