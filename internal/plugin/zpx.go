package plugin

import (
	"bufio"
	"compress/gzip"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

type asarHeader struct {
	Files map[string]asarEntry `json:"files"`
}

type asarEntry struct {
	Files  map[string]asarEntry `json:"files,omitempty"`
	Size   int                  `json:"size,omitempty"`
	Offset int                  `json:"offset,omitempty"`
}

func ExtractZPX(zpxPath, destDir string) error {
	f, err := os.Open(zpxPath)
	if err != nil {
		return fmt.Errorf("open zpx: %w", err)
	}
	defer f.Close()

	magic := make([]byte, 2)
	if _, err := io.ReadFull(f, magic); err != nil {
		return fmt.Errorf("read magic: %w", err)
	}
	f.Seek(0, 0)

	var reader io.Reader
	if magic[0] == 0x1f && magic[1] == 0x8b {
		gr, err := gzip.NewReader(f)
		if err != nil {
			return fmt.Errorf("gzip reader: %w", err)
		}
		defer gr.Close()
		reader = gr
	} else {
		reader = f
	}

	data, err := io.ReadAll(reader)
	if err != nil {
		return fmt.Errorf("read decompressed: %w", err)
	}

	return extractAsar(data, destDir)
}

func extractAsar(data []byte, destDir string) error {
	if len(data) < 4 {
		return fmt.Errorf("asar too short")
	}
	headerLen := binary.LittleEndian.Uint32(data[:4])
	if int(headerLen)+4 > len(data) {
		return fmt.Errorf("asar header length out of range")
	}
	headerJSON := data[4 : 4+headerLen]
	var header asarHeader
	if err := json.Unmarshal(headerJSON, &header); err != nil {
		return fmt.Errorf("asar header parse: %w", err)
	}
	contentStart := 4 + int(headerLen)
	return extractAsarFiles(header.Files, data, contentStart, destDir, "")
}

func extractAsarFiles(files map[string]asarEntry, data []byte, contentStart int, destDir, prefix string) error {
	for name, entry := range files {
		path := filepath.Join(prefix, name)
		fullPath := filepath.Join(destDir, path)
		if entry.Files != nil {
			if err := os.MkdirAll(fullPath, 0755); err != nil {
				return fmt.Errorf("mkdir %s: %w", fullPath, err)
			}
			if err := extractAsarFiles(entry.Files, data, contentStart, destDir, path); err != nil {
				return err
			}
		} else {
			start := contentStart + entry.Offset
			end := start + entry.Size
			if end > len(data) {
				return fmt.Errorf("file %s out of range", path)
			}
			if err := os.MkdirAll(filepath.Dir(fullPath), 0755); err != nil {
				return fmt.Errorf("mkdir %s: %w", fullPath, err)
			}
			if err := os.WriteFile(fullPath, data[start:end], 0644); err != nil {
				return fmt.Errorf("write %s: %w", fullPath, err)
			}
		}
	}
	return nil
}

func PackZPX(srcDir, zpxPath string) error {
	out, err := os.Create(zpxPath)
	if err != nil {
		return fmt.Errorf("create zpx: %w", err)
	}
	defer out.Close()

	gw := gzip.NewWriter(out)
	defer gw.Close()

	bw := bufio.NewWriter(gw)
	defer bw.Flush()

	if err := packAsar(srcDir, bw); err != nil {
		return fmt.Errorf("pack asar: %w", err)
	}
	return nil
}

type fileEntry struct {
	path     string
	size     int
}

func packAsar(srcDir string, w io.Writer) error {
	var entries []fileEntry
	var offset int
	var buildTree func(prefix string) (map[string]asarEntry, error)
	buildTree = func(prefix string) (map[string]asarEntry, error) {
		dirEntries, err := os.ReadDir(filepath.Join(srcDir, prefix))
		if err != nil {
			return nil, err
		}
		files := make(map[string]asarEntry)
		for _, e := range dirEntries {
			fullPrefix := filepath.Join(prefix, e.Name())
			if e.IsDir() {
				subFiles, err := buildTree(fullPrefix)
				if err != nil {
					return nil, err
				}
				files[e.Name()] = asarEntry{Files: subFiles}
			} else {
				data, err := os.ReadFile(filepath.Join(srcDir, fullPrefix))
				if err != nil {
					return nil, err
				}
				sz := len(data)
				files[e.Name()] = asarEntry{Size: sz, Offset: offset}
				offset += sz
				entries = append(entries, fileEntry{path: fullPrefix, size: sz})
			}
		}
		return files, nil
	}
	tree, err := buildTree(".")
	if err != nil {
		return err
	}
	header := asarHeader{Files: tree}
	headerJSON, err := json.Marshal(header)
	if err != nil {
		return err
	}
	headerLen := make([]byte, 4)
	binary.LittleEndian.PutUint32(headerLen, uint32(len(headerJSON)))
	if _, err := w.Write(headerLen); err != nil {
		return err
	}
	if _, err := w.Write(headerJSON); err != nil {
		return err
	}
	for _, entry := range entries {
		data, err := os.ReadFile(filepath.Join(srcDir, entry.path))
		if err != nil {
			return err
		}
		if _, err := w.Write(data); err != nil {
			return err
		}
	}
	return nil
}
