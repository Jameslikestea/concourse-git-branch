package models

type CheckRequest struct {
	Source Source `json:"source"`
}

type CheckOutput []GitBranch
