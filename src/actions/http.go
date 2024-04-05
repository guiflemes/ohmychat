package actions

import (
	"encoding/json"
	"log"

	"github.com/go-resty/resty/v2"
	"github.com/iv-p/mapaccess"

	"oh-my-chat/src/utils"
)

type (
	HttpClient interface{ R() HttpReq }

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

func NewHttpGetAction(url, auth string, tag *TagAcess) *httpGetAction {

	client_adapter := &restyHttpClientAdapter{client: resty.New()}

	return &httpGetAction{
		url:       url,
		auth:      auth,
		client:    client_adapter,
		tag:       tag,
		mapAccess: mapaccess.Get,
	}
}

type TagAcess struct {
	Key string
}

type httpGetAction struct {
	url       string
	auth      string
	client    HttpClient
	tag       *TagAcess
	mapAccess MapAcesss
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

	if a.tag == nil {
		return resp.String()
	}

	var deserialised interface{}
	err = json.Unmarshal(resp.Body(), &deserialised)

	if err != nil {
		log.Println("Unmarshal: ", err)
		return "some error ocurred"
	}

	if a.mapAccess == nil {
		log.Println("map acess is not set")
		return "some error ocurred"
	}

	value, err := a.mapAccess(deserialised, a.tag.Key)
	if err != nil {
		log.Println("MapAcesss: ", err)
		return "some error ocurred"
	}

	return utils.Parse(value)
}

type MapAcesss func(data interface{}, key string) (interface{}, error)
