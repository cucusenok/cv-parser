package CVParser

import (
	"cv-parser/parser"
	"cv-parser/parser/work_duration"
	"cv-parser/spell"
	"database/sql"
	"fmt"
	"math"
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
	regexPhone                *regexp.Regexp = regexp.MustCompile(`\+?\d{1,4}?[-.\s]?\(?\d{1,3}?\)?[-.\s]?\d{1,4}[-.\s]?\d{1,4}[-.\s]?\d{1,9}`)
	regexGithub               *regexp.Regexp = regexp.MustCompile(`https?://github.com/([a-zA-Z0-9._-]+)$`)
	regexEmail                *regexp.Regexp = regexp.MustCompile(`^[\w-\.]+@([\w-]+\.)+[\w-]{2,4}$`)
	regexSocials              *regexp.Regexp = regexp.MustCompile(`((https|http):\/\/)?[-a-zA-Z0-9@:%._\+~#=]{1,256}\.[a-zA-Z0-9()]{1,6}\b([-a-zA-Z0-9()@:%_\+.~#?&//=]*)`)
)

type AddressData struct {
	Country     string `json:"country"`
	City        string `json:"city"`
	State       string `json:"state"`
	CountryCode string `json:"country_code"`
	StateCode   string `json:"state_code"`
}

type EducationData struct {
	Date        *work_duration.WorkPeriod `json:"date"`
	Place       string                    `json:"place"`
	Levels      []string                  `json:"levels"`
	Description []string                  `json:"description"`
}

type ContactsData struct {
	Emails         []string    `json:"emails"`
	Github         string      `json:"github"`
	Address        AddressData `json:"address"`
	Phones         []string    `json:"phones"`
	SocialNetworks []string    `json:"social_networks"`
}

type ExperienceString struct {
	Date        *work_duration.WorkPeriod `json:"date"`
	Title       string                    `json:"title"`
	Positions   []string                  `json:"position"`
	Skills      []string                  `json:"skills"`
	Level       []string                  `json:"level"`
	Description []string                  `json:"description"`
}

type SentenceData struct {
	Index                     int                       `json:"index"`
	WordsCount                int                       `json:"words_count"`
	Commas                    int                       `json:"commas"`
	Words                     []string                  `json:"words"`
	Skills                    []string                  `json:"skills"`
	Positions                 []string                  `json:"position"`
	Level                     []string                  `json:"level"`
	Date                      *work_duration.WorkPeriod `json:"date"`
	Sentence                  string                    `json:"sentence"`
	UpperCaseWords            []string                  `json:"upper_case_words"`
	IsPossibleBullet          bool                      `json:"is_possible_bullet"`
	Github                    string                    `json:"github"`
	Phones                    []string                  `json:"phones"`
	SocialNetworks            []string                  `json:"social_networks"`
	Emails                    []string                  `json:"emails"`
	EducationLevels           []string                  `json:"education_levels"`
	EducationPlace            string                    `json:"education_place"`
	IsPossiblePartOfEducation bool                      `json:"is_possible_part_of_education"`
}

type CVData struct {
	JobTitle   string             `json:"job_title"`
	Contacts   ContactsData       `json:"contacts"`
	Skills     []string           `json:"skills"`
	Experience []ExperienceString `json:"experience"`
	Educations []EducationData    `json:"educations"`
}

var EducationHelpersKeys = []string{"GPA", "EDUCATION", "gpa", "education"}
var EducationPlaceKeys = []string{"university", "University", "Unversity", "unversity"}
var EducationLevelsKeys = []string{"master", "Master", "Bachelor", "bachelor", "MSc", "msc", "MS", "ms", "BSc", "bsc", "BS", "bs"}

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

