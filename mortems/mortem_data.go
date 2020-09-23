package mortems

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"
)

type MortemData struct {
	File     string        `json:"file"`
	Title    string        `json:"title"`
	Owner    string        `json:"owner"`
	Date     time.Time     `json:"date"`
	Severity string        `json:"severity"`
	Detect   time.Duration `json:"detect"`
	Resolve  time.Duration `json:"resolve"`
	Downtime time.Duration `json:"total_down"`
}

var ErrNoTitle = errors.New("no title of format \"# Title Here\"")
var ErrNoOwner = errors.New("no owner of format \"Owner: First Last\"")
var ErrNoDate = errors.New("no date of format \"July 1, 2020\" (no st, nd, th)")
var ErrNoSeverity = errors.New("no severity of format \"| Severity | sev |\"")
var ErrNoDetect = errors.New("no time to detect of format \"| Time to Detect | x unit[, y smaller_unit] |\"")
var ErrNoResolve = errors.New("no time to resolve of format \"| Time to Resolve | x unit[, y smaller_unit] |\"")
var ErrNoDowntime = errors.New("no total downtime format \"| Total Downtime | x unit[, y smaller_unit] |\"")

func NewMortemData(content, path string) (MortemData, error) {
	title, err := ParseTitle(content)
	if err != nil {
		return MortemData{}, err
	}

	owner, err := ParseOwner(content)
	if err != nil {
		return MortemData{}, err
	}

	date, err := ParseDate(content)
	if err != nil {
		return MortemData{}, err
	}

	severity, err := ParseSeverity(content)
	if err != nil {
		return MortemData{}, err
	}

	detect, err := ParseDetect(content)
	if err != nil {
		return MortemData{}, err
	}

	resolve, err := ParseResolve(content)
	if err != nil {
		return MortemData{}, err
	}

	downtime, err := ParseDowntime(content)
	if err != nil {
		return MortemData{}, err
	}

	return MortemData{
		File:     path,
		Title:    title,
		Owner:    owner,
		Date:     date,
		Severity: severity,
		Detect:   detect,
		Resolve:  resolve,
		Downtime: downtime,
	}, nil
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

	var date time.Time
	var err error

	date, err = time.Parse("January 2, 2006", dateStr[1])
	if err != nil {
		date, err = time.Parse("Jan 2, 2006", dateStr[1])
		if err != nil {
			return time.Now(), fmt.Errorf("incorrect date format, %s: %w", err, ErrNoDate)
		}
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

func ParseDetect(content string) (time.Duration, error) {
	re := regexp.MustCompile(` *\|.*Detect +\| *(.+) +\| *`)

	detectMatches := re.FindStringSubmatch(content)
	if detectMatches == nil {
		return 0, ErrNoOwner
	}

	detectString := detectMatches[1]

	detectTime, err := stringToTime(detectString)
	if err != nil {
		return 0, fmt.Errorf("%w: %s", ErrNoDetect, err)
	}

	return detectTime, nil
}

func ParseResolve(content string) (time.Duration, error) {
	re := regexp.MustCompile(` *\|.*Resolve +\| *(.+) +\| *`)

	resolveMatches := re.FindStringSubmatch(content)
	if resolveMatches == nil {
		return 0, ErrNoOwner
	}

	resolveString := resolveMatches[1]

	resolveTime, err := stringToTime(resolveString)
	if err != nil {
		return 0, fmt.Errorf("%w: %s", ErrNoResolve, err)
	}

	return resolveTime, nil
}

func ParseDowntime(content string) (time.Duration, error) {
	re := regexp.MustCompile(` *\|.*Downtime +\| *(.+) +\| *`)

	downtimeMatches := re.FindStringSubmatch(content)
	if downtimeMatches == nil {
		return 0, ErrNoOwner
	}

	downtimeString := downtimeMatches[1]

	downtimeTime, err := stringToTime(downtimeString)
	if err != nil {
		return 0, fmt.Errorf("%w: %s", ErrNoDowntime, err)
	}

	return downtimeTime, nil
}

func stringToTime(timeStr string) (time.Duration, error) {
	noSpaceDetect := strings.ReplaceAll(timeStr, " ", "")
	timeGroups := strings.Split(noSpaceDetect, ",")

	totalTime := time.Duration(0)
	for _, t := range timeGroups {
		re := regexp.MustCompile("[^0-9]*([0-9]+)[^0-9]*")
		goTimeString := ""

		if strings.Contains(t, "day") {
			timeString := re.FindStringSubmatch(t)
			if timeString == nil {
				return 0, errors.New("missing number of days")
			}

			days, err := strconv.Atoi(timeString[1])
			if err != nil {
				return 0, fmt.Errorf("could not parse days: %w", err)
			}

			hours := days * 24
			goTimeString = strconv.Itoa(hours) + "h"
		} else if strings.Contains(t, "hour") {
			timeString := re.FindStringSubmatch(t)
			if timeString == nil {
				return 0, errors.New("missing number of hours")
			}

			goTimeString = timeString[1] + "h"
		} else if strings.Contains(t, "min") {
			timeString := re.FindStringSubmatch(t)
			if timeString == nil {
				return 0, errors.New("missing number of minutes")
			}

			goTimeString = timeString[1] + "m"
		} else if strings.Contains(t, "sec") {
			timeString := re.FindStringSubmatch(t)
			if timeString == nil {
				return 0, errors.New("missing number of seconds")
			}

			goTimeString = timeString[1] + "s"
		}

		thisTime, err := time.ParseDuration(goTimeString)
		if err != nil {
			return 0, ErrNoDetect
		}

		totalTime += thisTime
	}

	return totalTime, nil
}
