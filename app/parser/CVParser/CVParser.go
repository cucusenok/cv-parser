package CVParser

import (
	"awesomeProject2/app/parser"
	"awesomeProject2/app/parser/work_duration"
	"awesomeProject2/app/spell"
	"database/sql"
	"fmt"
	"os"
	"regexp"
	"strings"
	"unicode"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

var spellInstance *spell.Spell

var (
	regexCorrectEndOfSentence *regexp.Regexp = regexp.MustCompile(`\b([^0-9\s]+)1\b`)
)

type ExperienceString struct {
	Date        string `json:"date"`
	Title       string `json:"title"`
	Description string `json:"description"`
}

type SentenceData struct {
	Index          int      `json:"index"`
	Length         int      `json:"length"`
	Commas         int      `json:"commas"`
	Words          []string `json:"words"`
	Skills         []string `json:"skills"`
	Positions      []string `json:"position"`
	Level          []string `json:"level"`
	Date           string   `json:"date"`
	Sentence       string   `json:"sentence"`
	UpperCaseWords []string `json:"upperCaseWords"`
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

/*
Функция isUppercase проверяет написано ли слово в верхнем регистре
*/
func isUppercase(word string) bool {
	hasLetter := false
	for _, r := range word {
		if unicode.IsLetter(r) {
			if !unicode.IsUpper(r) {
				return false
			}
			hasLetter = true
		} else if !unicode.IsLetter(r) && !unicode.IsNumber(r) {
			// Если символ не буква и не цифра (спецсимвол), пропускаем его
			continue
		}
	}
	// Проверяем, есть ли в слове хотя бы одна буква
	return hasLetter
}

func ParseCV(text string) ([]ExperienceString, error) {
	err := godotenv.Load()
	spellInstance, err = LoadSpellFromDB()
	experienceList := []ExperienceString{}
	sentences := []SentenceData{}

	cvData := text
	cvData = regexCorrectEndOfSentence.ReplaceAllString(cvData, "$1.")
	cvData = strings.ReplaceAll(cvData, ". ' ", ". \n · ")
	cvData = strings.ReplaceAll(cvData, " ' ", " \n · ")

	paragraphs := strings.Split(cvData, "\n")
	filteredParagraphs := []string{}

	// фильтруем строки, убирая пустые
	for _, paragraph := range paragraphs {
		trimmedParagraph := strings.TrimSpace(paragraph) // Убираем лишние пробелы
		if trimmedParagraph != "" {
			filteredParagraphs = append(filteredParagraphs, trimmedParagraph)
		}
	}

	for index, paragraph := range filteredParagraphs {
		date := ""
		skills := []string{}
		positions := []string{}
		levels := []string{}
		upperCaseWords := []string{}

		combinations := parser.GenerateCombinations(strings.Split(paragraph, " "))

		words := strings.Fields(paragraph)

		for _, combination := range combinations {
			list, _ := spellInstance.Lookup(combination, spell.SuggestionLevel(spell.LevelClosest))
			if len(list) == 0 {
				continue
			}
			for _, l := range list {
				/*
					TODO некорректно находит скилл.
					Вход: middle Full Stack Developer 7+ YEARS OF EXPERIENCE
					Выход: ["Full Stack", "C++"]
				*/
				if l.WordData["type"] == "skill" && !parser.ContainsItem(skills, l.Word) {
					skills = append(skills, l.Word)
				}
				if l.WordData["type"] == "position" && !parser.ContainsItem(positions, l.Word) {
					positions = append(positions, l.Word)
				}
				// TODO добавить "Senior" в бд
				if l.WordData["type"] == "level" && !parser.ContainsItem(levels, l.Word) {
					levels = append(levels, l.Word)
				}
			}
		}

		for _, word := range words {
			// Убираем пунктуацию в конце слова
			word = strings.Trim(word, ",.!?")
			if isUppercase(word) {
				upperCaseWords = append(upperCaseWords, word)
			}
		}

		dataRange := work_duration.RegexDateRangeExcludeEnd.FindAllString(paragraph, -1)
		isDate := dataRange != nil
		if isDate {
			test, _ := work_duration.ParsePeriod(paragraph)
			fmt.Println("test: ", test)
			date = paragraph
		}

		sentences = append(sentences, SentenceData{
			Index:          index,
			Length:         len(strings.Fields(paragraph)),
			Words:          strings.Fields(paragraph),
			Commas:         strings.Count(paragraph, ","),
			Skills:         skills,
			Positions:      positions,
			Level:          levels,
			Date:           date,
			Sentence:       paragraph,
			UpperCaseWords: upperCaseWords,
		})
	}

	if err != nil {
		fmt.Println("err: ", err)
	}
	return experienceList, nil
}
