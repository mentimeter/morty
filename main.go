package main

import (
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/mentimeter/morty/mortems"
	"github.com/sethvargo/go-githubactions"
)

var ErrMissingInput = errors.New("missing required input")

func main() {
	if len(os.Args) <= 1 {
		exitHelp()
	}
	if os.Args[1] == "check" {
		if err := runCheck(); err != nil {
			log.Fatalf("Oh jeez: %s\n", err)
		}
	} else if len(os.Args) == 3 && os.Args[1] == "git" && os.Args[2] == "check" {
		if err := runGitCheck(); err != nil {
			log.Fatalf("Oh jeez: %s\n", err)
		}
	} else if len(os.Args) == 2 && os.Args[1] == "git" {
		if err := runGitCollect(); err != nil {
			log.Fatalf("Oh jeez: %s\n", err)
		}
	} else {
		exitHelp()
	}

	fmt.Println("Post-mortems successfully parsed and organized")
}

func exitHelp() {
	fmt.Println("What do you want me to do? Options are: ")
	fmt.Println("    `check` - for local parsing")
	fmt.Println("    `git` - to parse and commit stats")
	fmt.Println("    `git check` - to check parsing in your github action")
	log.Fatalln("Please give me an argument!")
}

func runGitCollect() error {
	token := githubactions.GetInput("token")
	if token == "" {
		fmt.Println("Want to check that your post-mortem is correctly formatted? Use `./morty check`")
		return fmt.Errorf("no github action token supplied: %w", ErrMissingInput)
	}

	repository := os.Getenv("GITHUB_REPOSITORY")
	if repository == "" {
		return fmt.Errorf("no repository name set as env variable: %w", ErrMissingInput)
	}

	ref := os.Getenv("GITHUB_REF")
	if ref == "" {
		return fmt.Errorf("no ref: %w", ErrMissingInput)
	}

	gitService := mortems.NewGitHubService(token, repository, ref)
	mortemCollector := mortems.NewMortemCollector(gitService)
	return mortemCollector.Collect()
}

func runGitCheck() error {
	token := githubactions.GetInput("token")
	if token == "" {
		fmt.Println("Want to check that your post-mortem is correctly formatted? Use `./morty check`")
		return fmt.Errorf("no github action token supplied: %w", ErrMissingInput)
	}

	repository := os.Getenv("GITHUB_REPOSITORY")
	if repository == "" {
		return fmt.Errorf("no repository name set as env variable: %w", ErrMissingInput)
	}

	ref := os.Getenv("GITHUB_REF")
	if ref == "" {
		return fmt.Errorf("no ref: %w", ErrMissingInput)
	}

	gitService := mortems.NewGitHubService(token, repository, ref)
	mortemCollector := mortems.NewMortemCollector(gitService)
	_, err := mortemCollector.Check()
	return err
}

func runCheck() error {
	repoPath, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("could not get wd: %w", err)
	}

	fileService := mortems.NewLocalFileService(repoPath)
	mortemCollector := mortems.NewMortemCollector(fileService)
	_, err = mortemCollector.Check()
	return err
}
