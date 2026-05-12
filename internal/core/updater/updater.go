package updater

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"time"
)

const version = "2.4.1"

type ReleaseInfo struct {
	Version     string `json:"version"`
	DownloadURL string `json:"download_url"`
	Notes       string `json:"notes"`
}

func CurrentVersion() string {
	return version
}

func CheckLatest(repo string) (*ReleaseInfo, error) {
	url := fmt.Sprintf("https://api.github.com/repos/%s/releases/latest", repo)
	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Get(url)
	if err != nil {
		return nil, fmt.Errorf("http: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read: %w", err)
	}

	var gh struct {
		TagName string `json:"tag_name"`
		Body    string `json:"body"`
	}
	if err := json.Unmarshal(body, &gh); err != nil {
		return nil, fmt.Errorf("parse: %w", err)
	}

	downloadURL := fmt.Sprintf(
		"https://github.com/%s/releases/download/%s/ztools-%s-%s.zip",
		repo, gh.TagName, gh.TagName, runtime.GOOS,
	)

	return &ReleaseInfo{
		Version:     gh.TagName,
		DownloadURL: downloadURL,
		Notes:       gh.Body,
	}, nil
}

func DownloadUpdate(url, destDir string) (string, error) {
	if err := os.MkdirAll(destDir, 0755); err != nil {
		return "", fmt.Errorf("mkdir: %w", err)
	}

	outPath := filepath.Join(destDir, "ztools-update.zip")
	out, err := os.Create(outPath)
	if err != nil {
		return "", fmt.Errorf("create: %w", err)
	}
	defer out.Close()

	resp, err := http.Get(url)
	if err != nil {
		return "", fmt.Errorf("download: %w", err)
	}
	defer resp.Body.Close()

	written, err := io.Copy(out, resp.Body)
	if err != nil {
		return "", fmt.Errorf("write: %w", err)
	}

	log.Printf("[Updater] downloaded %d bytes to %s", written, outPath)
	return outPath, nil
}
