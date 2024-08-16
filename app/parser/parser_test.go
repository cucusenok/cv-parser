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
	//cvData := "FOUNDER (START-UP) Nov, 2022 - May, 2023"
	//cvData := "FULLSTACK DEVELOPER June, 2019 - June, 2020"
	cvData := " Tereshkin Oleksandr\nSenior Full Stack Developer 7+ YEARS OF EXPERIENCE\ncucusenok.work@gmail.com linkedin.com/in/cucusenok github.com/Cucusenok @cucusenok Portfolio\nWith a passion for programming ignited at the age of 14, I have embarked on a journey of continuous growth and achievement. Since my first year of university, I have been actively involved in the professional realm, gaining invaluable experience along the way. Skilled in project management and diverse programming languages including TypeScript Go, Python.\nAdditionally, I co-founded a startup venture with a group of talented colleagues. Regrettably, unforeseen circumstances, including a war, led to the dissolution of the team. Nevertheless, this experience taught me valuable lessons about resilience, adaptability, and the importance of teamwork.\nSeeking a new opportunity to contribute to something incredible.\nFRONTEND\nTypeScript, JavaScript, React,\nVue.js, Vuex, redux-toolkit, MobX, WebRTC, Fabric, D3, ApexChart, HTML, atomic, CSS, SCSS, jest, ajv, MUI, webpack, css-modules, storybook\nEDUCATION\n2016 - 2020, BSc with Honours, KNU\nInformation computing technology\n2020 - 2022, MSc, KNU Informationcomputingtechnology\nBACKEND\nNode.js, Python, FastApi, Flask, Golang, Mongo, SQL, Postgres, MVC\nENVIRONMENT & COMMON Docker, VirualBox, HTTP, WebRTC, RPC, WebSocket, REST, OOP, SOLID, DRY, KISS\ncertificates\nCisco Networking Academy\nMANAGE & DESIGN\nFigma, Notion, Waterfall, Agile, SCRUM, Jira, Confluence, BitBacket, YouTrack, Trello, \nMiro, AdobeXD\n; Cisco,IntroductiontoNetworksBacheloN ; Cisco,ProgrammingEssentialsinPythoD ; Cisco,RoutingandSwitchingEssentialR ; Cisco,IntroductiontoCybersecurity\nexperience\nFOUNDER (START-UP) Nov, 2022 - May, 2023\n' Developed music recognition and sharing service across diverse sources1\n' Built a robust CI/CD pipeline encompassing backend services, IOS app deployment to TestFlight, Android app packaging to\nAPK, and web deployment using GitHub Actions and AWS (EC2, S3)1\n' Developed a cross-platform application using Flutter, ensuring seamless functionality on Android, IOS, and web platforms1\n' Implemented audio decoding and playback features for streaming services on Android and IOS, utilizing proprietary Flutter-\nplugins and programming languages such as Java/Kotlin and Swift1\n' Leveraged API integrations to enhance application functionality and connectivity, while employing GoLang for the central\nmonolith and Python and Node for microservices.\nSENIOR FULL STACK ENGINEER\n5g Networks\nNov, 2020 - Nov, 2022\n' Performed analysis and problem-solving for 3GPP, 4G/LTE, 5G, and O-RAN architectures1\n' Designed and developed a user-friendly interface, ensuring a comfortable and intuitive experience for users1\n' Led system design efforts, collaborating with a team of dozens of engineers to architect robust and scalable solutions1 ' Developed an automated interface generation system using YANG models, enabling efficient DSL-based interfaces1\n' Implemented a DSL-based interface with full-text search and XML-tree validation/generation capabilities.\nMiddleware Analysis Platform\nPortfolio\n' Designed and developed interfaces utilizing micro-frontend architecture1\n' Designed visually engaging interfaces capable of effectively visualizing and presenting large volumes of data1\n' Developed Python-based microservices to efficiently store and visualize data, ensuring seamless functionality1\n' Implemented robust logging using Zipkin for query tracking, while leveraging Logstash and Elastic for log storage1\n' Utilized 'out-of-the-box' solutions to design DSL and interfaces for hierarchical control system rights, employing RBAC (Role-\nBased Access Control) at the application level to ensure secure access across multiple products.\n\n SENIOR FRONTEND ENGINEER Freeze - online viewing room\nJun, 2020 - Nov, 2020\ncucusenok.work@gmail.com\nFGP linkedin.com/in/cucusenok\nPortfolio Product\ngithub.com/Cucusenok\n\u0083 Led and trained a team of engineers, overseeing task control, sprint preparation, and maintenance processes\u0080\n\u0083 Successfully finalized and redesigned both the mobile application and web version, enhancing user experi@encuecaunsdenok functionality\u0080\n\u0083 Implemented OAuth and two-step verification via SMS, ensuring secure access and integration with RESTful API. Neo4Web\nPortfolio\n\u0083 Developed a real-time streaming application specifically for Unreal Engine\u0080\n\u0083 Collaborated as part of a 10-member team, including Unreal Engine (UE) and backend developers\u0080\n\u0083 Designed and supported interface integration for WebRTC across various browsers and operating systems\u0080\n\u0083 Conducted thorough testing and bug fixing across different environments using remote machines, Virtual Boxes, and\nphysical devices\u0080\n\u0083 Implemented command and data exchange functionality between the client and the Unreal Engine application.\nFULLSTACK DEVELOPER June, 2019 - June, 2020\nBike.net: social network for bikers\n\u0083 Effectively allocated tasks and responsibilities between frontend and backend developers, ensuring efficient collaboration\nand workflow management\u0080\n\u0083 Streamlined and simplified the conversion process of work-flow documentation, descriptions, and processes, enhancing\nteam productivity and understanding\u0080\n\u0083 Successfully updated the current project state to align with modern requirements and industry best practices\u0080\n\u0083 Developed a search engine utilizing Solr and MySQL database, facilitating efficient and accurate data retrieval and search\nfunctionalities.\nTinder-like Game for Durex Portfolio\n\u0083 Conducted comprehensive revisions of an existing project, addressing issues, adding new features, and enhancing overall functionality\u0080\n\u0083 Performed code proofreading, debugging, and optimization to improve the project's stability and performance\u0080\n\u0083 Successfully implemented a push notification system, enabling real-time notifications and enhancing user engagement\u0080\n\u0083 Developed an admin panel to facilitate content control and management, providing an intuitive interface for efficient\nadministration.\nJUNIOR FULL STACK DEVELOPER Jan, 2017 - June, 2019 Private practice / freelance\n\u0083 Designed and developed a system utilizing a messenger bot and online schedule to facilitate audience discovery for the internal university website\u0080\n\u0083 Developed an application for visualizing and analyzing Apache server logs, providing valuable insights into server performance and user behavior\u0080\n\u0083 Created a cloud file manager with robust file-sharing capabilities, incorporating RBAC (Role-Based Access Control) for secure and controlled file access\u0080\n\u0083 Successfully completed an Internship at the Institute of Artificial Intelligence of Ukraine, gaining hands-on experience in various aspects of AI, including text classification and fundamental methodologies\u0080\n\u0083 Studied and implemented key techniques such as tokenization, vectorization, lemmatization, logistic regression, and Random Forest in the context of natural language processing and machine learning.\ncucusenok.work@gmail.com linkedin.com/in/cucusenok github.com/Cucusenok @cucusenok Portfolio\nPortfolio Product\n"
	cvData = strings.ToLower(cvData)
	Parse(cvData)
}
