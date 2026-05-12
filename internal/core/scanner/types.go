package scanner

type Command struct {
	Name    string   `json:"name"`
	Path    string   `json:"path"`
	Icon    string   `json:"icon,omitempty"`
	Aliases []string `json:"aliases,omitempty"`
	Acronym string   `json:"acronym,omitempty"`
}

type AppScanner interface {
	ScanApplications() ([]Command, error)
}
