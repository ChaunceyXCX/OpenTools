package api

import (
	"context"
	"log"

	"github.com/ChaunceyXCX/OpenTools/internal/plugin"
)

type PluginService struct {
	system *plugin.SystemPlugin
}

func NewPluginService() *PluginService {
	return &PluginService{
		system: plugin.NewSystemPlugin(),
	}
}

func (s *PluginService) GetSystemManifest() plugin.Manifest {
	return s.system.Manifest()
}

type ExecResult struct {
	Success bool   `json:"success"`
	Data    any    `json:"data,omitempty"`
	Error   string `json:"error,omitempty"`
}

func (s *PluginService) Execute(code string, args map[string]any) ExecResult {
	result, err := s.system.Execute(context.Background(), code, args)
	if err != nil {
		log.Printf("[Plugin] execute %s error: %v", code, err)
		return ExecResult{Success: false, Error: err.Error()}
	}
	return ExecResult{Success: true, Data: result}
}
