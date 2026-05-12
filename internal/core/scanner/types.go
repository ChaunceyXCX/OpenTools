package scanner

type Command struct {
	Name        string   `json:"name"`
	Path        string   `json:"path"`
	Icon        string   `json:"icon,omitempty"`
	Aliases     []string `json:"aliases,omitempty"`
	Acronym     string   `json:"acronym,omitempty"`
	Type        string   `json:"type,omitempty"`
	SubType     string   `json:"subType,omitempty"`
	PluginName  string   `json:"pluginName,omitempty"`
	FeatureCode string   `json:"featureCode,omitempty"`
	CmdType     string   `json:"cmdType,omitempty"`
}

type AppScanner interface {
	ScanApplications() ([]Command, error)
}

type MatchCmd struct {
	Name      string `json:"name"`
	Type      string `json:"type"`
	SubType   string `json:"subType,omitempty"`
	Path      string `json:"path,omitempty"`
	CmdType   string `json:"cmdType,omitempty"`
	Match     string `json:"match,omitempty"`
	Label     string `json:"label,omitempty"`
	MinLength int    `json:"minLength,omitempty"`
	Exclude   string `json:"exclude,omitempty"`
}
