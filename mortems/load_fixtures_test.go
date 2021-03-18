package mortems_test

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	. "github.com/mentimeter/morty/mortems"
)

func LoadTreeEntryFixtures() (map[string]*RepoFiles, error) {
	fixturesDir := "testdata"

	fixtureDirectories, err := ioutil.ReadDir(fixturesDir)
	if err != nil {
		return nil, err
	}

	fixtures := make(map[string]*RepoFiles)

	for _, dir := range fixtureDirectories {
		var files []*File

		err := os.Chdir(filepath.Join(fixturesDir, dir.Name()))
		if err != nil {
			return nil, fmt.Errorf("bad state! could not change fixture dir: %w", err)
		}

		err = filepath.Walk(".", func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if info.IsDir() {
				return nil
			}

			f, err := os.Open(path)
			if err != nil {
				return err
			}
			defer f.Close()

			content, err := ioutil.ReadAll(f)
			if err != nil {
				return err
			}

			file := &File{
				Path:    path,
				Content: string(content),
			}

			// fmt.Printf("file: %s, content: %s\n", path, string(content))

			files = append(files, file)

			return nil
		})
		if err != nil {
			return nil, err
		}

		err = os.Chdir("../../")
		if err != nil {
			return nil, fmt.Errorf("bad state! could not change fixture dir: %w", err)
		}

		fixtures[dir.Name()] = &RepoFiles{files}
	}

	return fixtures, nil
}
