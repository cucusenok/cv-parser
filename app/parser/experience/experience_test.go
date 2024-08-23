package experience

import (
	"awesomeProject2/app/parser"
	"reflect"
	"testing"
)

/*
Experience example:

"experience": [
   {
	 "date_start": "06-2019",
	 "date_end": "06-2020",
	 "place": "google",
	 "title": "SENIOR FULL STACK ENGINEER"
	 "position": ["ENGINEER"],
	 "skills": ["FULL STACK"],
	 "level": ["senior"],
	 "description": "Performed analysis and problem-solving for 3GPP, 4G/LTE, 5G, and O-RAN..."
   }
],
*/

type TestCase struct {
	input  string
	output []Experience
}

func Test_parseExperience(t *testing.T) {
	tests := []TestCase{
		{
			input: "typescript reactjs programmer designer",
			output: []Experience{
				{
					Start:       "",
					End:         "",
					Title:       "typescript reactjs programmer designer",
					Description: "",
				},
			},
		},
		{
			input: "typescript reactjs programmer designer \n Performed analysis and",
			output: []Experience{{
				Start:       "",
				End:         "",
				Title:       "typescript reactjs programmer designer",
				Description: "performed analysis and.",
			}},
		},
		{
			input: "typescript reactjs programmer designer \n Performed analysis and \n",
			output: []Experience{{
				Start:       "",
				End:         "",
				Title:       "typescript reactjs programmer designer",
				Description: "performed analysis and.",
			}},
		},
		{
			input: "Performed analysis and \n typescript reactjs programmer designer",
			output: []Experience{{
				Start:       "",
				End:         "",
				Title:       "typescript reactjs programmer designer",
				Description: "performed analysis and.",
			}},
		},
		{
			input: "typescript reactjs programmer designer \n Performed and",
			output: []Experience{{
				Start:       "",
				End:         "",
				Title:       "typescript reactjs programmer designer",
				Description: "performed and.",
			}},
		},
		{
			input: "typescript reactjs programmer designer \n Description 1 \n Description 2",
			output: []Experience{{
				Start:       "",
				End:         "",
				Title:       "typescript reactjs programmer designer",
				Description: "description 1.\ndescription 2.",
			}},
		},
		//{
		//	input: "FOUNDER (START-UP)\n\nNov, 2022 - May, 2023",
		//	output: Experience{
		//		Start:       "Nov, 2022",
		//		End:         "May, 2023",
		//		Title:       "FOUNDER (START-UP)",
		//		Description: "",
		//	},
		//},
	}

	for _, tt := range tests {
		t.Run("parseExperience", func(t *testing.T) {
			result, err := ParseExperience(tt.input)
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

type OutputData struct {
	Start    string `json:"start"`
	End      string `json:"end"`
	Position string `json:"position"`
}

func Test_Match(t *testing.T) {
	// CV_FULL_PARAGRAPH_FORMATED
	parser.Parse(parser.CV_ONLY_EXPIRIENCE)
	//for _, test := range tests {
	//	assert.Equal(t, test.output, test.input)
	//}
}
