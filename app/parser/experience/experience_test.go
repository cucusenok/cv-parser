package experience

import (
	"awesomeProject2/app/parser"
	"reflect"
	"strings"
	"testing"
)

type TestCase struct {
	input  string
	output []ExperienceString
}

func Test_parseExperience(t *testing.T) {
	tests := []TestCase{
		//// Without title
		//{
		//	input:  "Performed analysis and",
		//	output: []ExperienceString{},
		//},
		//// Only title
		//{
		//	input: "typescript reactjs programmer designer",
		//	output: []ExperienceString{
		//		{
		//			Date:        "",
		//			Title:       strings.ToLower("Typescript reactjs programmer designer"),
		//			Description: "",
		//		},
		//	},
		//},
		//
		//// Title and description
		//{
		//	input: "Typescript reactjs programmer designer \n Performed analysis and",
		//	output: []ExperienceString{{
		//		Date:        "",
		//		Title:       strings.ToLower("Typescript reactjs programmer designer"),
		//		Description: strings.ToLower("Performed analysis and"),
		//	}},
		//},
		//{
		//	input: "typescript reactjs programmer designer \n Performed analysis and \n",
		//	output: []ExperienceString{{
		//		Date:        "",
		//		Title:       strings.ToLower("Typescript reactjs programmer designer"),
		//		Description: strings.ToLower("Performed analysis and"),
		//	}},
		//},
		//// Title and description reverse
		//{
		//	input: "Performed analysis and \n typescript reactjs programmer designer",
		//	output: []ExperienceString{{
		//		Date:        "",
		//		Title:       strings.ToLower("Typescript reactjs programmer designer"),
		//		Description: strings.ToLower("Performed analysis and"),
		//	}},
		//},
		//{
		//	input: "typescript reactjs programmer designer \n Performed and",
		//	output: []ExperienceString{{
		//		Date:        "",
		//		Title:       strings.ToLower("Typescript reactjs programmer designer"),
		//		Description: strings.ToLower("Performed and"),
		//	}},
		//},
		//{
		//	input: "Typescript reactjs programmer designer \n Description 1 \n Description 2",
		//	output: []ExperienceString{{
		//		Date:        "",
		//		Title:       strings.ToLower("Typescript reactjs programmer designer"),
		//		Description: strings.ToLower("Description 1\nDescription 2"),
		//	}},
		//},
		//{
		//
		//	// Multi experience blocks
		//	input: "Typescript reactjs programmer designer 1 \n Description 1_1 \n Description 1_2 \n Typescript reactjs programmer designer 2 \n Description 2_1 \n Description 2_2",
		//	output: []ExperienceString{{
		//		Date:        "",
		//		Title:       strings.ToLower("Typescript reactjs programmer designer 1"),
		//		Description: strings.ToLower("Description 1_1\nDescription 1_2"),
		//	}, {
		//		Date:        "",
		//		Title:       strings.ToLower("Typescript reactjs programmer designer 2"),
		//		Description: strings.ToLower("Description 2_1\nDescription 2_2"),
		//	}},
		//},
		//
		//// With Dates
		//{
		//	input: "Typescript reactjs programmer designer \n Nov, 2022 \n Description 1",
		//	output: []ExperienceString{{
		//		Date:        strings.ToLower("Nov, 2022"),
		//		Title:       strings.ToLower("Typescript reactjs programmer designer"),
		//		Description: strings.ToLower("Description 1"),
		//	}},
		//},
		//{
		//	input: "Typescript reactjs programmer designer \n Nov, 2022 - May, 2023 \n Description 1",
		//	output: []ExperienceString{{
		//		Date:        strings.ToLower("Nov, 2022 - May, 2023"),
		//		Title:       strings.ToLower("Typescript reactjs programmer designer"),
		//		Description: strings.ToLower("Description 1"),
		//	}},
		//},
		//{
		//	input: "Typescript reactjs programmer designer 1 \n Nov, 2022 - May, 2023 \n Description 1 \n Typescript reactjs programmer designer 2 \n 2022 - 2023 \n Description 2_1",
		//	output: []ExperienceString{{
		//		Date:        strings.ToLower("Nov, 2022 - May, 2023"),
		//		Title:       strings.ToLower("Typescript reactjs programmer designer 1"),
		//		Description: strings.ToLower("Description 1"),
		//	}, {
		//		Date:        strings.ToLower("2022 - 2023"),
		//		Title:       strings.ToLower("Typescript reactjs programmer designer 2"),
		//		Description: strings.ToLower("Description 2_1"),
		//	}},
		//},
		//{
		//	input: "FULLSTACK DEVELOPER\n\nJune, 2019 - June, 2020\n\nsocial network for bikers Portfolio Product\n\n• Effectively allocated tasks and responsibilities between frontend and backend developers, ensuring efficient collaboration and workflow management\u0080\n\n• Streamlined and simplified the conversion process of work-flow documentation, descriptions, and processes, enhancing team productivity and understanding\u0080\n\nj\n\nq\n\n• Successfully updated the current pro ect state to align with modern re uirements and industry best practices\u0080\n\nQ\n\n• Developed a search engine utilizing Solr and MyS L database, facilitating efficient and accurate data retrieval and search functionalities.",
		//	output: []ExperienceString{{
		//		Date:        strings.ToLower("June, 2019 - June, 2020"),
		//		Title:       strings.ToLower("FULLSTACK DEVELOPER"),
		//		Description: strings.ToLower("social network for bikers Portfolio Product\n• Effectively allocated tasks and responsibilities between frontend and backend developers, ensuring efficient collaboration and workflow management\u0080\n• Streamlined and simplified the conversion process of work-flow documentation, descriptions, and processes, enhancing team productivity and understanding\u0080\nj\nq\n• Successfully updated the current pro ect state to align with modern re uirements and industry best practices\u0080\nQ\n• Developed a search engine utilizing Solr and MyS L database, facilitating efficient and accurate data retrieval and search functionalities"),
		//	}},
		//},
		{
			input: "FULLSTACK DEVELOPER June, 2019 - June, 2020\n\nsocial network for bikers Portfolio Product\n\n• Effectively allocated tasks and responsibilities between frontend and backend developers, ensuring efficient collaboration and workflow management\u0080\n\n• Streamlined and simplified the conversion process of work-flow documentation, descriptions, and processes, enhancing team productivity and understanding\u0080\n\nj\n\nq\n\n• Successfully updated the current pro ect state to align with modern re uirements and industry best practices\u0080\n\nQ\n\n• Developed a search engine utilizing Solr and MyS L database, facilitating efficient and accurate data retrieval and search functionalities.",
			output: []ExperienceString{{
				Date:        strings.ToLower("June, 2019 - June, 2020"),
				Title:       strings.ToLower("FULLSTACK DEVELOPER"),
				Description: strings.ToLower("social network for bikers Portfolio Product\n• Effectively allocated tasks and responsibilities between frontend and backend developers, ensuring efficient collaboration and workflow management\u0080\n• Streamlined and simplified the conversion process of work-flow documentation, descriptions, and processes, enhancing team productivity and understanding\u0080\nj\nq\n• Successfully updated the current pro ect state to align with modern re uirements and industry best practices\u0080\nQ\n• Developed a search engine utilizing Solr and MyS L database, facilitating efficient and accurate data retrieval and search functionalities"),
			}},
		},
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
