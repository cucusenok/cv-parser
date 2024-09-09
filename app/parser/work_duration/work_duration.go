package work_duration

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

var (
	RegexDateRangeExcludeEnd *regexp.Regexp = regexp.MustCompile(`(19|20)\d{2}(.*(19|20)\d{2})?`)
	regexFullDate            *regexp.Regexp = regexp.MustCompile(`(?:(?P<day>\d{1,2})\D{1})?(?:(?P<month>\d{1,2}|(?i)[a-zа-я]+)\D{1})?(?P<year>\d{4})`)
	RegexDates               *regexp.Regexp = regexp.MustCompile(`(?P<reverseDate>(19|20)\d{2}(\D{1,3}\d{1,2}|(?i)[a-zа-я]+)?(\D{1,3}\d{1,2})?)(\D|$)|(?P<normalDate>(\d{1,2}\D{1,3})?((\d{1,2}|(?i)[a-zа-я]+)\D{1,3})?(19|20)\d{2})`)
	regexFullDateReverse     *regexp.Regexp = regexp.MustCompile(`(?P<year>(19|20)\d{2})(\D{1}(?P<month>(\d{1,2}|(?i)[a-zа-я]+)))(\D{1}(?P<day>\d{1,2}))?`)
	regexNonDigit            *regexp.Regexp = regexp.MustCompile(`\D`)
	monthPattern             *regexp.Regexp = regexp.MustCompile(`(?i)[a-zа-я]+`)
	RegexDatesWithoutDigit   *regexp.Regexp = regexp.MustCompile(`(?P<reverseDate>(19|20)\d{2}(\D{1,3}\d{1,2})?(\D{1,3}\d{1,2})?)(\D|$)|(?P<normalDate>(\d{1,2}\D{1,3})?((\d{1,2})\D{1,3})?(19|20)\d{2})`)
)

type WorkPeriod struct {
	DateStart string `json:"start_date"`
	DateEnd   string `json:"end_date"`
}

type Month struct {
	MonthRoot string `json:"month"`
	Index     int    `json:"index"`
}

var months = []Month{
	{MonthRoot: "jan", Index: 1},
	{MonthRoot: "feb", Index: 2},
	{MonthRoot: "mar", Index: 3},
	{MonthRoot: "apr", Index: 4},
	{MonthRoot: "may", Index: 5},
	{MonthRoot: "jun", Index: 6},
	{MonthRoot: "jul", Index: 7},
	{MonthRoot: "aug", Index: 8},
	{MonthRoot: "sep", Index: 9},
	{MonthRoot: "oct", Index: 10},
	{MonthRoot: "nov", Index: 11},
	{MonthRoot: "dec", Index: 12},
	{MonthRoot: "янв", Index: 1},
	{MonthRoot: "фев", Index: 2},
	{MonthRoot: "мар", Index: 3},
	{MonthRoot: "апр", Index: 4},
	{MonthRoot: "май", Index: 5},
	{MonthRoot: "июн", Index: 6},
	{MonthRoot: "июл", Index: 7},
	{MonthRoot: "авг", Index: 8},
	{MonthRoot: "сен", Index: 9},
	{MonthRoot: "окт", Index: 10},
	{MonthRoot: "ноя", Index: 11},
	{MonthRoot: "дек", Index: 12},
}

const TYPE_DATE_NORMAL = "normalDate"   // 20.12.2024
const TYPE_DATE_REVERSE = "reverseDate" // 2024.12.20

/*
separateDates Функция разделяющая даты.

•	Примеры:
•	"1998.1.11_2000.11.11" => ["1998.1.11", "2000.11.11"]
•	"1998.2.12_2000.11" => ["1998.2.12", "2000.11"]
•	"1998.3.13_2000" => ["1998.3.13", "13_2000"]
•	"11.1998_2000.11.11" => ["11.1998", "2000.11.11"]
•	"1998 2000/11/11" => ["1998", "2000/11/11"]

и тд
*/
func SeparateDates(date string) []string {
	matches := RegexDates.FindAllStringSubmatch(date, -1)
	// тут мы ожидаем группы normal или reverse
	// и если одна из этих груп найдена - разбиваем по разделителю
	groupNames := RegexDates.SubexpNames()
	result := []string{}

	for _, match := range matches {
		for i, name := range groupNames {
			if len(name) == 0 {
				continue
			}
			if (name == TYPE_DATE_REVERSE || name == TYPE_DATE_NORMAL) && i < len(match) && len(match[i]) > 0 {
				result = append(result, match[i])
			}
		}
	}

	return result
}

/*
convertMonthToNumber Функция для замены текстового месяца на числовой.
Если совпадений не найдено, возвращаем 0.

•	Примеры строковых дат:
•	"Ноября" => "11"
•	"November" => "11"
•	"Nov" => "11"

и тд
*/
func convertMonthToNumber(month string) int {
	month = strings.ToLower(month)
	for _, m := range months {
		if strings.Contains(month, m.MonthRoot) {
			return m.Index
		}
	}
	return 0
}

