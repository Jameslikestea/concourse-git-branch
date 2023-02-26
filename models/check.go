package models

type CheckRequest struct {
	Source  Source    `json:"source"`
	Version GitBranch `json:"version"`
}

type CheckOutput []GitBranch
