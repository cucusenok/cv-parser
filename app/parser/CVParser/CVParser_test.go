package CVParser

import (
	"cv-parser/parser"
	"reflect"
	"testing"
)

type TestCase struct {
	input  string
	output []ExperienceString
}

func Test_ParseCV(t *testing.T) {
	tests := []TestCase{
		{
			//input: parser.FULL_CV_YEGOR,
			//input:  parser.CV_FULL_PARAGRAPH_FORMATED,
			input:  parser.FULL_CV_VLADISLAV,
			output: []ExperienceString{},
		},
	}

	for _, tt := range tests {
		t.Run("ParseCV", func(t *testing.T) {
			result, err := ParseCV(tt.input)
			if err != nil {
				// Отобрази в тестах появление ошибки
				t.Errorf("got %v, want %v", result, tt.output)
			}
			if !reflect.DeepEqual(result, tt.output) {
				t.Errorf("got %v, want %v", result, tt.output)
			}
		})
	}
}
