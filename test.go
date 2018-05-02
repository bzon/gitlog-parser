package main

import (
	"fmt"
	"os"
	"strings"

	git "gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing/object"
	yaml "gopkg.in/yaml.v2"
)

// DateFormat output format
const DateFormat = "Mon Jan 02,2006 15:04:05 -0700"

// change
type change struct {
	Author      string `yaml:"author"`
	AuthorEmail string `yaml:"author_email"`
	CommitID    string `yaml:"commit_id"`
	Date        string `yaml:"date"`
	Subject     string `yaml:"subject"`
}

// changelog
type changelog struct {
	RepoURL   string `yaml:"repo_url"`
	ParseType string `yaml:"type"`
	From      string `yaml:"from_commit"`
	To        string `yaml:"to_commit"`
	//changes is an array of change
	Changes []change `yaml:"changelogs"`
}

func main() {

	//test vars

	// clone repo
	r, _ := git.PlainOpen(os.Getenv("PWD"))
	list, _ := r.Remotes()
	var remoteurl = ""
	for _, values := range list {
		remoteurl = values.String()
	}

	changelog := changelog{
		// get only first remote url todo
		RepoURL:   strings.Split(strings.Split(remoteurl, "(")[0], "\t")[1],
		ParseType: "yaml",
		From:      "v1",
		To:        "v2",
	}
	ref, _ := r.Head()

	// get git log
	cIter, _ := r.Log(&git.LogOptions{From: ref.Hash()})

	// Add commits -> Changes
	cIter.ForEach(func(commit *object.Commit) error {
		// assign commit to change
		change := change{
			Author:      commit.Author.Name,
			AuthorEmail: commit.Author.Email,
			CommitID:    commit.Hash.String(),
			Date:        commit.Author.When.Format(DateFormat),
			Subject:     commit.Message,
		}

		// append single entry of change and add to Changes
		changelog.Changes = append(changelog.Changes, change)

		return nil
	})

	// marshal test into a yaml
	ymlf, _ := yaml.Marshal(changelog)

	// print yaml
	fmt.Printf("---start---\n%s\n---end---", string(ymlf))

}
