package ujeebu

import (
	"net/url"
	"reflect"
	"testing"
)

type testStruct struct {
	Field1  string `json:"field1"`
	Field2  int    `json:"field2"`
	Field22 int    `json:"field22,omitempty"`
	Field3  bool   `json:"field3"`
	Field33 bool   `json:"field33,omitempty"`
	Field4  string `json:"field4"`
	Field44 string `json:"field44,omitempty"`
}

func TestStructToQueryParams(t *testing.T) {
	tests := []struct {
		name     string
		input    interface{}
		expected url.Values
	}{
		{
			name: "basic struct",
			input: testStruct{
				Field1:  "value1",
				Field2:  42,
				Field3:  true,
				Field44: "",
			},
			expected: url.Values{
				"field1": []string{"value1"},
				"field2": []string{"42"},
				"field3": []string{"true"},
				"field4": []string{""},
			},
		},
		{
			name: "empty struct",
			input: testStruct{
				Field1:  "",
				Field2:  0,
				Field22: 0,
				Field3:  false,
				Field4:  "",
				Field44: "",
			},
			expected: url.Values{
				"field1": []string{""},
				"field2": []string{"0"},
				"field3": []string{"false"},
				"field4": []string{""},
			},
		},
		{
			name: "partial values",
			input: testStruct{
				Field2:  100,
				Field22: 200,
				Field3:  true,
				Field33: true,
			},
			expected: url.Values{
				"field1":  []string{""},
				"field2":  []string{"100"},
				"field22": []string{"200"},
				"field3":  []string{"true"},
				"field33": []string{"true"},
				"field4":  []string{""},
			},
		},
		{
			name:     "nil interface",
			input:    nil,
			expected: url.Values{},
		},
		{
			name:     "struct with no fields",
			input:    struct{}{},
			expected: url.Values{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := structToQueryParams(tt.input)
			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("got %v, expected %v", result, tt.expected)
			}
		})
	}
}
