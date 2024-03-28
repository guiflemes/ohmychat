package actions

import (
	"log"

	"github.com/go-resty/resty/v2"

	"oh-my-chat/src/utils"
)

type (
	HttpClient interface{ R() HttpReq }

	HttpReq interface {
		SetHeader(header, value string) HttpReq
		Get(url string) (HttpResp, error)
	}

	HttpResp interface{ String() string }
)

type restyHttpClientAdapter struct {
	client *resty.Client
}

func (r *restyHttpClientAdapter) R() HttpReq {
	return &restyReqAdapter{req: r.client.R()}
}

type restyRespAdapter struct {
	resp *resty.Response
}

func (r *restyRespAdapter) String() string {
	return r.resp.String()
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

func NewHttpGetAction(url, auth string) *httpGetAction {

	client_adapter := &restyHttpClientAdapter{client: resty.New()}

	return &httpGetAction{
		url:    url,
		auth:   auth,
		client: client_adapter,
	}
}

type httpGetAction struct {
	url    string
	auth   string
	client HttpClient
	tags   utils.TagsTranspiler
}

func (a *httpGetAction) Execute(message string) string {
	req := a.client.R()

	if a.auth != "" {
		req.SetHeader("Authorization", a.auth)
	}

	resp, err := req.SetHeader("Content-Type", "application/json").
		SetHeader("Accept", "application/json").Get(a.url)

	if err != nil {
		log.Println("error httpGetAction", err)
		return "some error ocurred"
	}

	if a.tags == nil {
		return resp.String()
	}

	return ""
}
