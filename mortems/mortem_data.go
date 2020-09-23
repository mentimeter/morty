package mortems

import (
	"errors"
	"regexp"
	"time"
)

type MortemData struct {
	File      string        `json:"file"`
	Title     string        `json:"title"`
	Owner     string        `json:"owner"`
	Date      time.Time     `json:"date"`
	Severity  string        `json:"severity"`
	Detect    time.Duration `json:"detect"`
	Resolve   time.Duration `json:"resolve"`
	TotalDown time.Duration `json:"total_down"`
}

var ErrNoTitle = errors.New("no title of format \"# Title Here\"")

func NewMortemData(fileContent string) {

}

func ParseTitle(content string) (string, error) {
	re := regexp.MustCompile(`#\s(?P<Title>.+)`)

	title := re.FindStringSubmatch(content)
	if title == nil {
		return "", ErrNoTitle
	}

	return title[1], nil
}
