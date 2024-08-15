package spell

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"sort"
	"strconv"
	"strings"
	"testing"
)

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

func Save() {
	s := New()
	db, err := sql.Open("postgres", "user=crawler dbname=andrey_words sslmode=disable password=crawler host=pg.dev.local port=5432")
	if err != nil {
		return
	}
	fmt.Println("loading...")
	if err := scanQuery(db, `
			with use as (select unnest(words) as id, count(*) as cnt from entities group by 1)
		select
			words.id,
			words.word,
			use.cnt
		from words
		left join use on use.id = words.id
		`, nil, func(rows *sql.Rows) error {
		var id int32
		var word string
		var count uint64
		if err := rows.Scan(&id, &word, &count); err != nil {
			return err
		}

		s.AddEntry(Entry{
			Frequency: count,
			Word:      word,
			WordData: WordData{
				"id": id,
			},
		})

		return nil
	}); err != nil {
		return
	}
	fmt.Println("saving...")
	s.Save("spell_data")
}

type WordStruct struct {
	Id int `json:"id"`
}

func ArraySort(arr []string) []string {
	sort.SliceStable(arr, func(i, j int) bool {
		val1, _ := strconv.Atoi(arr[i])
		val2, _ := strconv.Atoi(arr[j])
		return val1 < val2
	})

	return arr
}

func GenerateCombinations(arrays [][]string, currentCombo []string, index int, combinations *[][]string) {
	if index == len(arrays) {
		if len(currentCombo) > 0 {
			combinationCopy := make([]string, len(currentCombo))
			copy(combinationCopy, currentCombo)
			combinationCopy = ArraySort(combinationCopy)
			*combinations = append(*combinations, combinationCopy)
		}
		return
	}

	for _, elem := range arrays[index] {
		if len(currentCombo) == 3 {
			return
		}
		GenerateCombinations(arrays, append(currentCombo, elem), index+1, combinations)
	}
	GenerateCombinations(arrays, currentCombo, index+1, combinations)
}

func Test_Match(t *testing.T) {
	var combinations [][]string
	arrays := [][]string{
		[]string{"1", "2", "3", "4", "5"},
		[]string{"1", "2", "3", "4", "5"},
		[]string{"1", "2", "3", "4", "5"},
		[]string{"1", "2", "3", "4", "5"},
		[]string{"1", "2", "3", "4", "5"},
		[]string{"1", "2", "3", "4", "5"},
		[]string{"1", "2", "3", "4", "5"},
		[]string{"1", "2", "3", "4", "5"},
		[]string{"1", "2", "3", "4", "5"},
	}

	GenerateCombinations(arrays, []string{}, 0, &combinations)
	combinations_strs := []string{}
	for _, combination := range combinations {
		combination_str := fmt.Sprintf("'%s'", strings.Join(combination, ","))
		combinations_strs = append(combinations_strs, combination_str)
	}
	fmt.Println(len(combinations_strs))
	//fmt.Println(combinations_strs)
	return

	//s := New()
	//filename := "spell_data"
	/*	s, err := Load(filename)
		if err != nil {
			return
		}

		q := []string{"new", "yorl", "city"}
		arrays := [][]string{}
		for _, word := range q {
			array := []string{}
			prt := fmt.Sprintf("======= %s =====", word)
			fmt.Println(prt)
			start := time.Now()
			list, _ := s.Lookup(word, SuggestionLevel(LevelClosest))
			fmt.Println(word, time.Since(start))
			for _, item := range list {
				id := strconv.Itoa(int(item.WordData["id"].(float64)))
				array = append(array, id)
			}
			arrays = append(arrays, array)
		}

		var combinations [][]string
		PostgresWordCollector.GenerateCombinations(arrays, []string{}, 0, &combinations)
		combinations_strs := []string{}
		for _, combination := range combinations {
			combination_str := fmt.Sprintf("'%s'", strings.Join(combination, ","))
			combinations_strs = append(combinations_strs, combination_str)
		}*/
}
