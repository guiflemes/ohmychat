package actions

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"github.com/go-resty/resty/v2"
	"github.com/iv-p/mapaccess"
	"go.uber.org/zap"

	"oh-my-chat/src/logger"
	"oh-my-chat/src/models"
	"oh-my-chat/src/utils"
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

func NewHttpGetAction(model *models.HttpGetModel) *HttpGetAction {
	client_adapter := &restyHttpClientAdapter{client: resty.New()}

	return &HttpGetAction{
		url:         model.Url,
		auth:        model.Headers.Authorization,
		client:      client_adapter,
		tag:         &TagAcess{model.ResponseField},
		mapAccess:   mapaccess.Get,
		contentType: model.Headers.ContentType,
		timeOut:     model.TimeOut,
	}
}

type TagAcess struct {
	Key string
}

type SomeError struct{}

func (e *SomeError) Error() string {
	return "some error has ocurred"
}

type MapAcesss func(data interface{}, key string) (interface{}, error)

type HttpGetAction struct {
	url         string
	auth        string
	client      HttpClient
	tag         *TagAcess
	mapAccess   MapAcesss
	contentType string
	timeOut     int
}

func (a *HttpGetAction) Handle(ctx context.Context, message *models.Message) error {
	//create something to replace values in url for exempe www.test.com/{invoice_id}  -> message.Input = "invoice_id"
	log := logger.Logger.With(
		zap.String("action", "get_http"),
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
		return a.handleJson(req, message, log)
	default:
		log.Sugar().
			Errorf("contentType '%s' is not supported currently", a.contentType)
		return &SomeError{}
	}

}

func (a *HttpGetAction) handleJson(req HttpReq, message *models.Message, log *zap.Logger) error {

	resp, err := req.SetHeader("Content-Type", "application/json").
		SetHeader("Accept", "application/json").Get(a.url)

	if err != nil {
		log.Error("Failed to fetch url", zap.Error(err))
		return &SomeError{}
	}
	if a.tag == nil {
		message.Output = resp.String()
		return nil
	}

	var deserialised interface{}
	err = json.Unmarshal(resp.Body(), &deserialised)

	if err != nil {
		log.Error("Failed to Unmarshal response", zap.Error(err))
		return &SomeError{}
	}

	if a.mapAccess == nil {
		log.Error("MapAcesss is nill", zap.Error(errors.New("MapAccess should not be nil")))
		return &SomeError{}
	}

	value, err := a.mapAccess(deserialised, a.tag.Key)
	if err != nil {
		log.Error("Failed to access key", zap.Error(err))
		return &SomeError{}
	}

	log.Info("action executed sucessfully")
	message.Output = utils.Parse(value)
	return nil
}
