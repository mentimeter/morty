package mortems

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/google/go-github/v35/github"
	"golang.org/x/oauth2"
)

//go:generate go run github.com/maxbrunsfeld/counterfeiter/v6 . RepoFileService

type RepoFileService interface {
	GetFiles() (*RepoFiles, error)
	CommitNewFiles(*RepoFiles) error
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

func NewGitHubService(token, fullRepository, ref string) RepoFileService {
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

func (g *GitHub) GetFiles() (*RepoFiles, error) {
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

	var files []*File

	for _, entry := range tree.Entries {
		if entry.GetType() == "blob" {
			content, _, err := g.client.Git.GetBlobRaw(ctx, g.owner, g.repository, entry.GetSHA())
			if err != nil {
				return nil, fmt.Errorf("could not get blob content: %w", err)
			}

			entry.Content = github.String(string(content))

			file := &File{
				Path:    *entry.Path,
				Mode:    *entry.Mode,
				Type:    "blob",
				Content: string(content),
			}

			files = append(files, file)
		}
	}

	return &RepoFiles{files}, nil
}

func (g *GitHub) CommitNewFiles(updateEntries *RepoFiles) error {
	ctx := context.Background()

	newTree, _, err := g.client.Git.CreateTree(ctx, g.owner, g.repository, g.currentTree.GetSHA(), updateEntries.ToTreeEntries())
	if err != nil {
		return fmt.Errorf("could not create a new tree: %w", err)
	}

	authorDate := time.Now()
	author := &github.CommitAuthor{
		Date:  &authorDate,
		Name:  g.currentCommit.Author.Name,
		Email: g.currentCommit.Author.Email,
	}

	newCommitData := &github.Commit{
		Author:    author,
		Committer: author,
		Message:   github.String("Auto-morty-cally analyse post-mortems"),
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
