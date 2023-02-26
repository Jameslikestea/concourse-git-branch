package models

type InRequest struct {
	Source  Source    `json:"source"`
	Version GitBranch `json:"version"`
}

type InOutput struct {
	Version  GitBranch `json:"version"`
	Metadata []KVPair  `json:"metadata"`
}

type KVPair struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}
