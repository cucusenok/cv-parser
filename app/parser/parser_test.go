package parser

import (
	"reflect"
	"testing"
)

func Test_GenerateCombinations(t *testing.T) {
	tests := []struct {
		input    []string
		expected []string
	}{
		{
			input:    []string{"1", "2", "3", "4", "5"},
			expected: []string{"1", "2", "3", "4", "5", "1 2", "1 2 3", "2 3", "2 3 4", "3 4", "3 4 5", "4 5"},
		},
		{
			input:    []string{"a", "b", "c"},
			expected: []string{"a", "b", "c", "a b", "a b c", "b c"},
		},
		{
			input:    []string{"x", "y"},
			expected: []string{"x", "y", "x y"},
		},
		{
			input:    []string{"onlyone"},
			expected: []string{"onlyone"},
		},
		{
			input:    []string{},
			expected: []string{},
		},
	}

	for _, tt := range tests {
		t.Run("generateStrings", func(t *testing.T) {
			result := GenerateCombinations(tt.input)
			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("got %v, want %v", result, tt.expected)
			}
		})
	}

}

type OutputData struct {
	Start    string `json:"start"`
	End      string `json:"end"`
	Position string `json:"position"`
}

type TestCase struct {
	input  string
	output []OutputData
}

func Test_Match(t *testing.T) {
	// CV_FULL_PARAGRAPH_FORMATED
	Parse(CV_ONLY_EXPIRIENCE)
	//for _, test := range tests {
	//	assert.Equal(t, test.output, test.input)
	//}
}
