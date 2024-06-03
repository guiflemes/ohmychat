package models

type SummarizeField struct {
	Name string `yaml:"name" json:"name"`
	Path string `yaml:"path" json:"path"`
}

type Summarize struct {
	Separator       string           `yaml:"separator" json:"separator"`
	MaxInner        int              `yaml:"max_inner" json:"max_inner"`
	SummarizeFields []SummarizeField `yaml:"fields"    json:"fields"`
}

type JsonResponseConfig struct {
	ItemsPerMessage     int       `json:"items_per_message"    yaml:"items_per_message"`
	Summarize           Summarize `json:"summarize"            yaml:"summarize"`
	TruncationIndicator string    `json:"truncation_indicator" yaml:"truncation_indicator"`
}

type Headers struct {
	Authorization string `yaml:"authorization" json:"authorization"`
	ContentType   string `yaml:"content_type"  json:"content_type"`
}

type HttpGetModel struct {
	Url                string             `yaml:"url"                  json:"url"`
	Headers            Headers            `yaml:"headers"              json:"headers"`
	ResponseField      string             `yaml:"response_field"       json:"response_field"`
	TimeOut            int                `yaml:"timeout"              json:"timeout"`
	JsonResponseConfig JsonResponseConfig `yaml:"json_response_config" json:"json_response_config"`
}

func (h *HttpGetModel) GetType() string {
	return "http_get"
}
