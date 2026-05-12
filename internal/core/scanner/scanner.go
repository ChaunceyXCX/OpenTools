package scanner

import (
	"runtime"
)

func ScanApplications() ([]Command, error) {
	switch runtime.GOOS {
	case "linux":
		return ScanLinuxApplications()
	case "darwin":
		return ScanDarwinApplications()
	case "windows":
		return ScanWindowsApplications()
	default:
		return nil, nil
	}
}
