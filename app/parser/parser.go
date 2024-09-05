package parser

import (
	"cv-parser/spell"
	"database/sql"
	"fmt"
	"math"
	"os"
	"regexp"
	"strings"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

var spellInstance *spell.Spell

/*
	// Мы нашли должность, и понимает что должность занимает 90% строки
	Co-Founder, Chief Design OfficerCo-Founder, Chief Design Officer
	// мы ничего не достали отсюда, но понимаем, что строка коротка и до этого была должность - это может быть рекомендацией к тому чтобы быть компанией
	Berkana Tech SolutionsBerkana Tech Solutions
	Aug 2023 - Present · 1 yr 1 mo -- мы смогли спарсить промежуток
*/

// FULLSTACK DEVELOPER June, 2019 - June, 2020
// тут есть position + тут есть промежуток времени - значит скорее всего это exprience title

var (
	regexZipCode        *regexp.Regexp = regexp.MustCompile(`\d{4,5}(?:[-\s]\d{4})?`)
	regexPhone          *regexp.Regexp = regexp.MustCompile(`\+?\d{1,4}?[-.\s]?\(?\d{1,3}?\)?[-.\s]?\d{1,4}[-.\s]?\d{1,4}[-.\s]?\d{1,9}`)
	regexEmail          *regexp.Regexp = regexp.MustCompile("[a-z0-9!#$%&'*+/=?^_`{|}~-]+(?:\\.[a-z0-9!#$%&'*+/=?^_`{|}~-]+)*@(?:[a-z0-9](?:[a-z0-9-]*[a-z0-9])?\\.)+[a-z0-9](?:[a-z0-9-]*[a-z0-9])?")
	regexDomain         *regexp.Regexp = regexp.MustCompile(`(?:[_a-z0-9](?:[_a-z0-9-]{0,61}[a-z0-9])?\.)+(?:[a-z](?:[a-z0-9-]{0,61}[a-z0-9])?)?`)
	regexNickNameWithAt *regexp.Regexp = regexp.MustCompile(`^[a-zA-Z0-9@\-\.]*$`)
	// TODO: добавить учет present, для строк типа: Aug 2023 - Present · 1 yr 1 mo
	regexYear                    *regexp.Regexp = regexp.MustCompile(`(19|20)\d{2}`)
	regexDateRange               *regexp.Regexp = regexp.MustCompile(`(19|20)\d{2}.*(19|20)\d{2}`)
	regexAllSymbolsAfterDot      *regexp.Regexp = regexp.MustCompile(`\..*`)
	regexAllSymbolsAfterQuestion *regexp.Regexp = regexp.MustCompile(`\?.*`)
	regexURL                     *regexp.Regexp = regexp.MustCompile(`(?:[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}|[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}(/[a-zA-Z0-9._%+-]*)*)`)
	RegexDateRangeExcludeEnd     *regexp.Regexp = regexp.MustCompile(`(19|20)\d{2}(.*(19|20)\d{2})?`)
	regexListItems               *regexp.Regexp = regexp.MustCompile(`·.*?\.`)
	regexCorrentEndOfSentence    *regexp.Regexp = regexp.MustCompile(`\b([^0-9\s]+)1\b`)
)

func ContainsItem(arr []string, item string) bool {
	for _, val := range arr {
		if val == item {
			return true
		}
	}
	return false
}

func GenerateCombinations(input []string) []string {
	result := append([]string{}, input...)

	for i := 0; i < len(input)-1; i++ {
		for j := i + 1; j < len(input) && j-i+1 <= 3; j++ {
			substring := strings.Join(input[i:j+1], " ")
			result = append(result, substring)
		}
	}

	return result
}

func scanQuery(db *sql.DB, sql string, args []any, fn func(*sql.Rows) error) error {
	rows, err := db.Query(sql, args...)
	if err != nil {
		return err
	}
	defer rows.Close()
	for rows.Next() {
		if err := fn(rows); err != nil {
			return err
		}
	}
	return rows.Err()
}

func LoadSpellFromDB() (*spell.Spell, error) {
	err := godotenv.Load()

	if err != nil {
		panic(err)
	}
	pgConnectionStr := os.Getenv("PG_CONNECTION_STR")

	s := spell.New()
	db, err := sql.Open("postgres", pgConnectionStr)
	if err != nil {
		return nil, err
	}
	fmt.Println("loading...")
	if err = scanQuery(db, `select alias from skills_aliases`, nil, func(rows *sql.Rows) error {
		var alias string
		if err := rows.Scan(&alias); err != nil {
			if err.Error() == "sql: Scan error on column index 2, name \"cnt\": converting NULL to uint64 is unsupported" {
				return nil
			}
			return err
		}

		if alias != "" {
			s.AddEntry(spell.Entry{
				Frequency: 1,
				Word:      alias,
				WordData: spell.WordData{
					"type": "skill",
				},
			})
		}

		return nil
	}); err != nil {
		return nil, err
	}

	if err = scanQuery(db, `select alias from positions_aliases`, nil, func(rows *sql.Rows) error {
		var alias string
		if err := rows.Scan(&alias); err != nil {
			if err.Error() == "sql: Scan error on column index 2, name \"cnt\": converting NULL to uint64 is unsupported" {
				return nil
			}
			return err
		}

		if alias != "" {
			s.AddEntry(spell.Entry{
				Frequency: 1,
				Word:      alias,
				WordData: spell.WordData{
					"type": "position",
				},
			})
		}

		return nil
	}); err != nil {
		return nil, err
	}

	if err = scanQuery(db, `select level from levels`, nil, func(rows *sql.Rows) error {
		var level string
		if err := rows.Scan(&level); err != nil {
			if err.Error() == "sql: Scan error on column index 2, name \"cnt\": converting NULL to uint64 is unsupported" {
				return nil
			}
			return err
		}

		if level != "" {
			s.AddEntry(spell.Entry{
				Frequency: 1,
				Word:      level,
				WordData: spell.WordData{
					"type": "level",
				},
			})
		}

		return nil
	}); err != nil {
		return nil, err
	}

	return s, nil
}

type CVParseResult struct {
	Experiences []ExperienceString
}

func (cvp *CVParseResult) ToJsonDSL() {
	// TODO: конвертация в нужный формат для API
}

type ExperienceString struct {
	Title       string `json:"title"`
	Name        string `json:"name"`
	Start       string `json:"start"`
	End         string `json:"end"`
	Description string `json:"description"`
}

type Company struct {
	Name  string `json:"name"`
	Start string `json:"start"`
	End   string `json:"end"`
	Role  string `json:"role"`
	Tasks string `json:"tasks"`
}

type CVData struct {
	Experience []ExperienceString `json:"experience"`
	Company    []Company          `json:"company"`
	// Links      []string     `json:"links"`
	Mail string `json:"mail"`
}

func FindIndex(slice []int, condition func(int) bool) int {
	for i, v := range slice {
		if condition(v) {
			return i
		}
	}
	return -1
}

// test.domain.com -> test___domain___com

func Parse(text string) (*CVParseResult, error) {
	err := godotenv.Load()

	// frontTitles := []string{
	// 	"HTML", "HyperText Markup Language",
	// 	"CSS", "Cascading Style Sheets",
	// 	"JS", "JavaScript",
	// 	"DOM", "Document Object Model",
	// 	"AJAX", "Asynchronous JavaScript and XML",
	// 	"API", "Application Programming Interface",
	// 	"SPA", "Single Page Application",
	// 	"PWA", "Progressive Web Application",
	// 	"SSR", "Server-Side Rendering",
	// 	"CSR", "Client-Side Rendering",
	// 	"CDN", "Content Delivery Network",
	// 	"BEM", "Block Element Modifier",
	// 	"SASS/SCSS", "Syntactically Awesome Stylesheets",
	// 	"LESS", "Leaner Style Sheets",
	// 	"ECMAScript",
	// 	"ES6/ESNext", "ECMAScript 6/Next",
	// 	"JSX", "JavaScript XML",
	// 	"TS", "TypeScript",
	// 	"JSON", "JavaScript Object Notation",
	// 	"SEO", "Search Engine Optimization",
	// 	"UX", "User Experience",
	// 	"UI", "User Interface",
	// 	"WYSIWYG", "What You See Is What You Get",
	// 	"CLI", "Command Line Interface",
	// 	"VCS", "Version Control System",
	// 	"LTS", "Long Term Support",
	// 	"CI/CD", "Continuous Integration/Continuous Deployment",
	// 	"Continuous Integration",
	// 	"Continuous Deployment",
	// 	"NPM", "Node Package Manager",
	// 	"YARN", "Yet Another Resource Negotiator",
	// 	"CSSOM", "CSS Object Model",
	// 	"SVG", "Scalable Vector Graphics",
	// 	"AMP", "Accelerated Mobile Pages",
	// 	"SSR", "Server-Side Rendering",
	// 	"Lighthouse",
	// 	"GULP",
	// 	"WEBPACK",
	// }

	spellInstance, err = LoadSpellFromDB()
	skills := []string{}
	positions := []string{}
	// experiences := []Experience{}

	cvData := text
	cvData = regexCorrentEndOfSentence.ReplaceAllString(cvData, "$1.")
	cvData = strings.ReplaceAll(cvData, ". ' ", ". \n · ")
	cvData = strings.ReplaceAll(cvData, " ' ", " \n · ")
	lowerCaseText := strings.ToLower(cvData)
	fmt.Println("\n  cvData: ", cvData, "\n ")

	paragraphs := strings.Split(lowerCaseText, "\n")
	filteredParagraphs := []string{}

	experienceSectionsIndexes := []int{}
	sections := [][]string{}
	// links := regexURL.FindAllString(text, -1)

	for _, paragraph := range paragraphs {
		trimmedParagraph := strings.TrimSpace(paragraph) // Убираем лишние пробелы
		if trimmedParagraph != "" {
			filteredParagraphs = append(filteredParagraphs, trimmedParagraph)
		}
	}

	experienceList := []ExperienceString{}

	fmt.Println("\n ======================================= filteredParagraphs ===================================================\n ")

	/*
			Примеры filteredParagraphs:

			[
			  'Nov, 2022 - May, 2023 ' Developed music recognition and sharing service across diverse sources1 ' Built a robust CI/CD pipeline encompassing backend services, IOS app deployment to TestFlight, Android app packaging to APK, and web deployment using GitHub Actions and AWS (EC2, S3)1 ' Developed a cross-platform application using Flutter, ensuring seamless functionality on Android, IOS, and web platforms1 ' Implemented audio decoding and playback features for streaming services on Android and IOS, utilizing proprietary Flutterplugins and programming languages such as Java/Kotlin and Swift1 ' Leveraged API integrations to enhance application functionality and connectivity, while employing GoLang for the central monolith and Python and Node for microservices.',
		      'SENIOR FULL STACK ENGINEER',
		      'Nov, 2020 - Nov, 2022'
			]

	*/
	for i, paragraph := range filteredParagraphs {
		fmt.Println(paragraph)
		isDateParagraph := RegexDateRangeExcludeEnd.MatchString(paragraph)
		if isDateParagraph {
			experienceSectionsIndexes = append(experienceSectionsIndexes, i-1)
		}
	}

	fmt.Println("\n ======================================= filteredParagraphs end =================================================== \n ")

	for i, start := range experienceSectionsIndexes {
		var end int
		if i+1 < len(experienceSectionsIndexes) {
			end = experienceSectionsIndexes[i+1]
		} else {
			end = len(filteredParagraphs)
		}
		sections = append(sections, filteredParagraphs[start:end])
	}

	for _, section := range sections {
		sentenceList := []string{}
		years := []string{}
		sentencePositions := []string{}
		sentenceSkills := []string{}
		sentenceLevels := []string{}
		combinations := GenerateCombinations(strings.Split(section[0], " "))

		for _, combination := range combinations {
			list, _ := spellInstance.Lookup(combination, spell.SuggestionLevel(spell.LevelClosest))
			if len(list) == 0 {
				continue
			}
			for _, l := range list {
				// TODO: добавить проверку на кол-во символов в слове. Но думаю, что не корректно в случае с QA, PM и тд.
				if len(l.Word) < 3 {
					continue
				}

				if l.WordData["type"] == "level" {
					sentenceLevels = append(sentenceLevels, l.Word)
				}
				if l.WordData["type"] == "skill" {
					if !ContainsItem(skills, l.Word) {
						skills = append(skills, l.Word)
					}
					sentenceSkills = append(sentenceSkills, l.Word)

				}
				if l.WordData["type"] == "position" {
					if !ContainsItem(positions, l.Word) {
						positions = append(positions, l.Word)
					}
					sentencePositions = append(sentencePositions, l.Word)
				}
			}
		}

		for _, paragraph := range section {
			isDateParagraph := RegexDateRangeExcludeEnd.MatchString(paragraph)
			isListItem := regexListItems.MatchString(paragraph)

			if isDateParagraph {
				dataRange := RegexDateRangeExcludeEnd.FindAllString(paragraph, -1)
				years = regexYear.FindAllString(dataRange[0], -1)
				continue
			} else if isListItem {
				sentenceList = append(sentenceList, paragraph)
				continue
			}
		}
		endDate := fmt.Sprintf("%v", math.Inf(1))
		if len(years) >= 2 {
			endDate = years[1]
		}
		experienceList = append(experienceList, ExperienceString{
			Title:       section[0],
			Name:        strings.Join(sentenceLevels, ", "),
			Start:       years[0],
			End:         endDate,
			Description: strings.Join(sentenceList, "\n "),
		})
	}

	// fmt.Println("skills: ",skills)
	//fmt.Println(positions)

	// for _, experience := range experiences {
	// fmt.Println(experience)
	// }

	for _, experienceItem := range experienceList {

		fmt.Println("\n =============== experienceItem =============")
		fmt.Println("\n experienceItem.Title: ", experienceItem.Title)
		fmt.Println("\n experienceItem.Name: ", experienceItem.Name)
		fmt.Println("\n experienceItem.Start: ", experienceItem.Start)
		fmt.Println("\n experienceItem.End: ", experienceItem.End)
		fmt.Println("\n experienceItem.Description: ", experienceItem.Description)
		fmt.Println("\n =============== experienceItem end =============")
	}

	// cvData := CVData{
	// 	Experience: experiences,
	// 	Mail:       regexEmail.FindString(text),
	// 	// Links:      links,
	// }

	// fmt.Println("cvData: ", cvData)

	/*	for _, position := range positions {
		fmt.Println(position)
	}*/

	fmt.Println(err)
	return &CVParseResult{
		Experiences: experienceList,
	}, nil

}
