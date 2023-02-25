package models

type Source struct {
	URI        string `json:"uri"`
	PrivateKey string `json:"private_key"`
}

type GitBranch struct {
	Branch string `json:"branch"`
	Ref    string `json:"ref"`
}
