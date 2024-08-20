package work_duration

import (
	"fmt"
	"math"
	"regexp"
	"strconv"
	"strings"
)

var (
	regexYear                *regexp.Regexp = regexp.MustCompile(`(19|20)\d{2}`)
	regexDateRange           *regexp.Regexp = regexp.MustCompile(`(19|20)\d{2}.*(19|20)\d{2}`)
	regexDateRangeExcludeEnd *regexp.Regexp = regexp.MustCompile(`(19|20)\d{2}(.*(19|20)\d{2})?`)
	regexMonthYear           *regexp.Regexp = regexp.MustCompile(`(?P<month>\d{1,2})[-./_](?P<year>\d{4})`)
	regexFullDate            *regexp.Regexp = regexp.MustCompile(`(?:(?P<day>\d{1,2})\D+)?(?:(?P<month>\d{1,2})\D+)?(?P<year>\d{4})`)
	regexFullDateReverse     *regexp.Regexp = regexp.MustCompile(`^(?P<year>\d{4})\D+(?P<month>\d{1,2})(\D+(?P<day>\d{1,2}))?$`)
	regexNonDigit            *regexp.Regexp = regexp.MustCompile(`\D`)
)

type WorkPeriod struct {
	DateStart string `json:"start_date"`
	DateEnd   string `json:"end_date"`
}

/*
	isValidDate проверяет дату на корректность.
	month >= 1 && month <= 12
	day >= 1 && day <= 31
	TODO обработать некорректные даты
*/

func isValidDate(value string) bool {
	date := regexFullDate.FindStringSubmatch(value)
	if len(date) == 0 {
		return false
	}

	var month, day int
	var err error

	if len(date) == 3 {
		month, err = strconv.Atoi(date[1])
		day, err = strconv.Atoi(date[0])
	} else if len(date) == 2 {
		month, err = strconv.Atoi(date[0])
	}

	if err != nil {
		return false
	}

	isWithoutDays := day == 0
	onlyYear := month == 0 && day == 0

	if onlyYear {
		return true
	} else if isWithoutDays {
		return month >= 1 && month <= 12
	} else {
		return (month >= 1 && month <= 12) || (day >= 1 && day <= 31)
	}
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
	formatDate форматирует дату, заменяя любые знаки между элементами даты на "." .
	Примеры formatDate:

	•	"11/12/1923" -> "11.12.1923"
	•	"11👉12👉1923" -> "11.12.1923"
	•	"11ю12ю1923" -> "11.12.1923"
	•	"11/1925" -> "11.1925"
	•	"11👉1925" -> "11.1925"
	•	"1925" -> "1925"

*/

func formatDate(value string) string {
	result := regexNonDigit.ReplaceAllStringFunc(value, func(s string) string {
		return "."
	})
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

	if !isValidDate(startDate) {
		startDate = "Invalid start date"
	}

	endDate := fmt.Sprintf("%v", math.Inf(1))

	if len(result) >= 2 {
		endDate = formatDate(result[1])

		if !isValidDate(endDate) {
			endDate = "Invalid end date"
		}
	}

	return WorkPeriod{
		DateStart: startDate,
		DateEnd:   endDate,
	}
}
