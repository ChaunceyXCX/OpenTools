package plugin

import (
	"context"
	"os/exec"
	"runtime"
)

type SystemPlugin struct{}

func NewSystemPlugin() *SystemPlugin {
	return &SystemPlugin{}
}

func (p *SystemPlugin) Manifest() Manifest {
	return Manifest{
		Name:        "system",
		Title:       "System",
		Description: "ZTools system built-in plugin",
		Author:      "Zing",
		Version:     "1.0.0",
		Features: []Feature{
			{Code: "shutdown", Explain: "Shutdown", Cmds: []any{"Shutdown", "shutdown"}},
			{Code: "reboot", Explain: "Reboot", Cmds: []any{"Reboot", "reboot"}},
			{Code: "sleep", Explain: "Sleep", Cmds: []any{"Sleep", "sleep"}},
			{Code: "lock-screen", Explain: "Lock Screen", Cmds: []any{"Lock Screen", "lock screen"}},
			{Code: "logoff", Explain: "Log Off", Cmds: []any{"Log Off", "log off", "Sign Out"}},
		},
	}
}

func (p *SystemPlugin) Execute(ctx context.Context, code string, args map[string]any) (any, error) {
	switch code {
	case "shutdown":
		return nil, p.exec("shutdown", "-h", "now")
	case "reboot":
		return nil, p.exec("reboot")
	case "sleep":
		return nil, p.systemSleep()
	case "lock-screen":
		return nil, p.lockScreen()
	case "logoff":
		return nil, p.logOff()
	}
	return nil, nil
}

func (p *SystemPlugin) exec(name string, args ...string) error {
	return exec.Command(name, args...).Start()
}

func (p *SystemPlugin) systemSleep() error {
	switch runtime.GOOS {
	case "linux":
		return p.exec("systemctl", "suspend")
	case "darwin":
		return p.exec("pmset", "sleepnow")
	case "windows":
		return p.exec("rundll32.exe", "powrprof.dll,SetSuspendState", "Sleep")
	}
	return nil
}

func (p *SystemPlugin) lockScreen() error {
	switch runtime.GOOS {
	case "linux":
		return p.exec("loginctl", "lock-session")
	case "darwin":
		return p.exec("/System/Library/CoreServices/Menu Extras/User.menu/Contents/Resources/CGSession", "-suspend")
	case "windows":
		return p.exec("rundll32.exe", "user32.dll,LockWorkStation")
	}
	return nil
}

func (p *SystemPlugin) logOff() error {
	switch runtime.GOOS {
	case "linux":
		return p.exec("loginctl", "terminate-user", "")
	case "darwin":
		return p.exec("osascript", "-e", `tell app "System Events" to log out`)
	case "windows":
		return p.exec("shutdown", "-l")
	}
	return nil
}
