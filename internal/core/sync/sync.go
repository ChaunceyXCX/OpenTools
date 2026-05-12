package sync

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/ChaunceyXCX/OpenTools/internal/core/database"
	"github.com/studio-b12/gowebdav"
)

type Config struct {
	URL      string `json:"url"`
	Username string `json:"username"`
	Password string `json:"password"`
	Dir      string `json:"dir"`
}

type Engine struct {
	client *gowebdav.Client
	config Config
	db     *database.DB
	path   string
}

func NewEngine(db *database.DB, cfg Config) *Engine {
	client := gowebdav.NewClient(cfg.URL, cfg.Username, cfg.Password)
	return &Engine{
		client: client,
		config: cfg,
		db:     db,
		path:   cfg.Dir + "/ztools-backup.json",
	}
}

func (e *Engine) Push() error {
	docs, err := e.db.AllDocs("ZTOOLS/")
	if err != nil {
		return fmt.Errorf("read docs: %w", err)
	}
	data, err := json.Marshal(docs)
	if err != nil {
		return fmt.Errorf("marshal: %w", err)
	}
	if err := e.client.Write(e.path, data, 0644); err != nil {
		return fmt.Errorf("webdav write: %w", err)
	}
	log.Printf("[Sync] pushed %d docs to %s", len(docs), e.path)
	return nil
}

func (e *Engine) Pull() error {
	data, err := e.client.Read(e.path)
	if err != nil {
		return fmt.Errorf("webdav read: %w", err)
	}
	var docs []*database.DbDoc
	if err := json.Unmarshal(data, &docs); err != nil {
		return fmt.Errorf("unmarshal: %w", err)
	}
	for _, doc := range docs {
		if _, err := e.db.Put(doc); err != nil {
			log.Printf("[Sync] put error for %s: %v", doc.ID, err)
		}
	}
	log.Printf("[Sync] pulled %d docs", len(docs))
	return nil
}

type SyncService struct {
	engine *Engine
	config Config
	db     *database.DB
}

func NewSyncService(db *database.DB) *SyncService {
	return &SyncService{db: db}
}

func (s *SyncService) Configure(cfg Config) {
	s.config = cfg
	s.engine = NewEngine(s.db, cfg)
}

func (s *SyncService) Push() (map[string]any, error) {
	if s.engine == nil {
		return nil, fmt.Errorf("not configured")
	}
	start := time.Now()
	err := s.engine.Push()
	duration := time.Since(start).Milliseconds()
	return map[string]any{"success": err == nil, "duration_ms": duration}, err
}

func (s *SyncService) Pull() (map[string]any, error) {
	if s.engine == nil {
		return nil, fmt.Errorf("not configured")
	}
	start := time.Now()
	err := s.engine.Pull()
	duration := time.Since(start).Milliseconds()
	return map[string]any{"success": err == nil, "duration_ms": duration}, err
}
