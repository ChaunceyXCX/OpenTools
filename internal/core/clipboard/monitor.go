package clipboard

import (
	"crypto/sha256"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/atotto/clipboard"
)

const (
	pollInterval = 500 * time.Millisecond
	Prefix       = "CLIPBOARD/"
)

type ClipboardItem struct {
	ID        string `json:"id"`
	Content   string `json:"content"`
	Type      string `json:"type"`
	Timestamp int64  `json:"timestamp"`
	Hash      string `json:"hash"`
}

type Store interface {
	PutClipboardItem(item *ClipboardItem) error
}

type Monitor struct {
	mu       sync.Mutex
	running  bool
	lastHash string
	stopCh   chan struct{}
	store    Store
}

func NewMonitor(store Store) *Monitor {
	return &Monitor{
		stopCh: make(chan struct{}),
		store:  store,
	}
}

func (m *Monitor) Start() {
	m.mu.Lock()
	if m.running {
		m.mu.Unlock()
		return
	}
	m.running = true
	m.mu.Unlock()

	go func() {
		ticker := time.NewTicker(pollInterval)
		defer ticker.Stop()
		log.Println("[Clipboard] monitor started")
		for {
			select {
			case <-ticker.C:
				m.check()
			case <-m.stopCh:
				log.Println("[Clipboard] monitor stopped")
				return
			}
		}
	}()
}

func (m *Monitor) Stop() {
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.running {
		close(m.stopCh)
		m.running = false
	}
}

func (m *Monitor) check() {
	text, err := clipboard.ReadAll()
	if err != nil || text == "" {
		return
	}
	hash := fmt.Sprintf("%x", sha256.Sum256([]byte(text)))
	if hash == m.lastHash {
		return
	}
	m.lastHash = hash
	item := &ClipboardItem{
		ID:        Prefix + fmt.Sprintf("%d", time.Now().UnixNano()),
		Content:   truncate(text, 50000),
		Type:      "text",
		Timestamp: time.Now().UnixMilli(),
		Hash:      hash,
	}
	if err := m.store.PutClipboardItem(item); err != nil {
		log.Printf("[Clipboard] store error: %v", err)
	}
}

func truncate(s string, max int) string {
	if len(s) <= max {
		return s
	}
	return s[:max]
}
