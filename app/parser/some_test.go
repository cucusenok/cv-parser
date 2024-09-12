package parser

import (
	"cv-parser/spell"
	"fmt"
	"testing"
)

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

func TestSomething(test *testing.T) {
	spellInst, err := LoadSpellFromDB()
	if err != nil {
		panic(err)
	}

	processedWords := make(map[string]bool)
	var groupsList [][]string
	for _, v := range spellInst.Words {
		if processedWords[v.Word] {
			continue
		}

		list, _ := spellInst.Lookup(v.Word, spell.SuggestionLevel(spell.LevelClosestAndPossibleMatches))
		if len(list) == 0 {
			continue
		}
		if len(list) > 0 && !IsAllowedDistanceForWord(list[0]) {
			continue
		}

		// TODO: пропускать слова в которых не совпадает первая буква
		group := []string{v.Word}
		processedWords[v.Word] = true

		for _, l := range list {
			if l.Word[:2] != v.Word[:2] {
				continue
			}

			if !processedWords[l.Word] {
				group = append(group, l.Word)
				processedWords[l.Word] = true // Отмечаем найденное слово как обработанное
			}
		}
		if len(group) >= 2 {
			groupsList = append(groupsList, group)
			fmt.Println(group)
		}

		fmt.Println("========================================================")
	}
}
