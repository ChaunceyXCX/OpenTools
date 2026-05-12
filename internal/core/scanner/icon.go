package scanner

import (
	"os"
	"path/filepath"
	"strings"
)

var iconExtensions = []string{".png", ".svg", ".xpm"}
var preferredSizes = []string{"256x256", "128x128", "64x64", "48x48", "32x32", "scalable"}

func getIconSearchPaths() []string {
	home, err := os.UserHomeDir()
	if err != nil {
		home = "/root"
	}
	return []string{
		filepath.Join(home, ".local", "share", "icons"),
		"/usr/share/icons",
		"/usr/share/pixmaps",
		filepath.Join(home, ".icons"),
		"/usr/local/share/icons",
		"/usr/local/share/pixmaps",
	}
}

func FindIconPath(iconName string) string {
	if iconName == "" {
		return ""
	}
	if strings.HasPrefix(iconName, "/") {
		if _, err := os.Stat(iconName); err == nil {
			return "file://" + iconName
		}
	}
	baseName := strings.TrimSuffix(iconName, filepath.Ext(iconName))
	searchPaths := getIconSearchPaths()
	for _, searchPath := range searchPaths {
		entries, err := os.ReadDir(searchPath)
		if err != nil {
			continue
		}
		for _, entry := range entries {
			if !entry.IsDir() {
				continue
			}
			theme := entry.Name()
			for _, size := range preferredSizes {
				for _, category := range []string{"apps", "applications"} {
					for _, ext := range iconExtensions {
						iconPath := filepath.Join(searchPath, theme, size, category, baseName+ext)
						if _, err := os.Stat(iconPath); err == nil {
							return "file://" + iconPath
						}
					}
				}
			}
		}
		for _, ext := range iconExtensions {
			iconPath := filepath.Join(searchPath, baseName+ext)
			if _, err := os.Stat(iconPath); err == nil {
				return "file://" + iconPath
			}
		}
	}
	return ""
}
