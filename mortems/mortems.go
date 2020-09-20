package mortems

import "fmt"

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

	fmt.Println(files)

	return nil
}
