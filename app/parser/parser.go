package parser

import (
	"awesomeProject2/app/spell"
	"database/sql"
	"fmt"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"regexp"
	"strings"
)

var spellInstance *spell.Spell

var (
	regexZipCode                 *regexp.Regexp = regexp.MustCompile(`\d{4,5}(?:[-\s]\d{4})?`)
	regexPhone                   *regexp.Regexp = regexp.MustCompile(`\+?\d{1,4}?[-.\s]?\(?\d{1,3}?\)?[-.\s]?\d{1,4}[-.\s]?\d{1,4}[-.\s]?\d{1,9}`)
	regexEmail                   *regexp.Regexp = regexp.MustCompile("[a-z0-9!#$%&'*+/=?^_`{|}~-]+(?:\\.[a-z0-9!#$%&'*+/=?^_`{|}~-]+)*@(?:[a-z0-9](?:[a-z0-9-]*[a-z0-9])?\\.)+[a-z0-9](?:[a-z0-9-]*[a-z0-9])?")
	regexDomain                  *regexp.Regexp = regexp.MustCompile(`(?:[_a-z0-9](?:[_a-z0-9-]{0,61}[a-z0-9])?\.)+(?:[a-z](?:[a-z0-9-]{0,61}[a-z0-9])?)?`)
	regexNickNameWithAt          *regexp.Regexp = regexp.MustCompile(`^[a-zA-Z0-9@\-\.]*$`)
	regexYear                    *regexp.Regexp = regexp.MustCompile(`(19|20)\d{2}`)
	regexDateRange               *regexp.Regexp = regexp.MustCompile(`(19|20)\d{2}.*(19|20)\d{2}`)
	regexAllSymbolsAfterDot      *regexp.Regexp = regexp.MustCompile(`\..*`)
	regexAllSymbolsAfterQuestion *regexp.Regexp = regexp.MustCompile(`\?.*`)
)

func ContainsItem(arr []string, item string) bool {
	for _, val := range arr {
		if val == item {
			return true
		}
	}
	return false
}

func GenerateCombinations(input []string) []string {
	result := append([]string{}, input...)

	for i := 0; i < len(input)-1; i++ {
		for j := i + 1; j < len(input) && j-i+1 <= 3; j++ {
			substring := strings.Join(input[i:j+1], " ")
			result = append(result, substring)
		}
	}

	return result
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
	//pgConnectionStr := os.Getenv("PG_CONNECTION_STR")
	pgConnectionStr := "user=postgres dbname=parser sslmode=disable password=mysecretpassword host=localhost port=5432"

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

	return s, nil
}

type Experience struct {
	Name     string `json:"name"`
	Start    string `json:"start"`
	End      string `json:"end"`
	Sentence string `json:"sentence"`
}

func Parse(text string) {
	err := godotenv.Load()

	spellInstance, err = LoadSpellFromDB()
	skills := []string{}
	positions := []string{}
	experiences := []Experience{}

	paragraphs := strings.Split(text, "\n")

	for _, paragraph := range paragraphs {
		sentences := strings.Split(paragraph, ".")

		sentencePositions := []string{}
		for _, sentence := range sentences {
			combinations := GenerateCombinations(strings.Split(sentence, " "))
			//foundedPosition := false
			for _, combination := range combinations {
				list, _ := spellInstance.Lookup(combination, spell.SuggestionLevel(spell.LevelClosest))
				for _, l := range list {
					if l.WordData["type"] == "skill" && !ContainsItem(skills, l.Word) {
						skills = append(skills, l.Word)
					}
					if l.WordData["type"] == "position" {
						if !ContainsItem(positions, l.Word) {
							positions = append(positions, l.Word)
						}
						sentencePositions = append(sentencePositions, l.Word)
						//foundedPosition = true
					}
				}
			}

			if len(sentencePositions) > 0 {
				dataRange := regexDateRange.FindAllString(sentence, -1)
				if len(dataRange) > 0 {

					years := regexYear.FindAllString(dataRange[0], -1)
					if len(years) >= 2 {
						//fmt.Println(sentencePositions, years[0], years[1])
						experiences = append(experiences, Experience{
							Name:     sentencePositions[0],
							Start:    years[0],
							End:      years[1],
							Sentence: sentence,
						})
					}
				}

			}
			//fmt.Println(foundedPosition)

		}
	}
	//fmt.Println(skills)
	//fmt.Println(positions)

	for _, experience := range experiences {
		fmt.Println(experience)
	}
	/*	for _, position := range positions {
		fmt.Println(position)
	}*/
	fmt.Println(err)

}
