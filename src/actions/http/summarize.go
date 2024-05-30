package http

import (
	"github.com/tidwall/gjson"

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

func Summarize(response []byte, fields SummarizeFields, separator Separator) string {
	output := utils.NewStringBuilder()

	for _, field := range fields {
		var value string
		result := gjson.GetBytes(response, field.Path)

		if result.IsObject() {
			value = "ommitted"
			summarized := field.Name + string(separator) + value
			output.NextLine(summarized)
			continue
		}

		if result.IsArray() {
			max_inner := 10
			total_raw := len(result.Array())

			for index, raw := range result.Array() {
				if raw.IsObject() {
					value = "ommited"
					continue
				}
				if index == max_inner {
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
		value = "ommited"
	}
	switch result.Type {
	case gjson.String, gjson.Null, gjson.False, gjson.True, gjson.Number:
		value = result.String()
	default:
		value = "ommited"
	}
	return value
}
