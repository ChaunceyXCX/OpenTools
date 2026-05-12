package api

import (
	"log"

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
