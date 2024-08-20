package work_duration

import (
	"awesomeProject2/app/parser"
	"reflect"
	"testing"
)

func Test_reformatPeriod(t *testing.T) {
	tests := []TestCase{
		// Not date
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
		// Only years
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
		// Incorrect separator between dates
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
			input: "1923ю2033",
			output: WorkPeriod{
				DateStart: "1923",
				DateEnd:   "2033",
			},
		},
		// Without days
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
		// Incorrect separator between date
		{
			input: "11/12/1923",
			output: WorkPeriod{
				DateStart: "11.12.1923",
				DateEnd:   "+Inf",
			},
		},
		{
			input: "11/12/1923👉12.2033",
			output: WorkPeriod{
				DateStart: "11.12.1923",
				DateEnd:   "12.2033",
			},
		},
		{
			input: "11/12/1923 12/2033",
			output: WorkPeriod{
				DateStart: "11.12.1923",
				DateEnd:   "12.2033",
			},
		},
		{
			input: "11/12/1923_12/2033",
			output: WorkPeriod{
				DateStart: "11.12.1923",
				DateEnd:   "12.2033",
			},
		},
		{
			input: "1.9.1923_12.2.2033",
			output: WorkPeriod{
				DateStart: "01.09.1923",
				DateEnd:   "12.02.2033",
			},
		},

		// TODO обработать некорректные даты
		//// Incorrect date
		//{
		//	input: "29.1923",
		//	output: WorkPeriod{
		//		DateStart: "Invalid start date",
		//		DateEnd:   "",
		//	},
		//},
		//{
		//	input: "1.29.1923",
		//	output: WorkPeriod{
		//		DateStart: "Invalid start date",
		//		DateEnd:   "",
		//	},
		//},
		//{
		//	input: "33.12.1923",
		//	output: WorkPeriod{
		//		DateStart: "Invalid start date",
		//		DateEnd:   "",
		//	},
		//},
		//{
		//	input: "1.12.1923 1.29.1923",
		//	output: WorkPeriod{
		//		DateStart: "01.12.1923",
		//		DateEnd:   "Invalid end date",
		//	},
		//},

		// Reverse
		{
			input: "1923.9.1",
			output: WorkPeriod{
				DateStart: "01.09.1923",
				DateEnd:   "+Inf",
			},
		},
		//// TODO не понимаю, как разделить подобные строки
		//{
		//	input: "1923.9.1_1923.12.3",
		//	output: WorkPeriod{
		//		DateStart: "01.09.1923",
		//		DateEnd:   "03.12.1923",
		//	},
		//},
		//{
		//	input: "1923_9_1_1923_12_3",
		//	output: WorkPeriod{
		//		DateStart: "01.09.1923",
		//		DateEnd:   "03.12.1923",
		//	},
		//},
		//{
		//	input: "1923_9_1923_12_3",
		//	output: WorkPeriod{
		//		DateStart: "09.1923",
		//		DateEnd:   "03.12.1923",
		//	},
		//},
		///*
		//		TODO тут непонятно, в каком формате это должно быть
		//
		//			DateStart: "03.09.1923",
		//			DateEnd:   "09.1923",
		//	или
		//
		//			DateStart: "09.1923",
		//			DateEnd:   "03.09.1923",
		//
		//*/
		//
		//{
		//	input: "1923_9_3_9_1923",
		//	output: WorkPeriod{
		//		DateStart: "03.09.1923",
		//		DateEnd:   "09.1923",
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
