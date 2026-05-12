package scanner

import (
	"os"
	"path/filepath"
	"testing"
)

func writeDesktop(t *testing.T, dir, name, content string) string {
	t.Helper()
	path := filepath.Join(dir, name)
	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}
	return path
}

func TestParseDesktopFile(t *testing.T) {
	t.Setenv("LANG", "en_US.UTF-8")
	content := `[Desktop Entry]
Type=Application
Name=TestApp
Name[zh_CN]=测试应用
Exec=/usr/bin/testapp
Icon=testapp
NoDisplay=false
`
	dir := t.TempDir()
	path := writeDesktop(t, dir, "test.desktop", content)

	cmd := parseDesktopFileToCommand(path)
	if cmd == nil {
		t.Fatal("expected command, got nil")
	}
	if cmd.Name != "TestApp" {
		t.Fatalf("Name=%q, want TestApp", cmd.Name)
	}
	if cmd.Path != "/usr/bin/testapp" {
		t.Fatalf("Path=%q, want /usr/bin/testapp", cmd.Path)
	}
}

func TestParseDesktopFileLocalizedName(t *testing.T) {
	t.Setenv("LANG", "zh_CN.UTF-8")
	content := `[Desktop Entry]
Type=Application
Name=TestApp
Name[zh_CN]=测试应用
Exec=/usr/bin/testapp
`
	dir := t.TempDir()
	path := writeDesktop(t, dir, "test.desktop", content)
	cmd := parseDesktopFileToCommand(path)
	if cmd.Name != "测试应用" {
		t.Fatalf("localized name=%q, want 测试应用", cmd.Name)
	}
}

func TestParseDesktopFileFilters(t *testing.T) {
	dir := t.TempDir()

	tests := []struct {
		name    string
		content string
		wantNil bool
	}{
		{"hidden", "[Desktop Entry]\nType=Application\nHidden=true\nName=H\nExec=/bin/h\n", true},
		{"no display", "[Desktop Entry]\nType=Application\nNoDisplay=true\nName=N\nExec=/bin/n\n", true},
		{"no exec", "[Desktop Entry]\nType=Application\nName=X\n", true},
		{"valid", "[Desktop Entry]\nType=Application\nName=V\nExec=/bin/v\n", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			path := writeDesktop(t, dir, tt.name+".desktop", tt.content)
			cmd := parseDesktopFileToCommand(path)
			if tt.wantNil && cmd != nil {
				t.Fatal("expected nil")
			}
			if !tt.wantNil && cmd == nil {
				t.Fatal("expected command")
			}
		})
	}
}

func TestExtractPinyinAcronym(t *testing.T) {
	got := extractPinyinAcronym("微信")
	if got != "wx" {
		t.Fatalf("got %q, want wx", got)
	}
	got2 := extractPinyinAcronym("谷歌浏览器")
	if got2 == "" {
		t.Fatal("expected non-empty pinyin")
	}
}

func TestExtractAcronym(t *testing.T) {
	got := extractAcronym("Visual Studio Code")
	if got != "vsc" {
		t.Fatalf("got %q, want vsc", got)
	}
}

func TestCleanExecCommand(t *testing.T) {
	got := cleanExecCommand("/usr/bin/app %f %u")
	if got != "/usr/bin/app" {
		t.Fatalf("got %q, want /usr/bin/app", got)
	}
}
