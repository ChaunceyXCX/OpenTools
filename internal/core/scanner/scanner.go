package scanner

import (
	"log"
	"runtime"
)

func ScanApplications() ([]Command, error) {
	switch runtime.GOOS {
	case "linux":
		return ScanLinuxApplications()
	case "darwin":
		log.Println("[Scanner] macOS scanning not yet implemented")
		return nil, nil
	case "windows":
		log.Println("[Scanner] Windows scanning not yet implemented")
		return nil, nil
	default:
		return nil, nil
	}
}
