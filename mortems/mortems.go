package mortems

type MortemCollector struct {
	Git GitService
}

func NewMortemCollector(gitService GitService) MortemCollector {
	return MortemCollector{gitService}
}

func (m MortemCollector) Collect() error {
	_, err := m.Git.GetTreeEntries()
	if err != nil {
		return err
	}
	return nil
}
