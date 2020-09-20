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
	GetTreeEntries() ([]*github.TreeEntry, error)
	CommitNewBlobs([]*github.Blob) error
}

type GitHub struct {
	client     *github.Client
	ref        string
	owner      string
	repository string
}

func NewGitHubService(ref, token, fullRepository string) GitService {
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	tc := oauth2.NewClient(ctx, ts)
	client := github.NewClient(tc)
	repositoryArgs := strings.Split(fullRepository, "/")

	return &GitHub{client, ref, repositoryArgs[0], repositoryArgs[1]}
}

func (g *GitHub) GetTreeEntries() ([]*github.TreeEntry, error) {
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

	return tree.Entries, nil
}

func (g *GitHub) CommitNewBlobs([]*github.Blob) error {
	return nil
}
