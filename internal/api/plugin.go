package api

import (
	"log"
	"os"
	"path/filepath"

	"github.com/ChaunceyXCX/OpenTools/internal/plugin"
)

type PluginService struct {
	runtime *plugin.CompatRuntime
}

func NewPluginService() *PluginService {
	return &PluginService{
		runtime: plugin.NewCompatRuntime(),
	}
}

func (s *PluginService) LoadPlugin(dir string) *plugin.PluginInfo {
	info, err := s.runtime.LoadPluginDir(dir)
	if err != nil {
		log.Printf("[Plugin] load error: %v", err)
		return nil
	}
	return info
}

func (s *PluginService) LoadPluginFromZPX(zpxPath string) *plugin.PluginInfo {
	tmpDir, err := os.MkdirTemp("", "ztools-plugin-*")
	if err != nil {
		log.Printf("[Plugin] temp dir error: %v", err)
		return nil
	}
	defer os.RemoveAll(tmpDir)

	if err := plugin.ExtractZPX(zpxPath, tmpDir); err != nil {
		log.Printf("[Plugin] extract zpx error: %v", err)
		return nil
	}

	entries, err := os.ReadDir(tmpDir)
	if err != nil {
		return nil
	}
	for _, e := range entries {
		if e.IsDir() {
			if _, err := os.Stat(filepath.Join(tmpDir, e.Name(), "plugin.json")); err == nil {
				return s.LoadPlugin(filepath.Join(tmpDir, e.Name()))
			}
		}
	}
	if _, err := os.Stat(filepath.Join(tmpDir, "plugin.json")); err == nil {
		return s.LoadPlugin(tmpDir)
	}
	return nil
}

type ExecResult struct {
	Success bool   `json:"success"`
	Data    any    `json:"data,omitempty"`
	Error   string `json:"error,omitempty"`
}

func (s *PluginService) Execute(name string, featureCode string, args map[string]any) ExecResult {
	result, err := s.runtime.Execute(name, featureCode, args)
	if err != nil {
		log.Printf("[Plugin] execute %s/%s error: %v", name, featureCode, err)
		return ExecResult{Success: false, Error: err.Error()}
	}
	return ExecResult{Success: true, Data: result}
}

func (s *PluginService) ListPlugins() []plugin.PluginInfo {
	return s.runtime.ListPlugins()
}
