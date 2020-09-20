package main

import (
	"context"
	"fmt"
	"os"

	"github.com/google/go-github/v32/github"
	"golang.org/x/oauth2"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "Oh jeez: %s\n", err)
	}
}

func run() error {
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: "xxx"},
	)
	tc := oauth2.NewClient(ctx, ts)

	client := github.NewClient(tc)

	// list all repositories for the authenticated user
	repo, _, err := client.Repositories.Get(ctx, "ostenbom", "test-mortems")
	if err != nil {
		return fmt.Errorf("could not get repos: %w", err)
	}

	fmt.Println(repo)

	ref, _, err := client.Git.GetRef(ctx, "ostenbom", "test-mortems", "heads/master")
	if err != nil {
		return fmt.Errorf("could not get ref: %w", err)
	}

	fmt.Println(ref)

	commit, _, err := client.Git.GetCommit(ctx, "ostenbom", "test-mortems", ref.Object.GetSHA())
	if err != nil {
		return fmt.Errorf("could not get commit: %w", err)
	}

	tree, _, err := client.Git.GetTree(ctx, "ostenbom", "test-mortems", commit.GetTree().GetSHA(), true)
	if err != nil {
		return fmt.Errorf("could not get tree: %w", err)
	}

	for _, entry := range tree.Entries {
		fmt.Println(entry.GetPath())
		fmt.Println(entry.GetType())

		if entry.GetType() == "blob" {
			blob, _, err := client.Git.GetBlobRaw(ctx, "ostenbom", "test-mortems", entry.GetSHA())
			if err != nil {
				return fmt.Errorf("could not get tree: %w", err)
			}

			fmt.Println(string(blob))
		}
	}

	return nil
}
