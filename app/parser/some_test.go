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

	/*	list1, _ := spellInst.Lookup("data science", spell.SuggestionLevel(spell.LevelClosestAndPossibleMatches))

		fmt.Println(list1)*/

	//words := []string{"data science", "data"}
	for _, v := range spellInst.Words {
		list, _ := spellInst.Lookup(v.Word, spell.SuggestionLevel(spell.LevelClosestAndPossibleMatches))
		if len(list) == 0 {
			continue
		}
		if len(list) > 0 && !IsAllowedDistanceForWord(list[0]) {
			continue
		}
		// TODO: пропускать слова в которых не совпадает первая буква
		fmt.Println(v.Word)
		for _, l := range list {
			fmt.Println(l.Word, l.Frequency, l.Distance)
		}

		fmt.Println("========================================================")
	}
}
