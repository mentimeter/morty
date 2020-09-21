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

	templatePath := "post-mortems/template.md"
	templateFile := getFile(files, templatePath)

	if templateFile == nil || templateFile.GetContent() != postMortemTemplate {
		newFiles = append(newFiles, newTreeEntryFile(templatePath, postMortemTemplate))
	}

	howToPath := "post-mortems/README.md"
	howToFile := getFile(files, howToPath)

	if howToFile == nil || howToFile.GetContent() != howToPostMortem {
		newFiles = append(newFiles, newTreeEntryFile(howToPath, howToPostMortem))
	}

	if len(newFiles) > 0 {
		err := m.Git.CommitNewFiles(newFiles)
		if err != nil {
			return fmt.Errorf("could not commit new files: %w", err)
		}
	}

	return nil
}

func newTreeEntryFile(path string, content string) *github.TreeEntry {
	return &github.TreeEntry{
		Path:    github.String(path),
		Mode:    github.String("100644"),
		Type:    github.String("blob"),
		Content: github.String(content),
	}
}

func getFile(files []*github.TreeEntry, path string) *github.TreeEntry {
	for _, file := range files {
		if file.GetPath() == path {
			return file
		}
	}

	return nil
}