/*
ReplaceDateWithMonthNumber Функция для замены текстового месяца на числовой.
Если совпадений не найдено, возвращаем исходную строку.

•	Примеры строковых дат:
•	"1 Ноября 1923" => "1 11 1923"
•	"Ноябрь 1923" => "11 1923"
•	"1 November 1923" => "1 11 1923"
•	"Nov 1923" => "11 1923"
*/
func ReplaceDateWithMonthNumber(date string) string {
	replacedStr := monthPattern.ReplaceAllStringFunc(date, func(match string) string {
		index := strconv.Itoa(convertMonthToNumber(match))
		if index != "0" {
			res := formatDayOrMonth(index)
			return res
		}
		return match
	})

	return replacedStr
}

/*
	reverseDate возвращает дату в формате dd.mm.yyyy.
	Примеры reverseDate:

	2022/09/01 => 01.09.2022
	2022.09.01 => 01.09.2022
	2022.09 => 09.2022
*/

func reverseDate(date string) string {
	match := regexFullDateReverse.FindAllStringSubmatch(date, -1)
	groupNames := regexFullDateReverse.SubexpNames()

	var year, month, day string

	for _, m := range match {
		for i, name := range groupNames {
			if name == "year" && i < len(m) {
				year = m[i]
			} else if name == "month" && i < len(m) {
				month = m[i]
			} else if name == "day" && i < len(m) {
				day = m[i]
			}
		}
	}

	if len(match) == 0 {
		return date
	}

	isWithoutDays := len(day) == 0
	onlyYear := len(month) == 0 && len(day) == 0

	switch {
	case onlyYear:
		return year
	case isWithoutDays:
		return fmt.Sprintf("%s.%s", formatDayOrMonth(month), year)
	default:
		return fmt.Sprintf("%s.%s.%s", formatDayOrMonth(day), formatDayOrMonth(month), year)
	}
}

/*
	formatDayOrMonth добавляет 0 в начало, если это требуется.
	Примеры formatDayOrMonth:

	1.09.2022 => 01.09.2022
	01.9.2022 => 01.09.2022
	1.9.2022 => 01.09.2022

*/

func formatDayOrMonth(value string) string {
	if len(value) == 1 {
		return "0" + value
	}
	return value
}

/*
	FormatDate форматирует дату, заменяя любые знаки между элементами даты на "." .
	Примеры formatDate:

	•	"1 Ноября 1923" -> "01.11.1923"
	•	"январь 2033" -> "01.2023"
	•	"11/12/1923" -> "11.12.1923"
	•	"11👉12👉1923" -> "11.12.1923"
	•	"11ю12ю1923" -> "11.12.1923"
	•	"11/1925" -> "11.1925"
	•	"11👉1925" -> "11.1925"
	•	"1925" -> "1925"

*/

func FormatDate(value string) string {
	formattedDate := ReplaceDateWithMonthNumber(value)
	formattedDate = regexNonDigit.ReplaceAllStringFunc(formattedDate, func(s string) string {
		return "."
	})
	date := strings.Split(formattedDate, ".")
	result := []string{}

	for _, d := range date {
		if d != "" {
			result = append(result, d)
		}
	}

	var year string
	var month string
	var day string

	if len(result) == 3 {
		year = result[2]
		month = result[1]
		day = result[0]
	} else if len(result) == 2 {
		year = result[1]
		month = result[0]
	} else {
		year = result[0]
	}

	isWithoutDays := len(day) == 0
	onlyYear := len(month) == 0 && len(day) == 0

	if onlyYear {
		return year
	} else if isWithoutDays {
		return fmt.Sprintf("%s.%s", formatDayOrMonth(month), year)
	}
	return fmt.Sprintf("%s.%s.%s", formatDayOrMonth(day), formatDayOrMonth(month), year)
}

func IsMatchDate(text string) bool {
	return RegexDateRangeExcludeEnd.MatchString(text)
}

func isReverseDate(text string) bool {
	return regexFullDateReverse.MatchString(text)
}

func ParsePeriod(text string) (*WorkPeriod, error) {
	date := strings.TrimSpace(text)
	if date == "" || !IsMatchDate(date) {
		return &WorkPeriod{
			DateStart: "",
			DateEnd:   "",
		}, nil
	}

	separatedDates := SeparateDates(date)
	var result []string
	/*
		Если есть возможность привести дату к единому формату - это нужно делать чтобы дальше
		в коде не плодить isReverseDate(datePart)
	*/
	for _, datePart := range separatedDates {
		datePart = FormatDate(datePart)
		isReverse := isReverseDate(datePart)
		if isReverse {
			datePart = reverseDate(datePart)
		}
		result = append(result, datePart)
	}

	if len(result) == 0 {
		return &WorkPeriod{
			DateStart: "",
			DateEnd:   "",
		}, nil
	}

	startDate := result[0]
	endDate := ""

	if len(result) > 2 {
		// Тут нужна с-ма логирования, чтобы увидеть, какой пользовательский ввод поломал софт

		resStr := fmt.Sprintf("Text pasrsin failed: %s", text)
		fmt.Println(resStr)

		return nil, errors.New(resStr)
	}

	if len(result) == 2 {
		endDate = result[1]
	}

	return &WorkPeriod{
		DateStart: startDate,
		DateEnd:   endDate,
	}, nil
}
