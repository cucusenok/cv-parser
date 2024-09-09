package parser

import (
	"encoding/json"
	"fmt"
	"testing"
)

type Contact struct {
	Emails         []string          `json:"emails"`
	GitHub         string            `json:"github"`
	Address        []Address         `json:"address"`
	Phones         []string          `json:"phones"`
	SocialNetworks map[string]string `json:"social_networks"`
}

type Address struct {
	Country     string `json:"country"`
	City        string `json:"city"`
	State       string `json:"state"`
	CountryCode string `json:"country_code"`
	StateCode   string `json:"state_code"`
}

type Skill struct {
	Alias string `json:"alias"`
	Type  string `json:"type"`
	ID    int    `json:"id"`
}

type Experience struct {
	DateStart   string   `json:"date_start"`
	DateEnd     string   `json:"date_end"`
	Place       string   `json:"place"`
	Title       string   `json:"title"`
	Position    []string `json:"position"`
	Skills      []string `json:"skills"`
	Level       []string `json:"level"`
	Description string   `json:"description"`
}

type Education struct {
	DateStart   string   `json:"date_start"`
	DateEnd     string   `json:"date_end"`
	Place       string   `json:"place"`
	Title       string   `json:"title"`
	Level       []string `json:"level"`
	Description string   `json:"description"`
}

type Certificate struct {
	Title       string `json:"title"`
	Description string `json:"description"`
}

type Profile struct {
	JobTitle     string        `json:"job_title"`
	Contacts     Contact       `json:"contacts"`
	Skills       []Skill       `json:"skills"`
	Experience   []Experience  `json:"experience"`
	Education    []Education   `json:"education"`
	Certificates []Certificate `json:"certificates"`
}

const jsonStr = `{
   "job_title": "senior full stack developer",
   "contacts": {
     "emails": ["cucusenok.work@gmail.com"],
     "github": "cucusenok",
     "address": [
        {
           "country": "usa",
           "city": "las vegas",
           "state": "nevada",
           "country_code": "us",
           "state_code": "nv"
        }
     ],
     "phones": ["88005553535"],
     "social_networks": {
       "telegram": "@cucusenok",
       "linkedin": "linkedin.com/in/cucusenok",
     }
   },
   "skills": [
     { "alias": "typescript", "type": "frontend", "id": 1 },
     { "alias": "javascript", "type": "frontend", "id": 2 },
     { "alias": "react", "type": "frontend", "id": 3 },
     { "alias": "vue.js", "type": "frontend", "id": 4 },
     { "alias": "node.js", "type": "backend", "id": 5 },
     { "alias": "python", "type": "backend", "id": 6 },
     { "alias": "fastapi", "type": "backend", "id": 7 },
     { "alias": "flask", "type": "backend", "id": 8 },
     { "alias": "golang", "type": "backend", "id": 9 },
     { "alias": "docker", "type": "environment & common", "id": 10 },
     { "alias": "virtualbox", "type": "environment & common", "id": 11 },
     { "alias": "rest", "type": "environment & common", "id": 12 },
     { "alias": "solid", "type": "environment & common", "id": 13 }
   ],
   "experience": [
       {
         "date_start": "11-2022",
         "date_end": "05-2023",
         "place": "self-founded startup",
         "title": "founder",
         "position": ["founder"],
         "skills": ["golang", "python", "node.js", "flutter", "aws"],
         "level": ["founder"],
         "description": "Developed music recognition and sharing service across diverse sources Built a robust CI/CD pipeline encompassing backend services, IOS app deployment to TestFlight, Android app packaging to APK, and web deployment using GitHub Actions and AWS (EC2, S3)1 Developed a cross-platform application using Flutter, ensuring seamless functionality on Android, IOS, and web platforms Implemented audio decoding and playback features for streaming services on Android and IOS, utilizing proprietary Flutter-plugins and programming languages such as Java/Kotlin and Swift Leveraged API integrations to enhance application functionality and connectivity, while employing GoLang for the central monolith and Python and Node for microservices. SENIOR FULL STACK ENGINEER 5g Networks"
       },
       {
         "date_start": "11-2020",
         "date_end": "11-2022",
         "place": "5g networks",
         "title": "senior full stack engineer",
         "position": ["engineer"],
         "skills": ["3gpp", "4g/lte", "5g", "o-ran"],
         "level": ["senior"],
         "description": "performed analysis and problem-solving for 3gpp, 4g/lte, 5g, and o-ran architectures. designed user-friendly interfaces and developed automated interface generation systems."
       },
       {
         "date_start": "06-2020",
         "date_end": "11-2020",
         "place": "freeze - online viewing room",
         "title": "senior frontend engineer",
         "position": ["engineer"],
         "skills": ["oauth", "restful api", "webrtc"],
         "level": ["senior"],
         "description": "led and trained a team of engineers. finalized and redesigned mobile and web versions, implemented oauth and two-step verification."
       },
       {
         "date_start": "06-2019",
         "date_end": "06-2020",
         "place": "bike.net",
         "title": "fullstack developer",
         "position": ["developer"],
         "skills": ["solr", "mysql"],
         "level": ["mid-level"],
         "description": "allocated tasks between frontend and backend developers. developed a search engine utilizing solr and mysql."
       },
       {
         "date_start": "01-2017",
         "date_end": "06-2019",
         "place": "private practice / freelance",
         "title": "junior full stack developer",
         "position": ["developer"],
         "skills": ["apache", "rbac", "natural language processing"],
         "level": ["junior"],
         "description": "developed a cloud file manager with rbac. completed an internship at the institute of artificial intelligence of ukraine, working on nlp and machine learning."
       }
   ],
   "education": [
       {
         "date_start": "2016",
         "date_end": "2020",
         "place": "kiev national university",
         "title": "bsc with honours, knu",
         "level": ["bachelor"],
         "description": "information computing technology"
       },
       {
         "date_start": "2020",
         "date_end": "2022",
         "place": "kiev national university",
         "title": "msc, knu",
         "level": ["master"],
         "description": "information computing technology"
       }
   ],
   "certificates": [
       {
         "title": "cisco networking academy",
         "description": "cisco, introduction to networks; cisco, programming essentials in python; cisco, routing and switching essentials; cisco, introduction to cybersecurity."
       }
   ]
}`

func MainTest(t *testing.T) {
	// тут польностью проверить тип

	var profile Profile
	err := json.Unmarshal([]byte(jsonStr), &profile)
	if err != nil {
		fmt.Println("Error unmarshaling JSON:", err)
		return
	}

	// Тут получить результат распознавания
	// res := Parse(CV_FULL_PARAGRAPH_FORMATED)

	// TODO: сравнить полученное с profile

	// Output the struct
	fmt.Printf("%+v\n", profile)

}
