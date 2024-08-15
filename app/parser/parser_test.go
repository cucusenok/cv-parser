package parser

import (
	"reflect"
	"strings"
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

func Test_Match(t *testing.T) {
	cvData := "FOUNDER (START-UP) Nov, 2022 - May, 2023"
	//cvData := "FULLSTACK DEVELOPER June, 2019 - June, 2020"
	cvData = strings.ToLower(cvData)
	Parse(cvData)
}
