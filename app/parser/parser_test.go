package parser

import (
	"github.com/stretchr/testify/assert"
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

type TestCase struct {
	input    []string
	output    []string
}

const tests = []TestCase{
	{
		input:  "FOUNDER (START-UP)\n\n\t Nov, 2022 - May, 2023",
		output: {
			start: "2020 nov",
			end: "2023 may",
			position: ""
		},
	},
}

func Test_Match(t *testing.T) {
	// CV_FULL_PARAGRAPH_FORMATED
	Parse(CV_ONLY_EXPIRIENCE)
	for _, test := range tests {
		assert.Equal(test.input, test.output)
	}
}
