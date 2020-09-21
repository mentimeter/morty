package mortems

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/google/go-github/v32/github"
)

type MortemCollector struct {
	Git GitService
}

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

func NewMortemCollector(gitService GitService) MortemCollector {
	return MortemCollector{gitService}
}

func (m MortemCollector) Collect() error {
	existingFiles, err := m.Git.GetFiles()
	if err != nil {
		return err
	}

	newFiles := RepoFiles{}

	templatePath := "post-mortems/template.md"
	templateFile := existingFiles.GetFile(templateContent)

	if templateFile == nil || templateFile.GetContent() != templateContent {
		newFiles.AddFile(templatePath, templateContent)
	}

	howToPath := "post-mortems/README.md"
	howToFile := existingFiles.GetFile(howToPath)

	if howToFile == nil || howToFile.GetContent() != howToContent {
		newFiles.AddFile(howToPath, howToContent)
	}

	databasePath := "mortems.json"
	modifiedDatabase := false

	databaseFile := existingFiles.GetFile(databasePath)
	if databaseFile == nil {
		modifiedDatabase = true
	}

	var database []MortemData

	// databaseBytes := []byte(databaseFile.GetContent())
	//
	// err = json.Unmarshal(databaseBytes, &database)
	// if err != nil {
	// 	return fmt.Errorf("could not load database from file: %w", err)
	// }

	for _, file := range existingFiles.Files {
		if strings.HasPrefix(file.GetPath(), "post-mortems/") &&
			file.GetPath() != howToPath &&
			file.GetPath() != templatePath {
			modifiedDatabase = true

			mortem, err := ParseMortem(file.GetContent())
			if err != nil {
				return fmt.Errorf("could not parse data from mortem %s: %w", file.GetPath(), err)
			}

			database = append(database, mortem)
		}
	}

	databaseBytes, err := json.Marshal(database)
	if err != nil {
		return fmt.Errorf("could not marshal database to json: %w", err)
	}

	databaseString := string(databaseBytes)

	if modifiedDatabase {
		newFiles.AddFile(databasePath, databaseString)
	}

	if newFiles.Size() > 0 {
		err := m.Git.CommitNewFiles(&newFiles)
		if err != nil {
			return fmt.Errorf("could not commit new files: %w", err)
		}
	}

	return nil
}

func ParseMortem(mortem string) (MortemData, error) {
	detectTime, err := time.ParseDuration("4m")
	if err != nil {
		return MortemData{}, err
	}

	resolveTime, err := time.ParseDuration("6h14m")
	if err != nil {
		return MortemData{}, err
	}

	totalDownTime, err := time.ParseDuration("6h28m")
	if err != nil {
		return MortemData{}, err
	}

	return MortemData{
		File:      "0001-first-mortem.md",
		Title:     "Love Lost Globally: Jerry Develops Malicious App",
		Owner:     "Morty Smith",
		Date:      time.Date(2020, time.July, 1, 0, 0, 0, 0, time.UTC),
		Severity:  "1",
		Detect:    detectTime,
		Resolve:   resolveTime,
		TotalDown: totalDownTime,
	}, nil
}

func newTreeEntryFile(path string, content string) *github.TreeEntry {
	return &github.TreeEntry{
		Path:    github.String(path),
		Mode:    github.String("100644"),
		Type:    github.String("blob"),
		Content: github.String(content),
	}
}
