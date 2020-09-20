package mortems

import (
	"fmt"

	"github.com/google/go-github/v32/github"
)

type MortemCollector struct {
	Git GitService
}

func NewMortemCollector(gitService GitService) MortemCollector {
	return MortemCollector{gitService}
}

func (m MortemCollector) Collect() error {
	files, err := m.Git.GetFiles()
	if err != nil {
		return err
	}

	var newFiles []*github.TreeEntry

	if !containsMortemDirectory(files) {
		newFiles = append(newFiles,
			&github.TreeEntry{
				Path:    github.String("mortems/template.md"),
				Mode:    github.String("100644"),
				Type:    github.String("blob"),
				Content: github.String("<!-- Make sure that"),
			},
		)
	}

	if len(newFiles) > 0 {
		err := m.Git.CommitNewFiles(newFiles)
		if err != nil {
			return fmt.Errorf("could not commit new files: %w", err)
		}
	}

	return nil
}

func containsMortemDirectory(files []*github.TreeEntry) bool {
	for _, file := range files {
		if file.GetPath() == "mortems" {
			return true
		}
	}

	return false
}
