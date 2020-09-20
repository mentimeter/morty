package mortems_test

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/google/go-github/v32/github"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/ostenbom/morty/mortems"
	"github.com/ostenbom/morty/mortems/mortemsfakes"
)

var _ = Describe("Mortems", func() {
	var mortems MortemCollector
	var gitService *mortemsfakes.FakeGitService
	var treeEntryFixtures map[string][]*github.TreeEntry

	BeforeEach(func() {
		var err error
		treeEntryFixtures, err = loadTreeEntryFixtures()
		Expect(err).NotTo(HaveOccurred())

		gitService = new(mortemsfakes.FakeGitService)
		mortems = NewMortemCollector(gitService)
	})

	Context("when the repository contains no post-mortems", func() {
		BeforeEach(func() {
			gitService.GetFilesReturns(treeEntryFixtures["no-mortems-dir"], nil)
		})

		It("Gets the TreeEntries from Git", func() {
			Expect(mortems.Collect()).To(Succeed())
			Expect(gitService.GetFilesCallCount()).To(Equal(1))
		})
	})
})

func loadTreeEntryFixtures() (map[string][]*github.TreeEntry, error) {
	fixturesDir := "testdata"

	fixtureDirectories, err := ioutil.ReadDir(fixturesDir)
	if err != nil {
		return nil, err
	}

	fixtures := make(map[string][]*github.TreeEntry)

	for _, dir := range fixtureDirectories {
		var files []*github.TreeEntry

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

			file := &github.TreeEntry{
				Path:    &path,
				Content: github.String(string(content)),
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

		fixtures[dir.Name()] = files
	}

	return fixtures, nil
}
