package models

type Headers struct {
	Authorization string `yaml:"authorization" json:"authorization"`
	ContentType   string `yaml:"content_type"  json:"content_type"`
}

type HttpGetModel struct {
	Url           string  `yaml:"url"            json:"url"`
	Headers       Headers `yaml:"headers"        json:"headers"`
	ResponseField string  `yaml:"response_field" json:"response_field"`
	TimeOut       int     `yaml:"timeout"        json:"timeout"`
}

func (h *HttpGetModel) GetType() string {
	return "http_get"
}
