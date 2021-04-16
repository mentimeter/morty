package mortems

import (
	_ "embed"
	"encoding/json"
	"fmt"
	"strings"
)

//go:embed install_morty
var installScriptContent string

type MortemCollector struct {
	Repo            RepoFileService
	Reporter        ReportingService
	previousMortems []MortemData
	nextMortems     []MortemData
}

func NewMortemCollector(fileService RepoFileService) MortemCollector {
	return MortemCollector{fileService, nil, []MortemData{}, []MortemData{}}
}

func NewMortemReportingCollector(fileService RepoFileService, reportingService ReportingService) MortemCollector {
	return MortemCollector{fileService, reportingService, []MortemData{}, []MortemData{}}
}

func (m *MortemCollector) Check() (RepoFiles, error) {
	newFiles := RepoFiles{}

	existingFiles, err := m.Repo.GetFiles()
	if err != nil {
		return RepoFiles{}, err
	}

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

	installScriptPath := "install_morty"
	installScriptFile := existingFiles.GetFile(installScriptPath)

	if installScriptFile == nil || installScriptFile.GetContent() != installScriptContent {
		newFiles.AddExecutableFile(installScriptPath, installScriptContent)
	}

	databasePath := "mortems.json"
	modifiedDatabase := false

	databaseFile := existingFiles.GetFile(databasePath)
	if databaseFile == nil {
		modifiedDatabase = true
		databaseFile = &File{
			Path:    databasePath,
			Mode:    "100644",
			Type:    "blob",
			Content: "[]",
		}
	}

	databaseBytes := []byte(databaseFile.GetContent())

	err = json.Unmarshal(databaseBytes, &m.previousMortems)
	if err != nil {
		return RepoFiles{}, fmt.Errorf("could not load database from file: %w", err)
	}

	for _, file := range existingFiles.Files {
		if strings.HasPrefix(file.GetPath(), "post-mortems/") &&
			!strings.HasPrefix(file.GetPath(), "post-mortems/images/") &&
			file.GetPath() != howToPath &&
			file.GetPath() != templatePath {
			modifiedDatabase = true

			mortem, err := NewMortemData(file.GetContent(), file.GetPath())
			if err != nil {
				return RepoFiles{}, fmt.Errorf("could not parse data from mortem %s: %w", file.GetPath(), err)
			}

			m.nextMortems = append(m.nextMortems, mortem)
		}
	}

	readmePath := "README.md"
	readmeFile := existingFiles.GetFile(readmePath)
	readmeContent := GenerateReadme(m.nextMortems)

	if readmeFile == nil || readmeFile.GetContent() != readmeContent {
		newFiles.AddFile(readmePath, readmeContent)
	}

	databaseBytes, err = json.Marshal(m.nextMortems)
	if err != nil {
		return RepoFiles{}, fmt.Errorf("could not marshal database to json: %w", err)
	}

	databaseString := string(databaseBytes)

	if modifiedDatabase {
		newFiles.AddFile(databasePath, databaseString)
	}

	return newFiles, nil

}

func (m *MortemCollector) Collect() error {
	newFiles, err := m.Check()
	if err != nil {
		return fmt.Errorf("could not check files, there might be a parsing error: %w", err)
	}

	if m.Reporter != nil {
		toReport := diffMortems(m.previousMortems, m.nextMortems)
		for _, report := range toReport {

			m.Reporter.ReportSeverity(report.Severity)
			m.Reporter.ReportDetect(report.Detect, report.Severity)
			m.Reporter.ReportResolve(report.Resolve, report.Severity)
			m.Reporter.ReportDowntime(report.Downtime, report.Severity)
		}
	}

	if newFiles.Size() > 0 {
		err := m.Repo.CommitNewFiles(&newFiles)
		if err != nil {
			return fmt.Errorf("could not commit new files: %w", err)
		}
	}

	return nil
}

func diffMortems(old []MortemData, new []MortemData) []MortemData {
	diff := []MortemData{}
	for _, newMortem := range new {
		exists := false
		for _, oldMortem := range old {
			if oldMortem == newMortem {
				exists = true
			}
		}
		if !exists {
			diff = append(diff, newMortem)
		}
	}

	return diff
}
