package utils

import (
	"encoding/json"
	"fmt"
	"reflect"
	"testing"

	"github.com/iv-p/mapaccess"
	"github.com/stretchr/testify/assert"
)

func TestParseTags(t *testing.T) {
	assert := assert.New(t)

	type testCase struct {
		desc     string
		value    any
		expected string
	}

	for _, c := range []testCase{
		{
			desc:     "parse int",
			value:    1,
			expected: "1",
		},
		{
			desc:     "parse true",
			value:    true,
			expected: "true",
		},
		{
			desc:     "parse false",
			value:    false,
			expected: "false",
		},
		{
			desc:     "parse float",
			value:    2.2,
			expected: "2.2",
		},
	} {
		t.Run(c.desc, func(t *testing.T) {
			r := Parse(c.value)
			assert.Equal(r, c.expected)
		})
	}
}

func TestLib(t *testing.T) {
	j := []byte(`{
    "id": "9b92b11b-b57f-4fa6-af5e-e35a290dc764",	
    "name": "John Doe",
    "friends": [
        {
            "name": "Jaime Mckinney"
        },
        {
            "name": "Evangeline Alvarado"
        },
        {
            "name": "Beth Cantrell"
        }
    ],
    "coco": [1, 2]
}`)
	var deserialised interface{}
	json.Unmarshal(j, &deserialised)

	see, _ := mapaccess.Get(deserialised, "coco")
	valueOfSlice := reflect.ValueOf(see)
	sliceType := valueOfSlice.Type()
	fmt.Println("ROAAAAAAAA", sliceType.Elem(), see)
	bestFriendName, _ := mapaccess.Get(deserialised, "friends[0].name")
	assert := assert.New(t)
	assert.Equal(bestFriendName, "Jaime Mckinney")
}
