package api

type Command struct {
	ID      string `json:"id"`
	Name    string `json:"name"`
	Type    string `json:"type"`
	SubType string `json:"subType,omitempty"`
	Icon    string `json:"icon,omitempty"`
}

type ClipboardItem struct {
	ID        string `json:"id"`
	Content   string `json:"content"`
	Type      string `json:"type"`
	Timestamp int64  `json:"timestamp"`
}

type PluginManifest struct {
	Name        string `json:"name"`
	Version     string `json:"version"`
	Description string `json:"description"`
	Author      string `json:"author"`
	Main        string `json:"main,omitempty"`
	Features    []Feature `json:"features,omitempty"`
}

type Feature struct {
	Code string   `json:"code"`
	Explain string `json:"explain"`
	Cmds  []string `json:"cmds"`
}
