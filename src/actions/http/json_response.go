package http

import (
	"errors"

	"oh-my-chat/src/models"
	"oh-my-chat/src/utils"
)

type Summarizer func(response []byte, fields SummarizeFields, config SummarizeConfig) Summarized

type HttpJsonResponseHandler struct {
	summarizeConfig SummarizeConfig
	summarizeFields SummarizeFields
	summarize       Summarizer
}

func NewHttpJsonResponseHandler(config models.JsonResponseConfig) *HttpJsonResponseHandler {

	summarizeFields := SummarizeFields(
		utils.Map(config.Summarize.SummarizeFields, func(f models.SummarizeField) SummarizeField {
			return SummarizeField{Name: f.Name, Path: f.Path}
		}),
	)

	return &HttpJsonResponseHandler{
		summarizeConfig: SummarizeConfig{
			MaxInner:       config.Summarize.MaxInner,
			SeparatorStyle: Separator(config.Summarize.Separator),
		},
		summarizeFields: summarizeFields,
		summarize:       Summarize,
	}
}

func (h *HttpJsonResponseHandler) Handle(response []byte, message *models.Message) error {
	switch h.summarizeConfig.SeparatorStyle {
	case "write_space":
		h.summarizeConfig.SeparatorStyle = WriteSpaceStyle
	case "colon":
		h.summarizeConfig.SeparatorStyle = ColonStyle
	case "hyphen":
		h.summarizeConfig.SeparatorStyle = HyphenStyle
	case "pipe":
		h.summarizeConfig.SeparatorStyle = PipeStyle
	}

	if h.summarizeFields.IsEmpty() {
		return errors.New("summarizeFields must be set")
	}

	summarized := h.summarize(response, h.summarizeFields, h.summarizeConfig)

	if summarized.IsArray() {
		return errors.New("Not Implemented Error")
	}
	message.Output = summarized.String()
	return nil

}
