package utils

import (
	"fmt"
	"reflect"
	"strconv"
)

type TagTranspiler struct {
	Tags  []string
	index int
	next  bool
	key   string
}

func (t *TagTranspiler) Get() string {
	if len(t.Tags) == 0 {
		return ""
	}

	tag := t.Tags[t.index]
	t.index = (t.index + 1) % len(t.Tags)

	if t.index == 0 {
		t.next = false
	}

	if t.index == 1 {
		t.next = true
	}

	t.key = tag
	return tag
}

func (t *TagTranspiler) HasNext() bool { return t.next }

type TagsTranspiler []TagTranspiler

func FnTag(m map[string]any, tags TagsTranspiler) string {
	parser := ParseType{}
	result := ""

	for _, tag := range tags {
		value := getValue(m, &tag)
		str := parser.Parse(value)
		result += fmt.Sprintf("%s : %s/n", tag.key, str)
	}

	return result
}

func getValue(m map[string]any, tag *TagTranspiler) any {
	var r any
	for {
		key := tag.Get()

		if key == "" {
			break
		}

		value, ok := m[key]
		if !ok {
			break
		}

		fmt.Println(key, tag.HasNext())
		if !tag.HasNext() {
			r = value
			break
		}
		m = value.(map[string]any)
	}
	return r

}

type ParseType struct{}

func (p *ParseType) Parse(value any) string {
	t := reflect.TypeOf(value)

	if t.Kind() == reflect.Slice {
		return p.Slice(value.([]any), t)
	}
	return p.Primitive(value, t)
}

func (p *ParseType) Primitive(value any, t reflect.Type) string {

	switch t.Kind() {
	case reflect.Bool:
		return p.Boolean(value.(bool))

	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return p.Integer(value.(int))

	case reflect.Float32, reflect.Float64:
		return p.Float64(value.(float64))

	default:
		return value.(string)
	}

}

func (p *ParseType) Slice(value []any, t reflect.Type) string {
	v := ""
	for _, item := range value {
		str := p.Primitive(item, t.Elem())
		v += str + ", "
	}
	return v
}

func (p *ParseType) Integer(value int) string { return strconv.Itoa(value) }

func (p *ParseType) Float64(value float64) string { return strconv.FormatFloat(value, 'f', -1, 64) }

func (p *ParseType) Boolean(value bool) string { return strconv.FormatBool(value) }

// TODO SLICE, MAP, SLICE[MAP]
func TesteFn() {
	m := map[string]any{
		"string": "xixi",
		"map": map[string]any{
			"string": "xixi",
			"slice":  []int{1, 2, 3},
		},
		"bool":  true,
		"int":   123,
		"float": 123.2,
		"slice": []int{1, 2, 3},
	}
	tags := TagsTranspiler{
		//TagTranspiler{Tags: []string{"string"}},
		//TagTranspiler{Tags: []string{"bool"}},
		//TagTranspiler{Tags: []string{"int"}},
		//TagTranspiler{Tags: []string{"float"}},
		TagTranspiler{Tags: []string{"map", "string"}},
	}
	s := FnTag(m, tags)
	fmt.Println(s)
}
