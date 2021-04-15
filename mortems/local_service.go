package mortems

import (
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"strings"
)

type LocalFileService struct {
	repoPath string
}

func NewLocalFileService(repoPath string) RepoFileService {
	return &LocalFileService{
		repoPath,
	}
}

func (s *LocalFileService) GetFiles() (*RepoFiles, error) {
	var files RepoFiles

	err := filepath.Walk(s.repoPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		if strings.Contains(path, ".git") {
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

		absoluteRepo, err := filepath.Abs(s.repoPath)
		if err != nil {
			return err
		}

		relativePath, err := filepath.Rel(absoluteRepo, path)
		if err != nil {
			return err
		}

		files.AddFile(relativePath, string(content))

		return nil
	})

	if err != nil {
		return nil, err
	}

	// for _, f := range files.Files {
	// 	fmt.Printf("file: %v\n", f)
	// }

	return &files, nil
}

func (s *LocalFileService) CommitNewFiles(files *RepoFiles) error {
	for _, file := range files.Files {
		err := os.WriteFile(path.Join(s.repoPath, file.GetPath()), []byte(file.GetContent()), 0644)
		if err != nil {
			return err
		}
	}

	return nil
}
