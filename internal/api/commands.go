package api

import (
	"log"
	"runtime"

	"github.com/ChaunceyXCX/OpenTools/internal/core/launcher"
	"github.com/ChaunceyXCX/OpenTools/internal/core/scanner"
)

type CommandsService struct{}

func NewCommandsService() *CommandsService {
	return &CommandsService{}
}

func (s *CommandsService) ScanApplications() []scanner.Command {
	if runtime.GOOS != "linux" {
		log.Printf("[Commands] scanning not supported on %s", runtime.GOOS)
		return nil
	}
	commands, err := scanner.ScanLinuxApplications()
	if err != nil {
		log.Printf("[Commands] scan error: %v", err)
		return nil
	}
	return commands
}

type LaunchResult struct {
	Success bool   `json:"success"`
	Error   string `json:"error,omitempty"`
}

func (s *CommandsService) Launch(path string) LaunchResult {
	if err := launcher.LaunchApp(path); err != nil {
		return LaunchResult{Success: false, Error: err.Error()}
	}
	return LaunchResult{Success: true}
}
