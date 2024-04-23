package schemas

import (
	"encoding/json"
	"fmt"
)

type Schemas map[string]Schema

type SchemaType string

type Schema interface {
	GetID() string
	GetType() SchemaType
}

func (p *Schemas) UnmarshalJSON(data []byte) error {
	var raw map[string]any
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}
	props, err := parseSchemas(raw)

	if err != nil {
		return err
	}

	*p = props
	return nil
}

func parseSchemas(raw map[string]any) (map[string]Schema, error) {
	result := make(map[string]Schema)

	for k, v := range raw {
		switch rawSchema := v.(type) {
		case map[string]any:
			p, err := decodeSchema(rawSchema)

			if err != nil {
				return nil, err
			}

			b, err := json.Marshal(rawSchema)
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

func decodeSchema(raw map[string]any) (Schema, error) {
	var p Schema
	switch SchemaType(raw["type"].(string)) {
	case SchemaTypeHttpGet:
		p = &HttpGetSchema{}
	default:
		return nil, fmt.Errorf("unsupported property type: %s", raw["type"].(string))
	}
	return p, nil
}

func parseGuidedEgineSchema() {}
