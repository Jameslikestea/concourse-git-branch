package main

import (
	"crypto/rand"
	"crypto/sha256"
	"fmt"
)

func main() {
	data := make([]byte, 10)
	rand.Read(data)
	s := fmt.Sprintf("%x", sha256.Sum256(data))

	fmt.Printf("[{\"ref\":\"399a572\",\"branch\":\"main\"},{\"ref\":\"%s\",\"branch\":\"develop\"}]\n", s[:7])
}
