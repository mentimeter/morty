package mortems

import (
	"fmt"
	"sort"
	"time"
)

func GenerateReadme(mortems []MortemData) string {
	months := make(map[string][]MortemData)

	monthYearFormat := "January 2006"

	for _, mortem := range mortems {
		month := mortem.Date.Format(monthYearFormat)
		months[month] = append(months[month], mortem)
	}

	readme := "# Post-Mortems\n"

	readme += "## Overall Statistics\n"
	readme += metricsTable(metrics(mortems))

	sortedMonths := []string{}
	for month := range months {
		sortedMonths = append(sortedMonths, month)
	}

	sort.Slice(sortedMonths, func(i, j int) bool {
		timeI, _ := time.Parse(monthYearFormat, sortedMonths[i])
		timeJ, _ := time.Parse(monthYearFormat, sortedMonths[j])

		return timeI.After(timeJ)
	})

	for _, month := range sortedMonths {
		readme += monthSection(month, months[month])
	}

	return readme
}

func monthSection(month string, mortems []MortemData) string {
	section := "### " + month + "\n"
	section += metricsTable(metrics(mortems))

	for _, m := range mortems {
		section += fmt.Sprintf("- [%s](%s)\n", m.Title, m.File)
	}

	section += "\n"

	return section
}

func metricsTable(detect, resolve, down string) string {
	return fmt.Sprintf(`| Average Detection Time | Average Resolve Time | Average Downtime |
| --- | --- | --- |
| %s | %s | %s |
`, detect, resolve, down)
}

func metrics(mortems []MortemData) (string, string, string) {
	if len(mortems) == 0 {
		return "", "", ""
	}

	detectTotal := 0
	resolveTotal := 0
	downTotal := 0

	for _, m := range mortems {
		detectTotal += int(m.Detect)
		resolveTotal += int(m.Resolve)
		downTotal += int(m.Downtime)
	}

	detectAvg := detectTotal / len(mortems)
	resolveAvg := resolveTotal / len(mortems)
	downAvg := downTotal / len(mortems)

	detect := prettyTime(detectAvg)
	resolve := prettyTime(resolveAvg)
	down := prettyTime(downAvg)

	return detect, resolve, down
}

func prettyTime(t int) string {
	dur := time.Duration(t)

	if dur > time.Hour {
		return fmt.Sprintf("%.0f hours", dur.Hours())
	} else if dur > time.Minute {
		return fmt.Sprintf("%.0f minutes", dur.Minutes())
	} else {
		return fmt.Sprintf("%.0f seconds", dur.Seconds())
	}
}
