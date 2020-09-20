package mortems

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/google/go-github/v32/github"
	"golang.org/x/oauth2"
)

//go:generate go run github.com/maxbrunsfeld/counterfeiter/v6 . GitService

type GitService interface {
	GetFiles() ([]*github.TreeEntry, error)
	CommitNewFiles([]*github.TreeEntry) error
}

type GitHub struct {
	client        *github.Client
	ref           string
	owner         string
	repository    string
	currentTree   *github.Tree
	currentCommit *github.Commit
	currentRef    *github.Reference
}

func NewGitHubService(token, fullRepository, ref string) GitService {
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	tc := oauth2.NewClient(ctx, ts)
	client := github.NewClient(tc)
	repositoryArgs := strings.Split(fullRepository, "/")

	return &GitHub{
		client:     client,
		ref:        ref,
		owner:      repositoryArgs[0],
		repository: repositoryArgs[1],
	}
}

func (g *GitHub) GetFiles() ([]*github.TreeEntry, error) {
	ctx := context.Background()

	ref, _, err := g.client.Git.GetRef(ctx, g.owner, g.repository, g.ref)
	if err != nil {
		return nil, fmt.Errorf("could not get ref: %w", err)
	}

	g.currentRef = ref

	commit, _, err := g.client.Git.GetCommit(ctx, g.owner, g.repository, ref.Object.GetSHA())
	if err != nil {
		return nil, fmt.Errorf("could not get commit: %w", err)
	}

	g.currentCommit = commit

	tree, _, err := g.client.Git.GetTree(ctx, g.owner, g.repository, commit.GetTree().GetSHA(), true)
	if err != nil {
		return nil, fmt.Errorf("could not get tree: %w", err)
	}

	g.currentTree = tree

	var files []*github.TreeEntry

	for _, entry := range tree.Entries {
		if entry.GetType() == "blob" {
			content, _, err := g.client.Git.GetBlobRaw(ctx, g.owner, g.repository, entry.GetSHA())
			if err != nil {
				return nil, fmt.Errorf("could not get blob content: %w", err)
			}

			entry.Content = github.String(string(content))
			files = append(files, entry)
		}
	}

	return files, nil
}

func (g *GitHub) CommitNewFiles(updateEntries []*github.TreeEntry) error {
	ctx := context.Background()

	newTree, _, err := g.client.Git.CreateTree(ctx, g.owner, g.repository, g.currentTree.GetSHA(), updateEntries)
	if err != nil {
		return fmt.Errorf("could not create a new tree: %w", err)
	}

	authorDate := time.Now()
	author := &github.CommitAuthor{
		Date:  &authorDate,
		Name:  github.String("Morty Smith"),
		Email: github.String("morty@your-post-mortems.now"),
	}

	newCommitData := &github.Commit{
		Author:    author,
		Committer: author,
		Message:   github.String("morty: collect your mortems"),
		Tree:      newTree,
		Parents:   []*github.Commit{g.currentCommit},
	}

	createdCommit, _, err := g.client.Git.CreateCommit(ctx, g.owner, g.repository, newCommitData)
	if err != nil {
		return fmt.Errorf("could not create a new commit: %w", err)
	}

	g.currentRef.Object.SHA = createdCommit.SHA

	_, _, err = g.client.Git.UpdateRef(ctx, g.owner, g.repository, g.currentRef, false)
	if err != nil {
		return fmt.Errorf("could not update the current ref: %w", err)
	}

	return nil
}
