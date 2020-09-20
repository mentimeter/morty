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

	fmt.Printf("Token: %s\n", token[:5])
	fmt.Printf("Env token: %s\n", os.Getenv("GITHUB_REPOSITORY"))

	return nil
}
