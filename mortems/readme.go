package mortems

import (
	"fmt"
	"time"
)

func GenerateReadme(mortems []MortemData) string {
	months := make(map[string][]MortemData)

	for _, mortem := range mortems {
		month := mortem.Date.Format("January 2006")
		months[month] = append(months[month], mortem)
	}

	readme := "# Post-Mortems\n"

	readme += "## Overall Statistics\n"
	readme += metricsTable(metrics(mortems))

	for month, monthMortems := range months {
		readme += monthSection(month, monthMortems)
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
