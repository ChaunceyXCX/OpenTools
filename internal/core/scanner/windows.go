package scanner

import (
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/parsiya/golnk"
)

func ScanWindowsApplications() ([]Command, error) {
	searchPaths := []string{}
	if progData := os.Getenv("ProgramData"); progData != "" {
		searchPaths = append(searchPaths,
			filepath.Join(progData, "Microsoft", "Windows", "Start Menu", "Programs"))
	}
	if appData := os.Getenv("APPDATA"); appData != "" {
		searchPaths = append(searchPaths,
			filepath.Join(appData, "Microsoft", "Windows", "Start Menu", "Programs"))
	}
	seen := make(map[string]bool)
	var commands []Command
	for _, dir := range searchPaths {
		filepath.Walk(dir, func(p string, info os.FileInfo, err error) error {
			if err != nil || !strings.HasSuffix(strings.ToLower(p), ".lnk") {
				return nil
			}
			f, err := lnk.File(p)
			if err != nil {
				return nil
			}
			name := strings.TrimSuffix(filepath.Base(p), ".lnk")
			if seen[name] {
				return nil
			}
			targetPath := f.LinkInfo.LocalBasePath + f.LinkInfo.CommonPathSuffix
			if targetPath == "" {
				targetPath = f.StringData.RelativePath
			}
			seen[name] = true
			commands = append(commands, Command{
				Name:    name,
				Path:    targetPath,
				Type:    "direct",
				SubType: "app",
			})
			return nil
		})
	}
	sort.Slice(commands, func(i, j int) bool {
		return commands[i].Name < commands[j].Name
	})
	return commands, nil
}
