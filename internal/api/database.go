package api

import (
	"encoding/json"
	"log"
	"sort"

	"github.com/ChaunceyXCX/OpenTools/internal/core/clipboard"
	"github.com/ChaunceyXCX/OpenTools/internal/core/database"
)

var globalDB *database.DB

type DatabaseService struct{}

func NewDatabaseService() *DatabaseService {
	return &DatabaseService{}
}

func InitDatabase(dbPath string) error {
	var err error
	globalDB, err = database.Open(dbPath)
	if err != nil {
		return err
	}
	log.Printf("[DB] opened at %s", dbPath)
	return nil
}

func CloseDatabase() {
	if globalDB != nil {
		globalDB.Close()
	}
}

func (s *DatabaseService) Get(key string) *database.DbDoc {
	doc, err := globalDB.Get(key)
	if err != nil {
		log.Printf("[DB] Get error: %v", err)
		return nil
	}
	return doc
}

func (s *DatabaseService) Put(id string, data string) *database.DbResult {
	doc := &database.DbDoc{
		ID:   id,
		Data: data,
	}
	result, err := globalDB.Put(doc)
	if err != nil {
		log.Printf("[DB] Put error: %v", err)
		return &database.DbResult{ID: id, Error: true, Message: err.Error()}
	}
	return result
}

func (s *DatabaseService) Remove(key string) *database.DbResult {
	result, err := globalDB.Remove(key)
	if err != nil {
		log.Printf("[DB] Remove error: %v", err)
		return &database.DbResult{ID: key, Error: true, Message: err.Error()}
	}
	return result
}

func (s *DatabaseService) AllDocs(prefix string) []*database.DbDoc {
	docs, err := globalDB.AllDocs(prefix)
	if err != nil {
		log.Printf("[DB] AllDocs error: %v", err)
		return nil
	}
	return docs
}

func (s *DatabaseService) PutClipboardItem(item *clipboard.ClipboardItem) error {
	if globalDB == nil {
		return nil
	}
	data, err := json.Marshal(item)
	if err != nil {
		return err
	}
	doc := &database.DbDoc{ID: item.ID, Data: string(data)}
	_, err = globalDB.Put(doc)
	return err
}

func (s *DatabaseService) GetClipboardHistory(page, pageSize int) []*clipboard.ClipboardItem {
	if globalDB == nil {
		return nil
	}
	docs, err := globalDB.AllDocs(clipboard.Prefix)
	if err != nil {
		return nil
	}
	var items []*clipboard.ClipboardItem
	for _, doc := range docs {
		var item clipboard.ClipboardItem
		if err := json.Unmarshal([]byte(doc.Data), &item); err != nil {
			continue
		}
		items = append(items, &item)
	}
	sort.Slice(items, func(i, j int) bool {
		return items[i].Timestamp > items[j].Timestamp
	})
	if pageSize > 0 {
		start := (page - 1) * pageSize
		if start >= len(items) {
			return nil
		}
		end := start + pageSize
		if end > len(items) {
			end = len(items)
		}
		return items[start:end]
	}
	return items
}

func (s *DatabaseService) DeleteClipboardItem(id string) bool {
	if globalDB == nil {
		return false
	}
	result, err := globalDB.Remove(id)
	return err == nil && result.Ok
}

func (s *DatabaseService) ClearClipboard() int {
	if globalDB == nil {
		return 0
	}
	docs, err := globalDB.AllDocs(clipboard.Prefix)
	if err != nil {
		return 0
	}
	count := 0
	for _, doc := range docs {
		if r, err := globalDB.Remove(doc.ID); err == nil && r.Ok {
			count++
		}
	}
	return count
}
