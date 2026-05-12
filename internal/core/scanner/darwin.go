package scanner

import (
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"
)

func ScanDarwinApplications() ([]Command, error) {
	searchPaths := []string{
		"/Applications",
		"/System/Applications",
		"/System/Applications/Utilities",
	}
	if home, err := os.UserHomeDir(); err == nil {
		searchPaths = append(searchPaths, filepath.Join(home, "Applications"))
	}

	seen := make(map[string]bool)
	var commands []Command

	for _, dir := range searchPaths {
		entries, err := os.ReadDir(dir)
		if err != nil {
			continue
		}
		for _, e := range entries {
			if !e.IsDir() || !strings.HasSuffix(e.Name(), ".app") {
				continue
			}
			appPath := filepath.Join(dir, e.Name())
			name := readPlistString(appPath, "CFBundleName")
			if name == "" {
				name = readPlistString(appPath, "CFBundleDisplayName")
			}
			if name == "" {
				name = strings.TrimSuffix(e.Name(), ".app")
			}
			if seen[name] {
				continue
			}
			seen[name] = true

			execName := readPlistString(appPath, "CFBundleExecutable")
			execPath := appPath
			if execName != "" {
				execPath = filepath.Join(appPath, "Contents", "MacOS", execName)
			}

			icon := filepath.Join(appPath, "Contents", "Resources",
				readPlistString(appPath, "CFBundleIconFile")+".icns")

			commands = append(commands, Command{
				Name:    name,
				Path:    execPath,
				Icon:    "file://" + icon,
				Type:    "direct",
				SubType: "app",
			})
		}
	}

	sort.Slice(commands, func(i, j int) bool {
		return commands[i].Name < commands[j].Name
	})
	return commands, nil
}

func readPlistString(appPath, key string) string {
	plistPath := filepath.Join(appPath, "Contents", "Info.plist")
	if _, err := os.Stat(plistPath); err != nil {
		return ""
	}
	cmd := exec.Command("defaults", "read", plistPath, key)
	out, err := cmd.Output()
	if err != nil {
		return ""
	}
	return strings.TrimSpace(string(out))
}
