package main

import (
	"encoding/json"
	"log"
	"os"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/config"
	"github.com/go-git/go-git/v5/plumbing/transport/ssh"
	"github.com/go-git/go-git/v5/storage/memory"
	ssh2 "golang.org/x/crypto/ssh"

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
	conf, err := pkey.ClientConfig()
	if err != nil {
		log.Fatalln(err)
	}

	conf.HostKeyCallback = ssh2.InsecureIgnoreHostKey()

	refs, err := repo.List(
		&git.ListOptions{
			Auth: pkey,
		},
	)
	if err != nil {
		log.Fatalln(err)
	}

	output := make(models.CheckOutput, len(refs))

	for i, ref := range refs {
		output[i] = models.GitBranch{
			Branch: ref.Name().Short(),
			Ref:    ref.Hash().String(),
		}
	}

	json.NewEncoder(os.Stdout).Encode(output)
}
