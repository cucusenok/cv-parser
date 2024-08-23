package experience

import (
	"awesomeProject2/app/spell"
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
	regexYear                 *regexp.Regexp = regexp.MustCompile(`(19|20)\d{2}`)
	regexDateRangeExcludeEnd  *regexp.Regexp = regexp.MustCompile(`(19|20)\d{2}(.*(19|20)\d{2})?`)
	regexListItems            *regexp.Regexp = regexp.MustCompile(`·.*?\.`)
	regexCorrentEndOfSentence *regexp.Regexp = regexp.MustCompile(`\b([^0-9\s]+)1\b`)
)

type Experience struct {
	Start       string `json:"start_date"`
	End         string `json:"end_date"`
	Title       string `json:"title"`
	Description string `json:"description"`
}

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

type CVData struct {
	Experience []Experience `json:"experience"`
	// Links      []string     `json:"links"`
	Mail string `json:"mail"`
}

// test.domain.com -> test___domain___com

func ParseExperience(text string) (*[]Experience, error) {
	err := godotenv.Load()

	spellInstance, err = LoadSpellFromDB()
	skills := []string{}
	positions := []string{}
	experienceList := []Experience{}

	cvData := text
	cvData = regexCorrentEndOfSentence.ReplaceAllString(cvData, "$1.")
	cvData = strings.ReplaceAll(cvData, ". ' ", ". \n · ")
	cvData = strings.ReplaceAll(cvData, " ' ", " \n · ")
	lowerCaseText := strings.ToLower(cvData)
	fmt.Println("\n  cvData: ", cvData, "\n ")

	paragraphs := strings.Split(lowerCaseText, "\n")
	filteredParagraphs := []string{}

	// links := regexURL.FindAllString(text, -1)

	for _, paragraph := range paragraphs {
		trimmedParagraph := strings.TrimSpace(paragraph) // Убираем лишние пробелы
		if trimmedParagraph != "" {
			filteredParagraphs = append(filteredParagraphs, trimmedParagraph)
		}
	}

	for _, paragraph := range filteredParagraphs {
		sentences := strings.Split(paragraph, ".")
		// TODO: обрабатывать домены, из-за split(paragraph, ".") они тоже разбиваются test.com -> [test, com]
		sentencePositions := []string{}
		for _, sentence := range sentences {
			combinations := GenerateCombinations(strings.Split(sentence, " "))
			//foundedPosition := false

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
					// fmt.Println("l: ", l.Word, l.WordData["type"])
					if l.WordData["type"] == "skill" && !ContainsItem(skills, l.Word) {
						skills = append(skills, l.Word)
					}
					if l.WordData["type"] == "position" {
						if !ContainsItem(positions, l.Word) {
							positions = append(positions, l.Word)
						}
						sentencePositions = append(sentencePositions, l.Word)
						//foundedPosition = true
					}
				}
			}
			// fmt.Println("sentencePositions: ", sentencePositions)
			if len(sentencePositions) > 0 {
				dataRange := regexDateRangeExcludeEnd.FindAllString(sentence, -1)
				if len(dataRange) == 0 {
					continue
				}
				// fmt.Println("dataRange", dataRange)

				years := regexYear.FindAllString(dataRange[0], -1)
				endDate := fmt.Sprintf("%v", math.Inf(1))
				if len(years) >= 2 {
					endDate = years[1]
				}
				experienceList = append(experienceList, Experience{
					Title:       sentencePositions[0],
					Start:       years[0],
					End:         endDate,
					Description: "",
				})

			}
			//fmt.Println(foundedPosition)

		}
	}
	fmt.Println("err: ", err)
	return &experienceList, nil

}
