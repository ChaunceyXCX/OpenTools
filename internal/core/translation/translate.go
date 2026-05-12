package translation

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

type LangPair struct {
	From string `json:"from"`
	To   string `json:"to"`
}

type Translation struct {
	Original    string   `json:"original"`
	Translation string   `json:"translation"`
	From        string   `json:"from"`
	To          string   `json:"to"`
	Confidence  float64  `json:"confidence"`
}

type Manager struct {
	modelsDir string
	models    map[string]bool
}

func NewManager(dataDir string) *Manager {
	dir := filepath.Join(dataDir, "translation-models")
	os.MkdirAll(dir, 0755)
	return &Manager{
		modelsDir: dir,
		models:    make(map[string]bool),
	}
}

func (m *Manager) SupportedPairs() []LangPair {
	return []LangPair{
		{From: "en", To: "zh"},
		{From: "zh", To: "en"},
		{From: "en", To: "ja"},
		{From: "ja", To: "en"},
	}
}

func (m *Manager) Translate(text, from, to string) (*Translation, error) {
	if len(text) > 5000 {
		return nil, fmt.Errorf("text too long: %d chars (max 5000)", len(text))
	}
	return nil, fmt.Errorf("offline translation engine not loaded - requires Bergamot WASM models")
}

func (m *Manager) DownloadModel(from, to string) error {
	modelName := fmt.Sprintf("%s-%s", from, to)
	if m.models[modelName] {
		return nil
	}

	url := fmt.Sprintf(
		"https://files.mozilla.org/firefox/translations-models/release/%s/%s/%s/model.%s-%s.intgemm8.bin",
		modelName, modelName, modelName, from, to,
	)

	resp, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("download model: %w", err)
	}
	defer resp.Body.Close()

	modelPath := filepath.Join(m.modelsDir, modelName+".bin")
	out, err := os.Create(modelPath)
	if err != nil {
		return fmt.Errorf("create model file: %w", err)
	}
	defer out.Close()

	if _, err := out.ReadFrom(resp.Body); err != nil {
		return fmt.Errorf("write model: %w", err)
	}

	m.models[modelName] = true
	log.Printf("[Translation] downloaded model: %s", modelName)
	return nil
}

func (m *Manager) HasModel(from, to string) bool {
	modelName := fmt.Sprintf("%s-%s", from, to)
	return m.models[modelName]
}
