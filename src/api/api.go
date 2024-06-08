package api

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"go.uber.org/zap"

	"oh-my-chat/src/config"
	"oh-my-chat/src/logger"
)

type OhMyChatApi struct {
	PreviewApi *PreviewApi
}

func NewOhMyChatApi() *OhMyChatApi {
	return &OhMyChatApi{
		PreviewApi: NewPreviewApi(),
	}
}

func RunApi(ctx context.Context, conf config.Api, handler http.Handler) {

	logging := logger.Logger.With(zap.String("context", "api_run"))
	server := &http.Server{
		Addr:    fmt.Sprintf("0.0.0.0:%d", conf.Port),
		Handler: handler,
	}

	servCtx, serverCtxCancel := context.WithCancel(context.Background())

	go func() {
		err := server.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			logging.Error("error serving API", zap.Error(err))
			serverCtxCancel()
		}
	}()

	go func() {
		<-ctx.Done()
		logging.Debug("Stopping API")

		shutdown, cancel := context.WithTimeout(servCtx, 30*time.Minute)
		err := server.Shutdown(shutdown)
		if err != nil {
			logging.Error("error shutdown API", zap.Error(err))
		}

		cancel()
		<-shutdown.Done()
		serverCtxCancel()

	}()

	<-servCtx.Done()
	logging.Debug("API server stopped")

}
