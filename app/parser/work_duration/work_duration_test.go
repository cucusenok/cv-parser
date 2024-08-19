package work_duration

import (
	"awesomeProject2/app/parser"
	"reflect"
	"testing"
)

func Test_reformatPeriod(t *testing.T) {
	tests := []TestCase{
		{
			input: "",
			output: WorkPeriod{
				DateStart: "",
				DateEnd:   "",
			},
		},
		{
			input: "👉",
			output: WorkPeriod{
				DateStart: "",
				DateEnd:   "",
			},
		},
		{
			input: "123",
			output: WorkPeriod{
				DateStart: "",
				DateEnd:   "",
			},
		},
		{
			input: "2000",
			output: WorkPeriod{
				DateStart: "2000",
				DateEnd:   "+Inf",
			},
		},
		{
			input: "1923 2033",
			output: WorkPeriod{
				DateStart: "1923",
				DateEnd:   "2033",
			},
		},
		{
			input: "1923/2033",
			output: WorkPeriod{
				DateStart: "1923",
				DateEnd:   "2033",
			},
		},
		{
			input: "1923👉2033",
			output: WorkPeriod{
				DateStart: "1923",
				DateEnd:   "2033",
			},
		},
		{
			input: "1923+2033",
			output: WorkPeriod{
				DateStart: "1923",
				DateEnd:   "2033",
			},
		},
		{
			input: "1923.12",
			output: WorkPeriod{
				DateStart: "12.1923",
				DateEnd:   "+Inf",
			},
		},
		{
			input: "12.1923",
			output: WorkPeriod{
				DateStart: "12.1923",
				DateEnd:   "+Inf",
			},
		},
		{
			input: "11/12/1923👉12.2033",
			output: WorkPeriod{
				DateStart: "11.12.1923",
				DateEnd:   "12.1923",
			},
		},
		{
			input: "1923.12",
			output: WorkPeriod{
				DateStart: "12.1923",
				DateEnd:   "+Inf",
			},
		},
		{
			input: "12.1923",
			output: WorkPeriod{
				DateStart: "12.1923",
				DateEnd:   "+Inf",
			},
		},
		{
			input: "2.1923",
			output: WorkPeriod{
				DateStart: "02.1923",
				DateEnd:   "+Inf",
			},
		},
		{
			input: "12ю1923",
			output: WorkPeriod{
				DateStart: "12.1923",
				DateEnd:   "+Inf",
			},
		},
		//{
		//	input: "11/12/1923",
		//	output: WorkPeriod{
		//		DateStart: "11.12.1923",
		//		DateEnd:   "+Inf",
		//	},
		//},
		//{
		//	input: "1.09.1923",
		//	output: WorkPeriod{
		//		DateStart: "01.09.1923",
		//		DateEnd:   "+Inf",
		//	},
		//},
		//{
		//	input: "1.9.1923",
		//	output: WorkPeriod{
		//		DateStart: "01.09.1923",
		//		DateEnd:   "+Inf",
		//	},
		//},
	}
	for _, tt := range tests {
		t.Run("reformatPeriod", func(t *testing.T) {
			result := reformatPeriod(tt.input)
			if !reflect.DeepEqual(result, tt.output) {
				t.Errorf("got %v, want %v", result, tt.output)
			}
		})
	}
}

type TestCase struct {
	input  string
	output WorkPeriod
}

func Test_Match(t *testing.T) {
	// CV_FULL_PARAGRAPH_FORMATED
	parser.Parse(parser.CV_ONLY_EXPIRIENCE)
	//for _, test := range tests {
	//	assert.Equal(t, test.output, test.input)
	//}
}