// calculateAverageDiff вычисляет среднее значение между индексами в массиве.
func calculateAverageDiff(sentences []SentenceData) (int, error) {
	if len(sentences) < 2 {
		return 0, fmt.Errorf("недостаточно данных для вычисления среднего")
	}

	totalDiff := 0
	for i := 1; i < len(sentences); i++ {
		diff := sentences[i].Index - sentences[i-1].Index
		totalDiff += diff
	}

	averageDiff := float64(totalDiff) / float64(len(sentences)-1)
	roundedAverageDiff := int(math.Round(averageDiff)) // Округляем до ближайшего целого числа

	return roundedAverageDiff, nil
}

// filterByAverageDiff фильтрует массив на основе средней разницы между индексами.
func filterByAverageDiff(sentences []SentenceData) []SentenceData {
	averageDiff, averageDiffErr := calculateAverageDiff(sentences)

	if averageDiffErr != nil || len(sentences) < 2 {
		return sentences
	}

	var filteredSentences []SentenceData
	filteredSentences = append(filteredSentences, sentences[0]) // Добавляем первый элемент

	for i := 1; i < len(sentences); i++ {
		diff := sentences[i].Index - sentences[i-1].Index
		if diff <= averageDiff {
			filteredSentences = append(filteredSentences, sentences[i])
		}
	}

	return filteredSentences
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

/*
Функция findTitleSentences фильтрует возможные job titles и возвращает данные исходя их условий.
Считаем, что заголовки, которые возвращаются этой функцией - job title
*/
func findTitleSentences(possibleJobTitles []SentenceData, allSentences []SentenceData) []SentenceData {

	/*
	   Ожидается что дата для job title
	   JUNIOR FULL STACK DEVELOPER Jan, 2017 - June, 2019
	   будет либо в этой либо несколько строк выше/ниже
	   JUNIOR FULL STACK DEVELOPER
	   Jan, 2017 - June, 2019
	   Это параметр который определяет максимальное расстояние между строками
	*/
	allowedInterval := 2

	titleSentences := []SentenceData{}
	for _, possibleJobTitle := range possibleJobTitles {
		startIndex := possibleJobTitle.Index - allowedInterval
		if startIndex < 0 {
			startIndex = 0
		}

		endIndex := possibleJobTitle.Index + allowedInterval
		if endIndex >= len(allSentences) {
			endIndex = len(allSentences) - 1
		}

		dateFound := false
		for i := startIndex; i <= endIndex; i++ {
			if allSentences[i].Date != nil &&
				allSentences[i].Date.DateStart != "" {
				dateFound = true
				break
			}
		}

		if dateFound {
			titleSentences = append(titleSentences, allSentences[possibleJobTitle.Index])
		}
	}

	return titleSentences
}

// splitByAverageDiff делит массив на несколько подмассивов на основе среднего значения разницы между индексами.
func splitByAverageDiff(indices []SentenceData, averageDiff int) [][]SentenceData {
	if len(indices) < 2 {
		return [][]SentenceData{indices}
	}

	var result [][]SentenceData
	var currentSubArray []SentenceData

	currentSubArray = append(currentSubArray, indices[0])

	for i := 1; i < len(indices); i++ {
		diff := indices[i].Index - indices[i-1].Index
		if diff >= averageDiff {
			result = append(result, currentSubArray)
			currentSubArray = []SentenceData{}
		}
		currentSubArray = append(currentSubArray, indices[i])
	}

	// Добавляем последний подмассив в результат
	if len(currentSubArray) > 0 {
		result = append(result, currentSubArray)
	}

	return result
}

func collectEducationDataInRange(sentences, educationSentences []SentenceData) []EducationData {
	averageDiff, averageDiffErr := calculateAverageDiff(educationSentences)
	var collectedData []EducationData
	if averageDiffErr != nil {
		return collectedData
	}
	subArrays := [][]SentenceData{}

	if averageDiff > 1 {
		subArrays = splitByAverageDiff(educationSentences, averageDiff)
	} else {
		subArrays = splitByAverageDiff(educationSentences, len(educationSentences))
	}

	for i := 0; i < len(subArrays); i++ {
		startIndex := 0
		endIndex := 0
		if averageDiff == 1 {
			startIndex = subArrays[i][0].Index
			endIndex = subArrays[i][len(subArrays[i])-1].Index + 1
		} else if i < len(subArrays)-1 {
			startIndex = subArrays[i][0].Index
			endIndex = subArrays[i+1][0].Index
		} else {
			startIndex = subArrays[i][0].Index
			endIndex = subArrays[i][0].Index + averageDiff
		}

		educationData := EducationData{}
		for _, sentence := range sentences {
			if sentence.Index >= startIndex && sentence.Index < endIndex {
				educationData.Description = append(educationData.Description, sentence.Sentence)
				if sentence.Date.DateStart != "" {
					educationData.Date = sentence.Date
				}
				if len(sentence.EducationLevels) > 0 {
					for _, level := range sentence.EducationLevels {
						if !parser.ContainsItem(educationData.Levels, level) {
							educationData.Levels = append(educationData.Levels, level)
						}
					}
				}
				if len(sentence.EducationPlace) > 0 {
					educationData.Place = sentence.EducationPlace
				}
			}
		}

		collectedData = append(collectedData, educationData)
	}
	return collectedData
}

func containsWord(word string, keys []string) bool {
	for _, key := range keys {
		if word == key {
			return true
		}
	}
	return false
}

func ParseCV(text string) (CVData, error) {
	err := godotenv.Load()
	spellInstance, err = LoadSpellFromDB()
	if err != nil {
		return CVData{}, err
	}

	contactsData := ContactsData{}
	experienceList := []ExperienceString{}
	educations := []EducationData{}
	sentences := []SentenceData{}
	titleSentences := []SentenceData{}

	addressData := AddressData{
		Country:     "",
		City:        "",
		CountryCode: "",
		State:       "",
		StateCode:   "",
	}

	cvData := text

	possibleJobTitles := []SentenceData{}
	possiblePartsOfEducation := []SentenceData{}
	sentencesWithDates := []SentenceData{}
	skills := []string{}

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
		educationLevels := []string{}
		educationPlace := ""
		upperCaseWords := []string{}
		combinations := parser.GenerateCombinations(strings.Split(paragraph, " "))
		words := strings.Fields(specialCharsetRegex.ReplaceAllString(paragraph, ""))
		isPossibleBullet := false
		isPossiblePartOfEducation := false

		if strings.HasPrefix(paragraph, "•") {
			isPossibleBullet = true
		}

		// TODO после добавления ключевых слов в бд, перенести эту часть кода в цикл по комбинациям (combinations)
		for _, word := range words {
			if containsWord(word, EducationLevelsKeys) && !parser.ContainsItem(educationLevels, word) {
				educationLevels = append(educationLevels, word)
				isPossiblePartOfEducation = true
			}
			if containsWord(word, EducationPlaceKeys) {
				educationPlace = paragraph
				isPossiblePartOfEducation = true
			}
			//if containsWord(word, EducationHelpersKeys) {
			//	isPossiblePartOfEducation = true
			//}
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

		github := regexGithub.FindString(paragraph)
		emails := regexEmail.FindAllString(paragraph, -1)
		phones := regexPhone.FindAllString(paragraph, -1)

		//TODO некорректно собирает соц сети. В срез записывается Vue.js, Node.js, Bike.net:
		socialNetworks := regexSocials.FindAllString(paragraph, -1)

		if len(date.DateStart) > 0 {
			isPossiblePartOfEducation = true
		}

		sentences = append(sentences, SentenceData{
			Index:                     index,
			WordsCount:                len(strings.Fields(specialCharsetRegex.ReplaceAllString(paragraph, ""))),
			Words:                     strings.Fields(specialCharsetRegex.ReplaceAllString(paragraph, "")),
			Commas:                    strings.Count(paragraph, ","),
			Skills:                    skills,
			Positions:                 positions,
			Level:                     levels,
			Date:                      date,
			Sentence:                  paragraph,
			UpperCaseWords:            upperCaseWords,
			IsPossibleBullet:          isPossibleBullet,
			Github:                    github,
			Emails:                    emails,
			Phones:                    phones,
			SocialNetworks:            socialNetworks,
			EducationLevels:           educationLevels,
			EducationPlace:            educationPlace,
			IsPossiblePartOfEducation: isPossiblePartOfEducation,
		})
	}

	averageWordsCount = calculateAverageLength(sentences)

	for _, sentence := range sentences {
		if sentence.WordsCount == 0 {
			continue
		}

		for _, skill := range sentence.Skills {
			if !parser.ContainsItem(skills, skill) {
				skills = append(skills, skill)
			}
		}

		if len(sentence.Github) > 0 {
			contactsData.Github = sentence.Github
		}

		if len(sentence.Emails) > 0 {
			contactsData.Emails = append(contactsData.Emails, sentence.Emails...)
		}

		if len(sentence.Phones) > 0 {
			contactsData.Phones = append(contactsData.Phones, sentence.Phones...)
		}

		for _, socialNetwork := range sentence.SocialNetworks {
			if !parser.ContainsItem(contactsData.SocialNetworks, socialNetwork) {
				contactsData.SocialNetworks = append(contactsData.SocialNetworks, socialNetwork)
			}
		}

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
			(sentence.WordsCount <= averageWordsCount ||
				(sentence.WordsCount > maxWordCountParam) && (payloadPercent > 80) && (sentence.WordsCount-averageWordsCount < 6)) {
			jobTitle := sentence
			possibleJobTitles = append(possibleJobTitles, jobTitle)
		}

		if len(sentence.Date.DateStart) > 0 {
			sentencesWithDates = append(sentencesWithDates, sentence)
		}

		// Параметр, отвечающий за максимально допустимое расстояние между строками с информацией об обучении.
		// В основном, помогает отфильтровать даты из различных блоков
		maxIndexGapInEducationSentence := 20

		/*
			Условие для строк, которые могут быть частью информации об обучении.
			Условия:
				1) sentence.IsPossiblePartOfEducation в значении true
				2) Считаю, что кол-во слов в строке должно быть меньше либо равно среднему значению.
				3) В строке не находится информации о sentence.Positions.
				4) GAP между индексами строк не больше maxIndexGapInEducationSentence
		*/

		if sentence.IsPossiblePartOfEducation &&
			(sentence.WordsCount <= averageWordsCount || len(sentence.Date.DateStart) > 0) &&
			len(sentence.Positions) == 0 &&
			(len(possiblePartsOfEducation) == 0 || sentence.Index-possiblePartsOfEducation[len(possiblePartsOfEducation)-1].Index <= maxIndexGapInEducationSentence) {
			possiblePartsOfEducation = append(possiblePartsOfEducation, sentence)
		}
	}

	possibleEducationSentences := filterByAverageDiff(possiblePartsOfEducation)
	titleSentences = findTitleSentences(possibleJobTitles, sentences)

	upperCasingTitles := 0
	for _, titleSentence := range titleSentences {
		if titleSentence.WordsCount == len(titleSentence.UpperCaseWords) {
			upperCasingTitles++
		}
	}

	experienceList = collectDataInRange(sentences, titleSentences)
	// TODO необходимо научиться нормально определять конец блока, для того, чтобы не собирать лишнюю инфу в description. Это необходимо и для Experience
	educations = collectEducationDataInRange(sentences, possibleEducationSentences)

	//TODO собрать данные в addressData
	contactsData.Address = addressData

	// TODO в parsedCV добавить JobTitle
	parsedCV := CVData{
		Experience: experienceList,
		Contacts:   contactsData,
		Skills:     skills,
		Educations: educations,
		//JobTitle:
	}

	return parsedCV, nil
}
