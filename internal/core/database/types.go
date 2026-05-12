package database

const (
	NamespaceZtools = "ZTOOLS"
	PluginPrefix    = "PLUGIN"
)

type DbDoc struct {
	ID           string `json:"_id"`
	Rev          string `json:"_rev,omitempty"`
	LastModified int64  `json:"_lastModified,omitempty"`
	CloudSynced  bool   `json:"_cloudSynced,omitempty"`
	Data         string `json:"data,omitempty"`
}

type DbResult struct {
	ID      string `json:"id"`
	Rev     string `json:"rev,omitempty"`
	Ok      bool   `json:"ok,omitempty"`
	Error   bool   `json:"error,omitempty"`
	Name    string `json:"name,omitempty"`
	Message string `json:"message,omitempty"`
}

func GetPluginDataPrefix(pluginName string) string {
	return PluginPrefix + "/" + pluginName + "/"
}
