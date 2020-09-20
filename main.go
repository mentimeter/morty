package main

import (
	"errors"
	"fmt"
	"os"

	"github.com/sethvargo/go-githubactions"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "Oh jeez: %s\n", err)
		os.Exit(1)
	}
}

func run() error {
	token := githubactions.GetInput("token")
	if token == "" {
		return errors.New("missing 'token'")
	}

	repo := githubactions.GetInput("repository")
	if repo == "" {
		return errors.New("missing 'repository'")
	}

	fmt.Printf("Repo: %s\n", repo)
	fmt.Printf("Token: %s\n", token[:4])

	return nil
}
