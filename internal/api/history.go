package api

import (
	"encoding/json"
	"log"
	"time"

	"github.com/ChaunceyXCX/OpenTools/internal/core/database"
)

type HistoryEntry struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	Path      string `json:"path"`
	Timestamp int64  `json:"timestamp"`
	Pinned    bool   `json:"pinned"`
}

type HistoryService struct{}

func NewHistoryService() *HistoryService {
	return &HistoryService{}
}

func (s *HistoryService) AddEntry(name, path string) {
	if globalDB == nil {
		return
	}
	id := "HISTORY/" + name
	entry := &HistoryEntry{
		ID: id, Name: name, Path: path,
		Timestamp: time.Now().UnixMilli(),
	}
	existing := s.GetEntry(name)
	if existing != nil {
		entry.Pinned = existing.Pinned
	}
	data, _ := json.Marshal(entry)
	globalDB.Put(&database.DbDoc{ID: id, Data: string(data)})
}

func (s *HistoryService) GetEntry(name string) *HistoryEntry {
	if globalDB == nil {
		return nil
	}
	doc, err := globalDB.Get("HISTORY/" + name)
	if err != nil || doc == nil {
		return nil
	}
	var e HistoryEntry
	json.Unmarshal([]byte(doc.Data), &e)
	return &e
}

func (s *HistoryService) GetHistory(limit int) []HistoryEntry {
	if globalDB == nil {
		return nil
	}
	docs, err := globalDB.AllDocs("HISTORY/")
	if err != nil {
		return nil
	}
	var entries []HistoryEntry
	for _, doc := range docs {
		var e HistoryEntry
		if err := json.Unmarshal([]byte(doc.Data), &e); err != nil {
			continue
		}
		entries = append(entries, e)
	}
	sortByTimestampDesc(entries)
	if limit > 0 && len(entries) > limit {
		entries = entries[:limit]
	}
	return entries
}

func (s *HistoryService) TogglePin(name string) {
	e := s.GetEntry(name)
	if e == nil {
		return
	}
	e.Pinned = !e.Pinned
	data, _ := json.Marshal(e)
	globalDB.Put(&database.DbDoc{ID: "HISTORY/" + name, Data: string(data)})
}

func (s *HistoryService) RemoveEntry(name string) {
	if globalDB == nil {
		return
	}
	globalDB.Remove("HISTORY/" + name)
}

func (s *HistoryService) ClearHistory() {
	if globalDB == nil {
		return
	}
	docs, err := globalDB.AllDocs("HISTORY/")
	if err != nil {
		return
	}
	for _, doc := range docs {
		globalDB.Remove(doc.ID)
	}
	log.Printf("[History] cleared %d entries", len(docs))
}

func sortByTimestampDesc(entries []HistoryEntry) {
	n := len(entries)
	for i := 0; i < n; i++ {
		for j := i + 1; j < n; j++ {
			if entries[j].Timestamp > entries[i].Timestamp {
				entries[i], entries[j] = entries[j], entries[i]
			}
		}
	}
}
