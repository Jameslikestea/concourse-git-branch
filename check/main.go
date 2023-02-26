package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"sort"
	"strings"
	"time"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/config"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
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

	fmt.Fprintf(os.Stderr, "Scanning Github Branches from: %s\n", req.Version.Time.Format(time.RFC3339Nano))

	store := memory.NewStorage()

	repo := git.NewRemote(
		store, &config.RemoteConfig{
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

	err = repo.Fetch(
		&git.FetchOptions{
			Auth:       pkey,
			RemoteName: "origin",
			RefSpecs: []config.RefSpec{
				"+refs/heads/*:refs/remotes/origin/*",
			},
		},
	)
	if err != nil {
		log.Fatalln(err)
	}

	output := models.CheckOutput{
		req.Version,
	}

	iter, err := store.IterReferences()
	iter.ForEach(
		func(reference *plumbing.Reference) error {
			c, err := object.GetCommit(store, reference.Hash())
			if err != nil {
				return err
			}
			if req.Version.Time.Before(c.Author.When) {
				output = append(
					output, models.GitBranch{
						Branch: strings.TrimPrefix(reference.Name().Short(), "origin/"),
						Ref:    reference.Hash().String(),
						Time:   c.Author.When,
					},
				)
			}
			return nil
		},
	)

	sort.Slice(
		output, func(i, j int) bool {
			return output[i].Time.Before(output[j].Time)
		},
	)

	json.NewEncoder(os.Stdout).Encode(output)
}
