package scanner

import (
	"bufio"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"unicode"

	"github.com/mozillazg/go-pinyin"
)

type desktopEntry struct {
	values map[string]string
}

func parseDesktopFile(path string) *desktopEntry {
	f, err := os.Open(path)
	if err != nil {
		return nil
	}
	defer f.Close()

	entry := &desktopEntry{values: make(map[string]string)}
	inDesktopEntry := false
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "[Desktop Entry]" {
			inDesktopEntry = true
			continue
		}
		if strings.HasPrefix(line, "[") && strings.HasSuffix(line, "]") && inDesktopEntry {
			break
		}
		if !inDesktopEntry || line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		eqIdx := strings.IndexByte(line, '=')
		if eqIdx == -1 {
			continue
		}
		key := strings.TrimSpace(line[:eqIdx])
		value := strings.TrimSpace(line[eqIdx+1:])
		entry.values[key] = value
	}
	return entry
}

func (e *desktopEntry) get(key string) string {
	return e.values[key]
}

func (e *desktopEntry) getLocalizedName() string {
	lang := os.Getenv("LANG")
	if lang == "" {
		lang = os.Getenv("LANGUAGE")
	}
	langCode := strings.Split(lang, ".")[0]
	parts := strings.Split(langCode, "_")
	langBase := ""
	if len(parts) > 0 {
		langBase = parts[0]
	}

	candidates := []string{}
	if langCode != "" {
		candidates = append(candidates, "Name["+langCode+"]")
	}
	if langBase != "" {
		candidates = append(candidates, "Name["+langBase+"]")
	}
	candidates = append(candidates, "Name")

	for _, key := range candidates {
		if v := e.values[key]; v != "" {
			return strings.TrimSpace(v)
		}
	}
	return ""
}

func cleanExecCommand(exec string) string {
	result := exec
	for _, c := range []string{"%f", "%F", "%u", "%U", "%d", "%D", "%n", "%N", "%i", "%c", "%k", "%v", "%m"} {
		result = strings.ReplaceAll(result, c, "")
	}
	result = strings.Join(strings.Fields(result), " ")
	return strings.TrimSpace(result)
}

func extractAcronym(name string) string {
	var result strings.Builder
	for _, r := range name {
		if unicode.IsUpper(r) {
			result.WriteRune(unicode.ToLower(r))
		}
	}
	if result.Len() > 0 {
		return result.String()
	}
	words := strings.Fields(name)
	for _, w := range words {
		if len(w) > 0 {
			result.WriteString(strings.ToLower(string(w[0])))
		}
	}
	return result.String()
}

func extractPinyinAcronym(name string) string {
	var result strings.Builder
	for _, r := range name {
		if unicode.Is(unicode.Han, r) {
			py := pinyin.Pinyin(string(r), pinyin.Args{Style: pinyin.FirstLetter})
			if len(py) > 0 && len(py[0]) > 0 {
				result.WriteString(string(py[0][0]))
			}
		} else if unicode.IsLetter(r) {
			result.WriteRune(unicode.ToLower(r))
		}
	}
	return result.String()
}

func isChinese(s string) bool {
	for _, r := range s {
		if unicode.Is(unicode.Han, r) {
			return true
		}
	}
	return false
}

func getLinuxDesktopPaths() []string {
	home, err := os.UserHomeDir()
	if err != nil {
		home = "/root"
	}
	xdgDataDirs := os.Getenv("XDG_DATA_DIRS")
	if xdgDataDirs == "" {
		xdgDataDirs = "/usr/local/share:/usr/share"
	}
	baseDirs := strings.Split(xdgDataDirs, ":")
	paths := []string{filepath.Join(home, ".local", "share", "applications")}
	for _, dir := range baseDirs {
		if dir != "" {
			paths = append(paths, filepath.Join(dir, "applications"))
		}
	}
	seen := make(map[string]bool)
	unique := []string{}
	for _, p := range paths {
		if !seen[p] {
			seen[p] = true
			unique = append(unique, p)
		}
	}
	return unique
}

func scanDesktopDir(dirPath string) []string {
	entries, err := os.ReadDir(dirPath)
	if err != nil {
		return nil
	}
	var files []string
	for _, e := range entries {
		if !e.IsDir() && strings.HasSuffix(e.Name(), ".desktop") {
			files = append(files, filepath.Join(dirPath, e.Name()))
		}
	}
	return files
}

func parseDesktopFileToCommand(desktopPath string) *Command {
	entry := parseDesktopFile(desktopPath)
	if entry == nil {
		return nil
	}
	if entry.get("Type") != "Application" ||
		entry.get("NoDisplay") == "true" ||
		entry.get("Hidden") == "true" ||
		entry.get("Exec") == "" {
		return nil
	}
	name := entry.getLocalizedName()
	if name == "" {
		return nil
	}
	exec := cleanExecCommand(entry.get("Exec"))
	if exec == "" {
		return nil
	}
	iconPath := ""
	if iconName := entry.get("Icon"); iconName != "" {
		iconPath = FindIconPath(iconName)
	}
	aliases := []string{}
	rawName := entry.get("Name")
	if rawName != "" && rawName != name {
		aliases = append(aliases, rawName)
	}
	acronym := extractAcronym(name)
	if acronym != "" {
		if rawName != "" {
			if a2 := extractAcronym(rawName); a2 != "" && a2 != acronym {
				aliases = append(aliases, a2)
			}
		}
	}
	if isChinese(name) {
		py := extractPinyinAcronym(name)
		if py != "" {
			aliases = append(aliases, py)
		}
		if rawName != "" && isChinese(rawName) {
			if py2 := extractPinyinAcronym(rawName); py2 != "" && py2 != py {
				aliases = append(aliases, py2)
			}
		}
	}
	return &Command{
		Name:    name,
		Path:    exec,
		Icon:    iconPath,
		Aliases: aliases,
		Acronym: acronym,
	}
}

func ScanLinuxApplications() ([]Command, error) {
	seen := make(map[string]bool)
	var commands []Command
	for _, dirPath := range getLinuxDesktopPaths() {
		for _, f := range scanDesktopDir(dirPath) {
			cmd := parseDesktopFileToCommand(f)
			if cmd == nil {
				continue
			}
			if seen[cmd.Name] {
				continue
			}
			seen[cmd.Name] = true
			commands = append(commands, *cmd)
		}
	}
	sort.Slice(commands, func(i, j int) bool {
		return commands[i].Name < commands[j].Name
	})
	return commands, nil
}
