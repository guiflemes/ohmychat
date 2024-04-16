package models

import (
	"encoding/json"
	"fmt"
)

type Properties map[string]Property

type PropertyType string

type Property interface {
	GetID() string
	GetType() PropertyType
}

func (p *Properties) UnmarshalJSON(data []byte) error {
	var raw map[string]any
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}
	props, err := parseProperties(raw)

	if err != nil {
		return err
	}

	*p = props
	return nil
}

func parseProperties(raw map[string]any) (map[string]Property, error) {
	result := make(map[string]Property)

	for k, v := range raw {
		switch rawProperty := v.(type) {
		case map[string]any:
			p, err := decodeProperty(rawProperty)

			if err != nil {
				return nil, err
			}

			b, err := json.Marshal(rawProperty)
			if err != nil {
				return nil, err
			}

			err = json.Unmarshal(b, &p)
			if err != nil {
				return nil, err
			}
			result[k] = p
		default:
			return nil, fmt.Errorf("unsupported property format %T", v)
		}
	}

	return result, nil
}

func decodeProperty(raw map[string]any) (Property, error) {
	var p Property
	switch PropertyType(raw["type"].(string)) {
	case PropertyTypeHttpGet:
		p = &HttpGetProperty{}
	default:
		return nil, fmt.Errorf("unsupported property type: %s", raw["type"].(string))
	}
	return p, nil
}
