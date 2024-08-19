package experience

import (
	"awesomeProject2/app/parser"
	"reflect"
	"testing"
)

func Test_GenerateCombinations(t *testing.T) {
	tests := []TestCase{
		{
			input: "FOUNDER (START-UP)\n\n\t Nov, 2022 - May, 2023",
			output: []OutputData{
				{
					Start:    "2020 nov",
					End:      "2023 may",
					Position: "",
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run("generateStrings", func(t *testing.T) {
			result := parser.GenerateCombinations(tt.input)
			if !reflect.DeepEqual(result, tt.output) {
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

var tests = []TestCase{
	{
		input: "FOUNDER (START-UP)\n\n\t Nov, 2022 - May, 2023",
		output: []OutputData{
			{
				Start:    "2020 nov",
				End:      "2023 may",
				Position: "",
			},
		},
	},
}

func Test_Match(t *testing.T) {
	// CV_FULL_PARAGRAPH_FORMATED
	parser.Parse(parser.CV_ONLY_EXPIRIENCE)
	//for _, test := range tests {
	//	assert.Equal(t, test.output, test.input)
	//}
}
