package launcher

import (
	"os/exec"
	"runtime"
	"strings"
)

func LaunchApp(appPath string) error {
	executable, args := parseCommandString(appPath)
	if executable == "" {
		return nil
	}
	switch runtime.GOOS {
	case "darwin":
		return launchDarwin(executable, args)
	case "windows":
		return launchWindows(executable, args)
	default:
		return launchDefault(executable, args)
	}
}

func launchDefault(executable string, args []string) error {
	return exec.Command(executable, args...).Start()
}

func launchDarwin(executable string, args []string) error {
	a := append([]string{"-a", executable}, args...)
	return exec.Command("open", a...).Start()
}

func launchWindows(executable string, args []string) error {
	a := append([]string{"/C", "start", `""`, executable})
	a = append(a, args...)
	return exec.Command("cmd", a...).Start()
}

func parseCommandString(cmd string) (string, []string) {
	var parts []string
	var current strings.Builder
	inQuote := byte(0)
	for i := 0; i < len(cmd); i++ {
		ch := cmd[i]
		if inQuote != 0 {
			if ch == inQuote {
				inQuote = 0
			} else {
				current.WriteByte(ch)
			}
		} else if ch == '"' || ch == '\'' {
			inQuote = ch
		} else if ch == ' ' || ch == '\t' {
			if current.Len() > 0 {
				parts = append(parts, current.String())
				current.Reset()
			}
		} else {
			current.WriteByte(ch)
		}
	}
	if current.Len() > 0 {
		parts = append(parts, current.String())
	}
	if len(parts) == 0 {
		return "", nil
	}
	return parts[0], parts[1:]
}


