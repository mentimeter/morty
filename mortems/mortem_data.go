package mortems

import (
	"errors"
	"fmt"
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
var ErrNoOwner = errors.New("no owner of format \"Owner: First Last\"")
var ErrNoDate = errors.New("no date of format \"July 1, 2020\" (no st, nd, th)")
var ErrNoSeverity = errors.New("no severity of format \"| Severity | sev |\"")

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

func ParseOwner(content string) (string, error) {
	re := regexp.MustCompile(`.*Owner: (?P<Owner>.+)`)

	owner := re.FindStringSubmatch(content)
	if owner == nil {
		return "", ErrNoOwner
	}

	return owner[1], nil
}

func ParseDate(content string) (time.Time, error) {
	re := regexp.MustCompile(`.*Date: (?P<Date>.+)`)

	dateStr := re.FindStringSubmatch(content)
	if dateStr == nil {
		return time.Now(), ErrNoDate
	}

	date, err := time.Parse("January 2, 2006", dateStr[1])
	if err != nil {
		return time.Now(), fmt.Errorf("incorrect date format, %s: %w", err, ErrNoDate)
	}

	return date, nil
}

func ParseSeverity(content string) (string, error) {
	re := regexp.MustCompile(` *\| *Severity +\| *(.+) +\| *`)

	sev := re.FindStringSubmatch(content)
	if sev == nil {
		return "", ErrNoOwner
	}

	return sev[1], nil
}
