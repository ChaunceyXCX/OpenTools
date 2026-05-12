package launcher

import (
	"testing"
)

func TestParseNoArgs(t *testing.T) {
	exe, args := parseCommandString("/usr/bin/ls")
	if exe != "/usr/bin/ls" {
		t.Fatalf("exe=%q", exe)
	}
	if len(args) != 0 {
		t.Fatalf("expected 0 args, got %d", len(args))
	}
}

func TestParseWithArgs(t *testing.T) {
	exe, args := parseCommandString("/usr/bin/ls -la /tmp")
	if exe != "/usr/bin/ls" {
		t.Fatalf("exe=%q", exe)
	}
	if len(args) != 2 || args[0] != "-la" || args[1] != "/tmp" {
		t.Fatalf("args=%v", args)
	}
}

func TestParseQuotedPath(t *testing.T) {
	exe, args := parseCommandString(`"/usr/bin/my app" --flag`)
	if exe != "/usr/bin/my app" {
		t.Fatalf("exe=%q", exe)
	}
	if len(args) != 1 || args[0] != "--flag" {
		t.Fatalf("args=%v", args)
	}
}

func TestParseEmpty(t *testing.T) {
	exe, _ := parseCommandString("")
	if exe != "" {
		t.Fatal("expected empty")
	}
}
