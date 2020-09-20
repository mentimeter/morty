package main

import (
	"errors"
	"fmt"
	"os"

	"github.com/ostenbom/morty/mortems"
	"github.com/sethvargo/go-githubactions"
)

var ErrMissingInput = errors.New("missing required input")

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "Oh jeez: %s\n", err)
		os.Exit(1)
	}
}

func run() error {
	token := githubactions.GetInput("token")
	if token == "" {
		return fmt.Errorf("no github token: %w", ErrMissingInput)
	}

	repository := os.Getenv("GITHUB_REPOSITORY")
	if repository == "" {
		return fmt.Errorf("no repository: %w", ErrMissingInput)
	}

	ref := os.Getenv("GITHUB_REF")
	if ref == "" {
		return fmt.Errorf("no ref: %w", ErrMissingInput)
	}

	gitService := mortems.NewGitHubService(token, repository, ref)

	mortemCollector := mortems.NewMortemCollector(gitService)

	return mortemCollector.Collect()
}
