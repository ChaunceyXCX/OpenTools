package core

import (
	"fmt"
	"os"
	"path/filepath"
	"syscall"
)

type LockFile struct {
	file *os.File
	path string
}

func NewLockFile(dir string) (*LockFile, error) {
	if err := os.MkdirAll(dir, 0755); err != nil {
		return nil, fmt.Errorf("create lock dir: %w", err)
	}

	path := filepath.Join(dir, "ztools.lock")
	file, err := os.OpenFile(path, os.O_CREATE|os.O_RDWR, 0644)
	if err != nil {
		return nil, fmt.Errorf("open lock file: %w", err)
	}

	if err := syscall.Flock(int(file.Fd()), syscall.LOCK_EX|syscall.LOCK_NB); err != nil {
		file.Close()
		return nil, fmt.Errorf("another instance is already running")
	}

	return &LockFile{file: file, path: path}, nil
}

func (l *LockFile) Release() error {
	syscall.Flock(int(l.file.Fd()), syscall.LOCK_UN)
	return l.file.Close()
}
