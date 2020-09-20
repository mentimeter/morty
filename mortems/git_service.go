package mortems

import (
	"context"
	"fmt"
	"strings"

	"github.com/google/go-github/v32/github"
	"golang.org/x/oauth2"
)

//go:generate go run github.com/maxbrunsfeld/counterfeiter/v6 . GitService

type GitService interface {
	GetFiles() ([]*File, error)
	CommitNewBlobs([]*github.Blob) error
}

type GitHub struct {
	client     *github.Client
	ref        string
	owner      string
	repository string
}

type File struct {
	TreeEntry *github.TreeEntry
	Content   []byte
}

func NewGitHubService(token, fullRepository, ref string) GitService {
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	tc := oauth2.NewClient(ctx, ts)
	client := github.NewClient(tc)
	repositoryArgs := strings.Split(fullRepository, "/")

	return &GitHub{client, ref, repositoryArgs[0], repositoryArgs[1]}
}

func (g *GitHub) GetFiles() ([]*File, error) {
	ctx := context.Background()

	ref, _, err := g.client.Git.GetRef(ctx, g.owner, g.repository, g.ref)
	if err != nil {
		return nil, fmt.Errorf("could not get ref: %w", err)
	}

	commit, _, err := g.client.Git.GetCommit(ctx, g.owner, g.repository, ref.Object.GetSHA())
	if err != nil {
		return nil, fmt.Errorf("could not get commit: %w", err)
	}

	tree, _, err := g.client.Git.GetTree(ctx, g.owner, g.repository, commit.GetTree().GetSHA(), true)
	if err != nil {
		return nil, fmt.Errorf("could not get tree: %w", err)
	}

	var files []*File

	for _, entry := range tree.Entries {
		if entry.GetType() == "blob" {
			content, _, err := g.client.Git.GetBlobRaw(ctx, g.owner, g.repository, entry.GetSHA())
			if err != nil {
				return nil, fmt.Errorf("could not get blob content: %w", err)
			}

			file := &File{entry, content}
			files = append(files, file)
		}
	}

	return files, nil
}

func (g *GitHub) CommitNewBlobs([]*github.Blob) error {
	return nil
}
