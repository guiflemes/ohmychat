package http

import (
	"errors"

	"github.com/tidwall/gjson"

	"oh-my-chat/src/config"
	"oh-my-chat/src/utils"
)

// Name is used to represent the name the user will see in the message
// Path is which value you want to access
type SummarizeField struct {
	Name string
	Path string
}

type SummarizeFields []SummarizeField

// Separator defines different styles of separators that can be used in the Summarize function.
// Predefined styles include:
//   - WriteSpaceStyle: " "
//   - ColonStyle: ": "
//   - SemmiColonStyle: "; "
//   - UnderscoreStyle: "_ "
//   - HyphenStyle: "- "
//   - PipeStyle: "| "
//
// You can also define a custom separator by creating a variable of type Separator and assigning it a custom value.
// Example:
//
//	var CustomSeparator Separator = " -> "
//	result := Summarize(response, fields, CustomSeparator)
type Separator string

const (
	WriteSpaceStyle Separator = " "
	ColonStyle                = ": "
	SemmiColonStyle           = "; "
	UnderscoreStyle           = "_ "
	HyphenStyle               = "- "
	PipeStyle                 = "| "
)

type SummarizeConfig struct {
	MaxInner       int
	SeparatorStyle Separator
}

type Summarized struct {
	value []string
}

func (s Summarized) IsArray() bool {
	return len(s.value) > 1
}

func (s Summarized) String() string {
	if !s.IsArray() {
		return s.value[0]
	}
	return ""
}

func (s Summarized) Stream() error {
	return errors.New("Not implemented error")
}

func Summarize(response []byte, fields SummarizeFields, config SummarizeConfig) Summarized {
	parsed := gjson.ParseBytes(response)
	if parsed.IsArray() {
		return Summarized{}
	}
	msg := summarize(parsed, fields, config)
	return Summarized{value: []string{msg}}
}

func summarize(response gjson.Result, fields SummarizeFields, summConfig SummarizeConfig) string {
	output := utils.NewStringBuilder()
	separator := summConfig.SeparatorStyle

	for _, field := range fields {
		var value string
		result := response.Get(field.Path)

		if result.IsObject() {
			value = config.MessageOmitted
			summarized := field.Name + string(separator) + value
			output.NextLine(summarized)
			continue
		}

		if result.IsArray() {
			total_raw := len(result.Array())

			for index, raw := range result.Array() {
				if raw.IsObject() {
					value = config.MessageOmitted
					continue
				}
				if index == summConfig.MaxInner {
					value += "..."
					break
				}
				innerValue := innerSummarize(raw)
				value += innerValue
				if index+1 != total_raw {
					value += ", "
				}
			}

			summarized := field.Name + string(separator) + value
			output.NextLine(summarized)
			continue
		}

		value = innerSummarize(result)
		summarized := field.Name + string(separator) + value
		output.NextLine(summarized)
	}

	return output.String()

}

func innerSummarize(result gjson.Result) string {
	var value string
	if result.IsObject() {
		value = config.MessageOmitted
	}
	switch result.Type {
	case gjson.String, gjson.Null, gjson.False, gjson.True, gjson.Number:
		value = result.String()
	default:
		value = config.MessageOmitted
	}
	return value
}
