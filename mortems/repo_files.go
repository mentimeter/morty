package mortems

import "github.com/google/go-github/v32/github"

type RepoFiles struct {
	Files []*File
}

type File struct {
	Path    string
	Mode    string
	Type    string
	Content string
}

func (r *RepoFiles) Size() int {
	return len(r.Files)
}

func (r *RepoFiles) GetFile(path string) *File {
	for _, file := range r.Files {
		if file.GetPath() == path {
			return file
		}
	}

	return nil
}

func (r *RepoFiles) AddFile(path, content string) {
	file := &File{
		Path:    path,
		Mode:    "100644",
		Type:    "blob",
		Content: content,
	}
	r.Files = append(r.Files, file)
}

func (r *RepoFiles) ToTreeEntries() []*github.TreeEntry {
	var entries []*github.TreeEntry

	for _, file := range r.Files {
		entry := &github.TreeEntry{
			Path:    github.String(file.Path),
			Mode:    github.String(file.Mode),
			Type:    github.String(file.Type),
			Content: github.String(file.Content),
		}
		entries = append(entries, entry)
	}

	return entries
}

func (f *File) GetContent() string {
	return f.Content
}

func (f *File) GetPath() string {
	return f.Path
}

func (f *File) GetMode() string {
	return f.Mode
}

func (f *File) GetType() string {
	return f.Type
}
