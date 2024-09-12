package CVParser

import (
	"cv-parser/parser"
	"cv-parser/parser/work_duration"
	"cv-parser/spell"
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
	specialCharsetRegex       *regexp.Regexp = regexp.MustCompile("[^а-яА-Яa-zA-Z0-9 ]+")
)

type ExperienceString struct {
	Date        *work_duration.WorkPeriod `json:"date"`
	Title       string                    `json:"title"`
	Positions   []string                  `json:"position"`
	Skills      []string                  `json:"skills"`
	Level       []string                  `json:"level"`
	Description []string                  `json:"description"`
}

type SentenceData struct {
	Index            int                       `json:"index"`
	WordsCount       int                       `json:"wordsCount"`
	Commas           int                       `json:"commas"`
	Words            []string                  `json:"words"`
	Skills           []string                  `json:"skills"`
	Positions        []string                  `json:"position"`
	Level            []string                  `json:"level"`
	Date             *work_duration.WorkPeriod `json:"date"`
	Sentence         string                    `json:"sentence"`
	UpperCaseWords   []string                  `json:"upperCaseWords"`
	IsPossibleBullet bool                      `json:"isPossibleBullet"`
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
			} else {
				description = append(description, sentence.Sentence)
			}

			if sentence.Date.DateStart != "" {
				date = sentence.Date
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
				Description: description,
			}
		}
		result = append(result, rangeData)
	}

	return result
}

func IsAllowedDistanceForWord(suggest spell.Suggestion) bool {
	length := len(suggest.Word)
	if length < 6 {
		return suggest.Distance == 0
	}
	if length < 10 {
		return suggest.Distance < 2
	}
	return suggest.Distance < 3
}



["1", "2", "3"]
[ "1,2", "1,3", "1,2,3" ]

func ParseCV(text string) ([]ExperienceString, error) {
	err := godotenv.Load()
	spellInstance, err = LoadSpellFromDB()
	if err != nil {
		return nil, err
	}
	experienceList := []ExperienceString{}
	sentences := []SentenceData{}
	titleSentences := []SentenceData{}

	cvData := text

	possibleJobTitles := []SentenceData{}
	sentencesWithDates := []SentenceData{}

	var averageWordsCount int

	paragraphs := strings.Split(cvData, "\n")
	var filteredParagraphs []string

	for _, rawParagraph := range paragraphs {
		/*
			При копировании из PDF в конце строк может отображаться вместо точек "." -> "1"
		*/
		paragraph := regexCorrectEndOfSentence.ReplaceAllString(rawParagraph, "$1.")

		trimmedParagraph := strings.TrimSpace(paragraph) // Убираем лишние пробелы
		if trimmedParagraph != "" {
			filteredParagraphs = append(filteredParagraphs, trimmedParagraph)
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
		words := strings.Fields(specialCharsetRegex.ReplaceAllString(paragraph, ""))
		isPossibleBullet := false

		if strings.HasPrefix(paragraph, "•") {
			isPossibleBullet = true
		}

		for _, combination := range combinations {
			list, _ := spellInstance.Lookup(strings.ToLower(combination), spell.SuggestionLevel(spell.LevelClosest))

			if len(list) > 0 && !IsAllowedDistanceForWord(list[0]) {
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
			Index:            index,
			WordsCount:       len(strings.Fields(specialCharsetRegex.ReplaceAllString(paragraph, ""))),
			Words:            strings.Fields(specialCharsetRegex.ReplaceAllString(paragraph, "")),
			Commas:           strings.Count(paragraph, ","),
			Skills:           skills,
			Positions:        positions,
			Level:            levels,
			Date:             date,
			Sentence:         paragraph,
			UpperCaseWords:   upperCaseWords,
			IsPossibleBullet: isPossibleBullet,
		})
	}

	averageWordsCount = calculateAverageLength(sentences)

	for _, sentence := range sentences {
		if sentence.WordsCount == 0 {
			continue
		}
		jobTitle := SentenceData{}

		// Процент полезной нагрузки в предложении
		payloadPercent := (len(sentence.Skills) + len(sentence.Positions) + len(sentence.Level))
		if sentence.Date.DateStart != "" || sentence.Date.DateEnd != "" {
			// Добавим больше вероятности строкам с датами
			payloadPercent = payloadPercent + len(strings.Split(sentence.Date.DateStart, ".")) + len(strings.Split(sentence.Date.DateEnd, "."))
		}
		payloadPercent = (payloadPercent * 100) / sentence.WordsCount
		maxWordCountParam := 7

		// Если в строке Skills + Positions + Level
		// занимают >= 50% строки и короче средней длины строки
		// занимает 80% строки, но больше maxWordCountParam и тогда может быть длиннее averageWordsCount
		// тогда считаем, что это может быть jobTitle
		if len(sentence.Positions) > 0 &&
			(payloadPercent >= 50) &&
			!sentence.IsPossibleBullet &&
			(sentence.WordsCount < averageWordsCount ||
				(sentence.WordsCount > maxWordCountParam) && (payloadPercent > 80) && (sentence.WordsCount-averageWordsCount < 6)) {
			jobTitle = sentence
			possibleJobTitles = append(possibleJobTitles, jobTitle)
		}
		if len(sentence.Date.DateStart) > 0 {
			sentencesWithDates = append(sentencesWithDates, sentence)
		}
	}

	/*
	   Ожидается что дата для job title
	   JUNIOR FULL STACK DEVELOPER Jan, 2017 - June, 2019
	   будет либо в этой либо несколько строк выше/ниже
	   JUNIOR FULL STACK DEVELOPER
	   Jan, 2017 - June, 2019
	   Это параметр который определяет максимальное расстояние между строками
	*/
	allowedInterval := 2

	for _, possibleJobTitle := range possibleJobTitles {
		startIndex := possibleJobTitle.Index - allowedInterval
		if startIndex < 0 {
			startIndex = 0
		}

		endIndex := possibleJobTitle.Index + allowedInterval
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
			titleSentences = append(titleSentences, sentences[possibleJobTitle.Index])
		}
	}

	upperCasingTitles := 0
	for _, titleSentence := range titleSentences {
		if titleSentence.WordsCount == len(titleSentence.UpperCaseWords) {
			upperCasingTitles++
		}
	}

	experienceList = collectDataInRange(sentences, titleSentences)

	// TODO добавить проверку на заголовки в верхнем регистре
	if upperCasingTitles != 0 {
	}

	if err != nil {
		fmt.Println("err: ", err)
	}

	return experienceList, nil
}
