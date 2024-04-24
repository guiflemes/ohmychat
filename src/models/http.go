package models

type Headers struct {
	Authorization string `json:"authorization"`
	ContentType   string `json:"content_type"`
}

type HttpGetModel struct {
	ID            string    `json:"id"`
	Type          ModelType `json:"type"`
	Url           string    `json:"url"`
	Headers       Headers   `json:"headers"`
	ResponseField string    `json:"response_field"`
}
