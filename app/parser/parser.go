package parser

import (
	"awesomeProject2/app/spell"
	"database/sql"
	"fmt"
	"math"
	"regexp"
	"strings"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

var spellInstance *spell.Spell

var (
	regexZipCode                 *regexp.Regexp = regexp.MustCompile(`\d{4,5}(?:[-\s]\d{4})?`)
	regexPhone                   *regexp.Regexp = regexp.MustCompile(`\+?\d{1,4}?[-.\s]?\(?\d{1,3}?\)?[-.\s]?\d{1,4}[-.\s]?\d{1,4}[-.\s]?\d{1,9}`)
	regexEmail                   *regexp.Regexp = regexp.MustCompile("[a-z0-9!#$%&'*+/=?^_`{|}~-]+(?:\\.[a-z0-9!#$%&'*+/=?^_`{|}~-]+)*@(?:[a-z0-9](?:[a-z0-9-]*[a-z0-9])?\\.)+[a-z0-9](?:[a-z0-9-]*[a-z0-9])?")
	regexDomain                  *regexp.Regexp = regexp.MustCompile(`(?:[_a-z0-9](?:[_a-z0-9-]{0,61}[a-z0-9])?\.)+(?:[a-z](?:[a-z0-9-]{0,61}[a-z0-9])?)?`)
	regexNickNameWithAt          *regexp.Regexp = regexp.MustCompile(`^[a-zA-Z0-9@\-\.]*$`)
	regexYear                    *regexp.Regexp = regexp.MustCompile(`(19|20)\d{2}`)
	regexDateRange               *regexp.Regexp = regexp.MustCompile(`(19|20)\d{2}.*(19|20)\d{2}`)
	regexAllSymbolsAfterDot      *regexp.Regexp = regexp.MustCompile(`\..*`)
	regexAllSymbolsAfterQuestion *regexp.Regexp = regexp.MustCompile(`\?.*`)
	regexURL                     *regexp.Regexp = regexp.MustCompile(`(?:[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}|[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}(/[a-zA-Z0-9._%+-]*)*)`)
	regexDateRangeExcludeEnd     *regexp.Regexp = regexp.MustCompile(`(19|20)\d{2}(.*(19|20)\d{2})?`)
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
	//pgConnectionStr := os.Getenv("PG_CONNECTION_STR")
	pgConnectionStr := "user=postgres dbname=parser sslmode=disable password=mysecretpassword host=localhost port=1234"

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

	return s, nil
}

type Experience struct {
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
	Experience []Experience `json:"experience"`
	Company    []Company    `json:"company"`
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

func Parse(text string) {
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

	// for _, paragraph := range filteredParagraphs {
	// 	sentences := strings.Split(paragraph, ".")
	// 	// TODO: обрабатывать домены, из-за split(paragraph, ".") они тоже разбиваются test.com -> [test, com]
	// 	sentencePositions := []string{}
	// 	for _, sentence := range sentences {
	// 		combinations := GenerateCombinations(strings.Split(sentence, " "))
	// 		//foundedPosition := false

	// 		for _, combination := range combinations {
	// 			list, _ := spellInstance.Lookup(combination, spell.SuggestionLevel(spell.LevelClosest))
	// 			if len(list) == 0 {
	// 				continue
	// 			}
	// 			for _, l := range list {
	// 				// TODO: добавил проверку на кол-во символов в слове. Но думаю, что не корректно в случае с QA, PM и тд.
	// 				if len(l.Word) < 3 {
	// 					continue
	// 				}
	// 				// fmt.Println("l: ", l.Word, l.WordData["type"])
	// 				if l.WordData["type"] == "skill" && !ContainsItem(skills, l.Word) {
	// 					skills = append(skills, l.Word)
	// 				}
	// 				if l.WordData["type"] == "position" {
	// 					if !ContainsItem(positions, l.Word) {
	// 						positions = append(positions, l.Word)
	// 					}
	// 					sentencePositions = append(sentencePositions, l.Word)
	// 					//foundedPosition = true
	// 				}
	// 			}
	// 		}
	// 		// fmt.Println("sentencePositions: ", sentencePositions)
	// 		if len(sentencePositions) > 0 {
	// 			dataRange := regexDateRangeExcludeEnd.FindAllString(sentence, -1)
	// 			if len(dataRange) == 0 {
	// 				continue
	// 			}
	// 			// fmt.Println("dataRange", dataRange)

	// 			years := regexYear.FindAllString(dataRange[0], -1)
	// 			endDate := fmt.Sprintf("%v", math.Inf(1))
	// 			if len(years) >= 2 {
	// 				endDate = years[1]
	// 			}
	// 			experiences = append(experiences, Experience{
	// 				Name:     sentencePositions[0],
	// 				Start:    years[0],
	// 				End:      endDate,
	// 				Sentence: sentence,
	// 			})

	// 		}
	// 		//fmt.Println(foundedPosition)

	// 	}
	// }

	experienceList := []Experience{}

	fmt.Println("\n ======================================= filteredParagraphs ===================================================\n ")

	for i, paragraph := range filteredParagraphs {
		fmt.Println(paragraph)
		isDateParagraph := regexDateRangeExcludeEnd.MatchString(paragraph)
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
		sentencePositions := ""
		combinations := GenerateCombinations(strings.Split(section[0], " "))

		for _, combination := range combinations {
			list, _ := spellInstance.Lookup(combination, spell.SuggestionLevel(spell.LevelClosest))
			if len(list) == 0 {
				continue
			}
			for _, l := range list {
				// TODO: добавил проверку на кол-во символов в слове. Но думаю, что не корректно в случае с QA, PM и тд.
				if len(l.Word) < 3 {
					continue
				}
				if l.WordData["type"] == "skill" && !ContainsItem(skills, l.Word) {
					skills = append(skills, l.Word)
				}
				if l.WordData["type"] == "position" {
					if !ContainsItem(positions, l.Word) {
						positions = append(positions, l.Word)
					}
					sentencePositions = l.Word

				}
			}
		}

		for _, paragraph := range section {
			isDateParagraph := regexDateRangeExcludeEnd.MatchString(paragraph)
			isListItem := regexListItems.MatchString(paragraph)

			if isDateParagraph {
				dataRange := regexDateRangeExcludeEnd.FindAllString(paragraph, -1)
				years = regexYear.FindAllString(dataRange[0], -1)
				continue
			} else if isListItem {
				sentenceList = append(sentenceList, paragraph)
				continue
			}

			// sentences := strings.Split(paragraph, ".")
			// for _, sentence := range sentences {
			// 	combinations := GenerateCombinations(strings.Split(sentence, " "))

			// 	for _, combination := range combinations {
			// 		list, _ := spellInstance.Lookup(combination, spell.SuggestionLevel(spell.LevelClosest))
			// 		if len(list) == 0 {
			// 			continue
			// 		}
			// 		for _, l := range list {
			// 			// TODO: добавил проверку на кол-во символов в слове. Но думаю, что не корректно в случае с QA, PM и тд.
			// 			if len(l.Word) < 3 {
			// 				continue
			// 			}
			// 			if l.WordData["type"] == "skill" && !ContainsItem(skills, l.Word) {
			// 				skills = append(skills, l.Word)
			// 			}
			// 			if l.WordData["type"] == "position" {
			// 				if !ContainsItem(positions, l.Word) {
			// 					positions = append(positions, l.Word)
			// 				}
			// 				sentencePositions = l.Word

			// 			}
			// 		}
			// 	}
			// }
		}
		endDate := fmt.Sprintf("%v", math.Inf(1))
		if len(years) >= 2 {
			endDate = years[1]
		}
		experienceList = append(experienceList, Experience{
			Title:       section[0],
			Name:        sentencePositions,
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

}
