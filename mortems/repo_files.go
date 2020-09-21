package mortems

import "github.com/google/go-github/v32/github"

type RepoFiles struct {
	Files []*github.TreeEntry
}

func (r *RepoFiles) Size() int {
	return len(r.Files)
}

func (r *RepoFiles) GetFile(path string) *github.TreeEntry {
	for _, file := range r.Files {
		if file.GetPath() == path {
			return file
		}
	}

	return nil
}

func (r *RepoFiles) AddFile(path, content string) {
	file := &github.TreeEntry{
		Path:    github.String(path),
		Mode:    github.String("100644"),
		Type:    github.String("blob"),
		Content: github.String(content),
	}
	r.Files = append(r.Files, file)
}
