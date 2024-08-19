package work_duration

import (
	"fmt"
	"math"
	"regexp"
	"strings"
)

var (
	regexYear                *regexp.Regexp = regexp.MustCompile(`(19|20)\d{2}`)
	regexDateRange           *regexp.Regexp = regexp.MustCompile(`(19|20)\d{2}.*(19|20)\d{2}`)
	regexDateRangeExcludeEnd *regexp.Regexp = regexp.MustCompile(`(19|20)\d{2}(.*(19|20)\d{2})?`)
	regexMonthYear           *regexp.Regexp = regexp.MustCompile(`(?P<month>\d{1,2})[-./_](?P<year>\d{4})`)
	regexFullDate            *regexp.Regexp = regexp.MustCompile(`(?:(?P<day>\d{1,2})\D+)?(?:(?P<month>\d{1,2})\D+)?(?P<year>\d{4})`)
	regexFullDateReverse     *regexp.Regexp = regexp.MustCompile(`^(?P<year>\d{4})\D+(?P<month>\d{1,2})(\D+(?P<day>\d{1,2}))?$`)
	regexRemoveSymbolsOld2   *regexp.Regexp = regexp.MustCompile(`^(?:(?P<day>\d{1,2})\D+)?(?P<month>\d{1,2})?\D*(?P<year>\d{4})$`)
	regexRemoveSymbolsOld    *regexp.Regexp = regexp.MustCompile(`^(?:(\d{1,2})\D)?(?:(\d{1,2})\D)?(\d{4})$`)
	regexRemoveSymbols       *regexp.Regexp = regexp.MustCompile(`(\d+)\D+(\d+)`)
)

type WorkPeriod struct {
	DateStart string `json:"start_date"`
	DateEnd   string `json:"end_date"`
}

/*
	reverseDate возвращает дату в формате dd.mm.yyyy.
	Примеры reverseDate:

	2022.09.01 => 01.09.2022
	2022.09 => 09.2022
*/

func reverseDate(date string) string {
	match := regexFullDateReverse.FindStringSubmatch(date)
	if len(match) == 0 {
		return date
	}

	year := match[1]
	month := match[2]
	day := match[3]
	isWithoutDays := len(day) == 0
	onlyYear := len(month) == 0 && len(day) == 0

	switch {
	case onlyYear:
		return year
	case isWithoutDays:
		return fmt.Sprintf("%s.%s", month, year)
	default:
		return fmt.Sprintf("%s.%s.%s", day, month, year)
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
	formatDayOrMonth добавляет 0 в начало, если это требуется.
	Примеры formatDayOrMonth:

	1.09.2022 => 01.09.2022
	01.9.2022 => 01.09.2022
	1.9.2022 => 01.09.2022

*/

func padZeroToDate(date []string) string {
	var year string
	var month string
	var day string

	if len(date) == 3 {
		year = date[2]
		month = date[1]
		day = date[0]
	} else if len(date) == 2 {
		year = date[1]
		month = date[0]
	} else {
		year = date[0]
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
	formatDate форматирует дату, заменяя любые знаки между элементами даты на "." .
	Примеры formatDate:

	•	"11/12/1923" -> "11.12.1923"
	•	"11👉12👉1923" -> "11.12.1923"
	•	"11ю12ю1923" -> "11.12.1923"
	•	"11/1925" -> "11.1925"
	•	"11👉1925" -> "11.1925"
	•	"1925" -> "1925"

*/

func formatDateOld(value string) string {
	matches := regexRemoveSymbols.FindStringSubmatch(value)
	if len(matches) == 0 {
		return value
	}
	var year string
	var month string
	var day string

	date := matches[1:]

	if len(date) == 3 {
		year = date[2]
		month = date[1]
		day = date[0]
	} else if len(date) == 2 {
		year = date[1]
		month = date[0]
	} else {
		year = date[0]
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

func formatDate(value string) string {
	result := regexRemoveSymbols.ReplaceAllString(value, "$1.$2")
	date := strings.Split(result, ".")

	var year string
	var month string
	var day string

	if len(date) == 3 {
		year = date[2]
		month = date[1]
		day = date[0]
	} else if len(date) == 2 {
		year = date[1]
		month = date[0]
	} else {
		year = date[0]
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

func isMatchDate(text string) bool {
	return regexDateRangeExcludeEnd.MatchString(text)
}

func isReverseDate(text string) bool {
	return regexFullDateReverse.MatchString(text)
}

func reformatPeriod(text string) WorkPeriod {
	date := strings.TrimSpace(text)
	if date == "" || !isMatchDate(date) {
		return WorkPeriod{
			DateStart: "",
			DateEnd:   "",
		}
	}

	result := []string{}
	isReverse := isReverseDate(date)

	if isReverse {
		date := reverseDate(date)
		result = regexFullDate.FindAllString(date, -1)
	} else {
		result = regexFullDate.FindAllString(date, -1)
	}

	startDate := formatDate(result[0])
	endDate := fmt.Sprintf("%v", math.Inf(1))

	if len(result) >= 2 {
		endDate = formatDate(result[1])
	}

	return WorkPeriod{
		DateStart: startDate,
		DateEnd:   endDate,
	}
}
