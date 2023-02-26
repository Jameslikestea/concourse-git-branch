package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/transport/ssh"

	"github.com/Jameslikestea/concourse-git-branch/models"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintf(os.Stderr, "usage: %s <destination>\n", os.Args[0])
		os.Exit(1)
	}
	destination := os.Args[1]

	var req models.CheckRequest
	err := json.NewDecoder(os.Stdin).Decode(&req)
	if err != nil {
		log.Fatalln(err)
	}

	fmt.Fprintf(os.Stderr, "Destination: %s\n", destination)

	pkey, err := ssh.NewPublicKeys("git", []byte(req.Source.PrivateKey), "")
	if err != nil {
		log.Fatalln(err)
	}

	repo, err := git.PlainClone(
		destination, false, &git.CloneOptions{
			Auth:       pkey,
			URL:        req.Source.URI,
			RemoteName: "origin",
		},
	)
	if err != nil {
		log.Fatalln(err)
	}

	w, err := repo.Worktree()
	if err != nil {
		log.Fatalln(err)
	}

	err = w.Checkout(
		&git.CheckoutOptions{
			Hash: plumbing.NewHash(req.Version.Ref),
		},
	)
	if err != nil {
		log.Fatalln(err)
	}

	c, err := repo.CommitObject(plumbing.NewHash(req.Version.Ref))
	if err != nil {
		log.Fatalln(err)
	}

	output := models.InOutput{
		Version:  req.Version,
		Metadata: []models.KVPair{},
	}

	output.Metadata = append(output.Metadata, models.KVPair{"Commit Hash", c.Hash.String()})
	output.Metadata = append(output.Metadata, models.KVPair{"Author Name", c.Author.Name})
	output.Metadata = append(output.Metadata, models.KVPair{"Author Email", c.Author.Email})
	output.Metadata = append(output.Metadata, models.KVPair{"Commit Message", c.Message})

	os.MkdirAll(destination+"/.metadata", 0755)
	os.WriteFile(destination+"/.metadata/commit", []byte(c.Hash.String()), 0644)
	os.WriteFile(destination+"/.metadata/author_name", []byte(c.Author.Name), 0644)
	os.WriteFile(destination+"/.metadata/author_email", []byte(c.Author.Email), 0644)
	os.WriteFile(destination+"/.metadata/message", []byte(c.Message), 0644)
	os.WriteFile(destination+"/.metadata/branch", []byte(req.Version.Branch), 0644)

	json.NewEncoder(os.Stdout).Encode(output)
}
