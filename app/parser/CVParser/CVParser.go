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

// TODO добавить place
type ExperienceString struct {
	Date        *work_duration.WorkPeriod `json:"date"`
	Title       string                    `json:"title"`
	Positions   []string                  `json:"position"`
	Skills      []string                  `json:"skills"`
	Level       []string                  `json:"level"`
	Description string                    `json:"description"`
}

type JobTitles struct {
	Sentence string `json:"sentence"`
	Index    int    `json:"index"`
}

type SentenceData struct {
	Index          int                       `json:"index"`
	WordsCount     int                       `json:"wordsCount"`
	Commas         int                       `json:"commas"`
	Words          []string                  `json:"words"`
	Skills         []string                  `json:"skills"`
	Positions      []string                  `json:"position"`
	Level          []string                  `json:"level"`
	Date           *work_duration.WorkPeriod `json:"date"`
	Sentence       string                    `json:"sentence"`
	UpperCaseWords []string                  `json:"upperCaseWords"`
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
Функция splitByDate разбивает строку на несколько частей относительно даты.
Пример:
	Вход: "JUNIOR FULL STACK DEVELOPER Jan, 2017 - June, 2019 Private practice / freelance"
	Выход: ["JUNIOR FULL STACK DEVELOPER", "Jan, 2017 - June, 2019", "Private practice / freelance"]
*/

// TODO сделать проверку на месяцы из бд. Баг: JUNIOR меняется на 06 из-за проверки на корень слова

func splitByDate(text string) []string {
	formattedText := work_duration.ReplaceDateWithMonthNumber(text)
	matches := work_duration.RegexDatesWithoutDigit.FindAllStringSubmatch(formattedText, -1)
	groupNames := work_duration.RegexDatesWithoutDigit.SubexpNames()
	dates := []string{}

	for _, match := range matches {
		for i, name := range groupNames {
			if len(name) == 0 {
				continue
			}
			if (name == work_duration.TYPE_DATE_REVERSE || name == work_duration.TYPE_DATE_NORMAL) && i < len(match) && len(match[i]) > 0 {
				dates = append(dates, match[i])
			}
		}
	}
	if len(matches) == 0 {
		return []string{text}
	}

	var date string
	var beforeDate string
	var afterDate string

	if len(dates) == 0 || len(dates) > 2 {
		return []string{text}
	} else if len(dates) == 2 {
		date = strings.TrimSpace(formattedText[strings.Index(formattedText, dates[0]) : strings.Index(formattedText, dates[1])+len(dates[1])])
	} else {
		date = strings.TrimSpace(formattedText[strings.Index(formattedText, dates[0]):])
	}

	beforeDate = strings.TrimSpace(formattedText[:strings.Index(formattedText, date)])
	afterDate = strings.TrimSpace(formattedText[strings.Index(formattedText, date)+len(date):])

	return []string{beforeDate, date, afterDate}
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

/*
Функция calculateAverageLength возвращает среднее значение длины строки исходя из среза, который пробрасывается в аргументы
*/

func calculateAverageLength(data []SentenceData) int {
	totalLength := 0
	for _, item := range data {
		totalLength += item.WordsCount
	}
	averageLength := totalLength / len(data)
	return averageLength
}

func collectDataInRange(sentences, titleSentences []SentenceData) []ExperienceString {
	var result []ExperienceString

	for i := 0; i < len(titleSentences); i++ {
		startIndex := titleSentences[i].Index
		var endIndex int
		if i+1 < len(titleSentences) {
			endIndex = titleSentences[i+1].Index
		} else {
			endIndex = len(sentences)
		}

		var rangeData ExperienceString
		var title string
		date := &work_duration.WorkPeriod{
			DateStart: "",
			DateEnd:   "",
		}
		skills := []string{}
		positions := []string{}
		levels := []string{}
		description := []string{}

		for j := startIndex; j < endIndex; j++ {
			sentence := sentences[j]
			if j == startIndex {
				title = sentence.Sentence
			} else if len(sentence.Date.DateStart) > 0 {
				date = sentence.Date
			} else {
				description = append(description, sentence.Sentence)
			}

			if len(sentence.Level) > 0 {
				for _, level := range sentence.Level {
					if !parser.ContainsItem(levels, level) {
						levels = append(levels, level)
					}
				}
			}
			if len(sentence.Positions) > 0 {
				for _, position := range sentence.Positions {
					if !parser.ContainsItem(positions, position) {
						positions = append(positions, position)
					}
				}
			}
			if len(sentence.Skills) > 0 {
				for _, skill := range sentence.Skills {
					if !parser.ContainsItem(skills, skill) {
						skills = append(skills, skill)
					}
				}
			}

			rangeData = ExperienceString{
				Title:       title,
				Skills:      skills,
				Positions:   positions,
				Level:       levels,
				Date:        date,
				Description: strings.Join(description, " "),
			}
		}
		result = append(result, rangeData)
	}

	return result
}

func ParseCV(text string) ([]ExperienceString, error) {
	err := godotenv.Load()
	spellInstance, err = LoadSpellFromDB()
	experienceList := []ExperienceString{}
	sentences := []SentenceData{}
	titleSentences := []SentenceData{}

	cvData := text
	cvData = regexCorrectEndOfSentence.ReplaceAllString(cvData, "$1.")
	cvData = strings.ReplaceAll(cvData, ". ' ", ". \n · ")
	cvData = strings.ReplaceAll(cvData, " ' ", " \n · ")

	jobTitles := []JobTitles{}
	sentencesWithDates := []SentenceData{}

	var averageWordsCount int

	paragraphs := strings.Split(cvData, "\n")
	filteredParagraphs := []string{}

	// фильтруем строки, убирая пустые
	for _, paragraph := range paragraphs {
		trimmedParagraph := strings.TrimSpace(paragraph) // Убираем лишние пробелы
		if trimmedParagraph != "" {
			dataRange := work_duration.RegexDateRangeExcludeEnd.FindAllString(trimmedParagraph, -1)
			isDate := dataRange != nil
			if isDate {
				splitParagraphs := splitByDate(trimmedParagraph)
				for _, splitParagraph := range splitParagraphs {
					if len(splitParagraph) > 0 {
						filteredParagraphs = append(filteredParagraphs, splitParagraph)
					}
				}
			} else {
				filteredParagraphs = append(filteredParagraphs, trimmedParagraph)
			}
		}
	}

	for index, paragraph := range filteredParagraphs {
		date := &work_duration.WorkPeriod{
			DateStart: "",
			DateEnd:   "",
		}

		skills := []string{}
		positions := []string{}
		levels := []string{}
		upperCaseWords := []string{}
		combinations := parser.GenerateCombinations(strings.Split(paragraph, " "))
		words := strings.Fields(paragraph)

		for _, combination := range combinations {
			list, _ := spellInstance.Lookup(strings.ToLower(combination), spell.SuggestionLevel(spell.LevelClosest))

			if len(list) > 0 && list[0].Distance > 3 {
				continue
			}

			for _, l := range list {
				if l.WordData["type"] == "skill" && !parser.ContainsItem(skills, l.Word) {
					skills = append(skills, l.Word)
				}
				if l.WordData["type"] == "position" && !parser.ContainsItem(positions, l.Word) {
					positions = append(positions, l.Word)
				}
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
			parsedDate, _ := work_duration.ParsePeriod(paragraph)
			date = parsedDate
		}

		sentences = append(sentences, SentenceData{
			Index:          index,
			WordsCount:     len(strings.Fields(paragraph)),
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

	averageWordsCount = calculateAverageLength(sentences)

	for _, sentence := range sentences {
		jobTitle := JobTitles{}
		// Если в строке Skills + Positions + Level занимают >= 50% строки, тогда считаю, что это может быть jobTitle.
		if len(sentence.Positions) > 0 &&
			(((len(sentence.Skills) + len(sentence.Positions) + len(sentence.Level)) * 100 / sentence.WordsCount) >= 50) &&
			sentence.WordsCount < averageWordsCount {
			jobTitle = JobTitles{
				Sentence: sentence.Sentence,
				Index:    sentence.Index,
			}
			jobTitles = append(jobTitles, jobTitle)
		}
		if len(sentence.Date.DateStart) > 0 {
			sentencesWithDates = append(sentencesWithDates, sentence)
		}
	}

	interval := 2

	for _, jobTitle := range jobTitles {
		startIndex := jobTitle.Index - interval
		if startIndex < 0 {
			startIndex = 0
		}

		endIndex := jobTitle.Index + interval
		if endIndex >= len(sentences) {
			endIndex = len(sentences) - 1
		}

		dateFound := false
		for i := startIndex; i <= endIndex; i++ {
			if sentences[i].Date != nil &&
				sentences[i].Date.DateStart != "" {
				dateFound = true
				break
			}
		}

		if dateFound {
			titleSentences = append(titleSentences, sentences[jobTitle.Index])
		}
	}

	upperCasingTitles := 0
	for _, titleSentence := range titleSentences {
		if titleSentence.WordsCount == len(titleSentence.UpperCaseWords) {
			upperCasingTitles++
		}
	}

	experienceList = collectDataInRange(sentences, titleSentences)
	fmt.Println("experienceList: ", experienceList)

	// TODO добавить проверку на заголовки в верхнем регистре
	if upperCasingTitles == 0 {
	}

	if err != nil {
		fmt.Println("err: ", err)
	}

	return experienceList, nil
}
