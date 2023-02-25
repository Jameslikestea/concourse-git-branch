package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/Jameslikestea/concourse-git-branch/models"
)

func main() {
	var req models.CheckRequest
	err := json.NewDecoder(os.Stdin).Decode(&req)
	if err != nil {
		log.Fatalln(err)
	}
	json.NewEncoder(os.Stderr).Encode(req)
	fmt.Printf("[]\n")
}
