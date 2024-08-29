package experience

import (
	"awesomeProject2/app/parser/work_duration"
	"awesomeProject2/app/spell"
	"database/sql"
	"fmt"
	"os"
	"regexp"
	"strings"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

var spellInstance *spell.Spell

const MATCH_PERCENTAGE = 50

var (
	regexDateRangeExcludeEnd  *regexp.Regexp = regexp.MustCompile(`(19|20)\d{2}(.*(19|20)\d{2})?`)
	regexCorrectEndOfSentence *regexp.Regexp = regexp.MustCompile(`\b([^0-9\s]+)1\b`)
)

type ExperienceString struct {
	Date        string `json:"date"`
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
	Experience []ExperienceString `json:"experience"`
	// Links      []string     `json:"links"`
	Mail string `json:"mail"`
}

func addToDescription(experience *ExperienceString, newDescription string) {
	if experience.Description != "" {
		experience.Description += "\n"
	}
	experience.Description += newDescription
}

/*
	splitOnDate разбивает строку на две части, одна из которых дата, вторая - любой текст, который не относится к дате. И возвращает две строки соответственно
	Пример splitOnDate:

	FOUNDER (START-UP) Nov, 2022 - May, 2023 => ("FOUNDER (START-UP)","Nov, 2022 - May, 2023")
*/

func splitOnDate(text string) (string, string) {
	parts := strings.Split(text, " ")
	var beforeDateParts []string
	var dateParts []string

	for _, part := range parts {
		if work_duration.RegexDates.MatchString(part) {
			dateParts = append(dateParts, part)
		} else if len(dateParts) == 0 { // Мы еще не начали находить дату
			beforeDateParts = append(beforeDateParts, part)
		} else { // Мы уже начали находить дату, но часть не соответствует регулярному выражению
			dateParts = append(dateParts, part)
		}
	}

	beforeDate := strings.Join(beforeDateParts, " ")
	dateBlock := strings.Join(dateParts, " ")

	return beforeDate, dateBlock
}

// test.domain.com -> test___domain___com

func ParseExperience(text string) ([]ExperienceString, error) {
	err := godotenv.Load()

	spellInstance, err = LoadSpellFromDB()

	skills := []string{}
	positions := []string{}
	experienceList := []ExperienceString{}
	experience := ExperienceString{}
	currentTitle := ""

	cvData := text
	cvData = regexCorrectEndOfSentence.ReplaceAllString(cvData, "$1.")
	cvData = strings.ReplaceAll(cvData, ". ' ", ". \n · ")
	cvData = strings.ReplaceAll(cvData, " ' ", " \n · ")
	lowerCaseText := strings.ToLower(cvData)

	paragraphs := strings.Split(lowerCaseText, "\n")
	filteredParagraphs := []string{}

	// links := regexURL.FindAllString(text, -1)

	for _, paragraph := range paragraphs {
		trimmedParagraph := strings.TrimSpace(paragraph) // Убираем лишние пробелы
		if trimmedParagraph != "" {
			filteredParagraphs = append(filteredParagraphs, trimmedParagraph)
		}
	}
	fmt.Println("filteredParagraphs: ", filteredParagraphs)

	for index, paragraph := range filteredParagraphs {
		isMatchDate := work_duration.IsMatchDate(paragraph)
		if isMatchDate {
			text, date := splitOnDate(paragraph)
			fmt.Println("text, date: ", text, date)
		}

		sentences := strings.Split(paragraph, ".")
		var filteredSentences []string
		for _, sentence := range sentences {
			trimmed := strings.TrimSpace(sentence)
			if trimmed != "" {
				filteredSentences = append(filteredSentences, trimmed)
			}
		}

		// TODO: обрабатывать домены, из-за split(paragraph, ".") они тоже разбиваются test.com -> [test, com]
		// TODO: обрабатывать числа формата nn.nn => 2.2 и тд
		sentencePositions := []string{}

		for _, sentence := range filteredSentences {
			splitSentence := strings.Split(sentence, " ")
			combinations := GenerateCombinations(splitSentence)

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
						sentencePositions = append(sentencePositions, l.Word)
					}
				}
			}
			if len(sentencePositions) > 0 {
				isTitle := (len(sentencePositions)*100)/len(splitSentence) >= MATCH_PERCENTAGE // если совпадений слов >= MATCH_PERCENTAGE(50%), считаю, что это заголовок

				if isTitle {
					currentTitle = sentence
				} else {
					addToDescription(&experience, sentence)
				}
				sentencePositions = []string{}

			} else {
				dataRange := regexDateRangeExcludeEnd.FindAllString(sentence, -1)
				isDate := dataRange != nil
				if isDate {
					experience.Date = sentence
				} else {
					addToDescription(&experience, sentence)
				}
			}
		}

		if len(currentTitle) > 0 {
			if len(experience.Title) == 0 {
				experience.Title = currentTitle
			}
			updateExperienceList := len(experience.Title) > 0 && (experience.Title != currentTitle || (len(filteredParagraphs)-1) == index)
			if updateExperienceList {
				experienceList = append(experienceList, experience)
				experience = ExperienceString{}
			}
		}
	}
	if err != nil {
		fmt.Println("err: ", err)
	}
	return experienceList, nil
}
