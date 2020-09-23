package mortems

import (
	"encoding/json"
	"fmt"
	"strings"
)

type MortemCollector struct {
	Git GitService
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

	var mortems []MortemData

	// databaseBytes := []byte(databaseFile.GetContent())
	//
	// err = json.Unmarshal(databaseBytes, &database)
	// if err != nil {
	// 	return fmt.Errorf("could not load database from file: %w", err)
	// }

	for _, file := range existingFiles.Files {
		if strings.HasPrefix(file.GetPath(), "post-mortems/") &&
			!strings.HasPrefix(file.GetPath(), "post-mortems/images/") &&
			file.GetPath() != howToPath &&
			file.GetPath() != templatePath {
			modifiedDatabase = true

			mortem, err := NewMortemData(file.GetContent(), file.GetPath())
			if err != nil {
				return fmt.Errorf("could not parse data from mortem %s: %w", file.GetPath(), err)
			}

			mortems = append(mortems, mortem)
		}
	}

	readmePath := "README.md"
	readmeFile := existingFiles.GetFile(readmePath)
	readmeContent := GenerateReadme(mortems)

	if readmeFile == nil || readmeFile.GetContent() != readmeContent {
		newFiles.AddFile(readmePath, readmeContent)
	}

	databaseBytes, err := json.Marshal(mortems)
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
