package plugin

import "context"

type Manifest struct {
	Name        string    `json:"name"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Author      string    `json:"author"`
	Version     string    `json:"version"`
	Features    []Feature `json:"features"`
}

type Feature struct {
	Code    string `json:"code"`
	Explain string `json:"explain"`
	Icon    string `json:"icon,omitempty"`
	Cmds    []any  `json:"cmds"`
}

type Plugin interface {
	Manifest() Manifest
	Execute(ctx context.Context, code string, args map[string]any) (any, error)
}
