package http

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/go-resty/resty/v2"
	"go.uber.org/zap"

	"oh-my-chat/src/logger"
	"oh-my-chat/src/models"
)

type (
	HttpClient interface {
		R() HttpReq
		SetTimeOut(timeout time.Duration)
	}

	HttpReq interface {
		SetHeader(header, value string) HttpReq
		Get(url string) (HttpResp, error)
	}

	HttpResp interface {
		String() string
		Body() []byte
	}
)

type restyHttpClientAdapter struct {
	client *resty.Client
}

func (r *restyHttpClientAdapter) R() HttpReq {
	return &restyReqAdapter{req: r.client.R()}
}

func (r *restyHttpClientAdapter) SetTimeOut(timeout time.Duration) {
	r.client.SetTimeout(timeout)
}

type restyRespAdapter struct {
	resp *resty.Response
}

func (r *restyRespAdapter) String() string {
	return r.resp.String()
}

func (r *restyRespAdapter) Body() []byte {
	return r.resp.Body()
}

type restyReqAdapter struct {
	req *resty.Request
}

func (r *restyReqAdapter) SetHeader(header, value string) HttpReq {
	r.req.SetHeader(header, value)
	return r
}

func (r *restyReqAdapter) Get(url string) (HttpResp, error) {
	resp, err := r.req.Get(url)
	if err != nil {
		return nil, err
	}
	return &restyRespAdapter{resp: resp}, nil
}

var logging = logger.Logger.With(zap.String("action", "get_http"))

func NewHttpGetAction(model *models.HttpGetModel) *HttpGetAction {
	client_adapter := &restyHttpClientAdapter{client: resty.New()}
	jsonResponseHandler := NewHttpJsonResponseHandler(model.JsonResponseConfig)

	return &HttpGetAction{
		url:          model.Url,
		auth:         model.Headers.Authorization,
		client:       client_adapter,
		contentType:  model.Headers.ContentType,
		timeOut:      model.TimeOut,
		jsonResponse: jsonResponseHandler,
	}
}

type HttpGetAction struct {
	url          string
	auth         string
	client       HttpClient
	contentType  string
	timeOut      int
	jsonResponse *HttpJsonResponseHandler
}

func (a *HttpGetAction) Handle(ctx context.Context, message *models.Message) error {
	//create something to replace values in url for exempe www.test.com/{invoice_id}  -> message.Input = "invoice_id"
	logging.With(
		zap.String("url", a.url),
		zap.String("provider", string(message.Connector)),
	)

	if a.timeOut == 0 {
		a.timeOut = 10
	}

	a.client.SetTimeOut(time.Duration(a.timeOut) * time.Second)

	req := a.client.R()

	if a.auth != "" {
		req.SetHeader("Authorization", a.auth)
	}

	switch a.contentType {
	case "application/json":
		return a.handleJson(req, message)
	case "":
		logging.Sugar().
			Error("contentType cannot be empty")
		return errors.New("contentType cannot be empty")
	default:
		logging.Sugar().
			Errorf("contentType '%s' is not supported currently", a.contentType)
		return fmt.Errorf("contentType '%s' is not supported currently", a.contentType)
	}

}

func (a *HttpGetAction) handleJson(req HttpReq, message *models.Message) error {

	resp, err := req.SetHeader("Content-Type", "application/json").
		SetHeader("Accept", "application/json").Get(a.url)

	if err != nil {
		logging.Error("Failed to fetch url", zap.Error(err))
		return errors.New("Failed to fetch url")
	}

	return a.jsonResponse.Handle(resp.Body(), message)

}
