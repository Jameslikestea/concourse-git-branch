package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/config"
	"github.com/go-git/go-git/v5/plumbing/transport/ssh"
	"github.com/go-git/go-git/v5/storage/memory"

	"github.com/Jameslikestea/concourse-git-branch/models"
)

func main() {
	var req models.CheckRequest
	err := json.NewDecoder(os.Stdin).Decode(&req)
	if err != nil {
		log.Fatalln(err)
	}

	repo := git.NewRemote(
		memory.NewStorage(), &config.RemoteConfig{
			Name: "origin",
			URLs: []string{req.Source.URI},
		},
	)

	pkey, err := ssh.NewPublicKeys("git", []byte(req.Source.PrivateKey), "")
	if err != nil {
		log.Fatalln(err)
	}

	refs, err := repo.List(
		&git.ListOptions{
			Auth: pkey,
		},
	)
	if err != nil {
		log.Fatalln(err)
	}

	for _, ref := range refs {
		fmt.Fprintf(os.Stderr, "%s: %s", ref.Name().Short(), ref.Hash().String())
	}

	fmt.Printf("[]\n")
}
