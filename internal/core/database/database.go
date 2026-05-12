package database

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"

	"go.etcd.io/bbolt"
)

var (
	mainBucket       = []byte("main")
	metaBucket       = []byte("meta")
	attachmentBucket = []byte("attachment")
)

type DB struct {
	bolt *bbolt.DB
	mu   sync.Mutex
	path string
}

func Open(dbPath string) (*DB, error) {
	dir := filepath.Dir(dbPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return nil, fmt.Errorf("create db dir: %w", err)
	}

	boltDB, err := bbolt.Open(dbPath, 0600, &bbolt.Options{Timeout: 1 * time.Second})
	if err != nil {
		return nil, fmt.Errorf("open bbolt: %w", err)
	}

	if err := boltDB.Update(func(tx *bbolt.Tx) error {
		for _, name := range [][]byte{mainBucket, metaBucket, attachmentBucket} {
			if _, err := tx.CreateBucketIfNotExists(name); err != nil {
				return fmt.Errorf("create bucket %s: %w", string(name), err)
			}
		}
		return nil
	}); err != nil {
		boltDB.Close()
		return nil, err
	}

	return &DB{bolt: boltDB, path: dbPath}, nil
}

func (db *DB) Close() error {
	return db.bolt.Close()
}

func (db *DB) Path() string {
	return db.path
}

func docKey(id string) []byte {
	return []byte(id)
}

func (db *DB) Get(key string) (*DbDoc, error) {
	var doc *DbDoc
	err := db.bolt.View(func(tx *bbolt.Tx) error {
		b := tx.Bucket(mainBucket)
		val := b.Get(docKey(key))
		if val == nil {
			return nil
		}
		var d DbDoc
		if err := json.Unmarshal(val, &d); err != nil {
			return fmt.Errorf("unmarshal doc: %w", err)
		}
		d.ID = key
		doc = &d
		return nil
	})
	return doc, err
}

func (db *DB) Put(doc *DbDoc) (*DbResult, error) {
	result := &DbResult{ID: doc.ID}
	err := db.bolt.Update(func(tx *bbolt.Tx) error {
		b := tx.Bucket(mainBucket)
		doc.Rev = revFromTime()
		doc.LastModified = time.Now().UnixMilli()

		data, err := json.Marshal(doc)
		if err != nil {
			return fmt.Errorf("marshal doc: %w", err)
		}
		if err := b.Put(docKey(doc.ID), data); err != nil {
			return fmt.Errorf("bbolt put: %w", err)
		}
		result.Rev = doc.Rev
		result.Ok = true
		return nil
	})
	return result, err
}

func (db *DB) Remove(key string) (*DbResult, error) {
	result := &DbResult{ID: key}
	err := db.bolt.Update(func(tx *bbolt.Tx) error {
		b := tx.Bucket(mainBucket)
		if val := b.Get(docKey(key)); val == nil {
			result.Error = true
			result.Name = "not_found"
			result.Message = "Document not found"
			return nil
		}
		if err := b.Delete(docKey(key)); err != nil {
			return fmt.Errorf("bbolt delete: %w", err)
		}
		result.Ok = true
		return nil
	})
	return result, err
}

func (db *DB) AllDocs(prefix string) ([]*DbDoc, error) {
	var docs []*DbDoc
	prefixBytes := []byte(prefix)
	err := db.bolt.View(func(tx *bbolt.Tx) error {
		b := tx.Bucket(mainBucket)
		c := b.Cursor()
		var k, v []byte
		if prefix == "" {
			k, v = c.First()
		} else {
			k, v = c.Seek(prefixBytes)
		}
		for ; k != nil; k, v = c.Next() {
			if prefix != "" && !bytesHasPrefix(k, prefixBytes) {
				break
			}
			var d DbDoc
			if err := json.Unmarshal(v, &d); err != nil {
				continue
			}
			d.ID = string(k)
			docs = append(docs, &d)
		}
		return nil
	})
	return docs, err
}

func (db *DB) BulkDocs(docs []*DbDoc) ([]*DbResult, error) {
	results := make([]*DbResult, len(docs))
	err := db.bolt.Update(func(tx *bbolt.Tx) error {
		b := tx.Bucket(mainBucket)
		for i, doc := range docs {
			r := &DbResult{ID: doc.ID}
			doc.Rev = revFromTime()
			doc.LastModified = time.Now().UnixMilli()
			data, err := json.Marshal(doc)
			if err != nil {
				r.Error = true
				r.Message = err.Error()
				results[i] = r
				continue
			}
			if err := b.Put(docKey(doc.ID), data); err != nil {
				r.Error = true
				r.Message = err.Error()
				results[i] = r
				continue
			}
			r.Ok = true
			r.Rev = doc.Rev
			results[i] = r
		}
		return nil
	})
	return results, err
}

func (db *DB) PostAttachment(id string, data []byte, contentType string) (*DbResult, error) {
	result := &DbResult{ID: id}
	err := db.bolt.Update(func(tx *bbolt.Tx) error {
		b := tx.Bucket(attachmentBucket)
		if err := b.Put(docKey(id), data); err != nil {
			return fmt.Errorf("bbolt put attachment: %w", err)
		}
		// store content type in meta
		mb := tx.Bucket(metaBucket)
		typeKey := "att_type:" + id
		if err := mb.Put([]byte(typeKey), []byte(contentType)); err != nil {
			return fmt.Errorf("bbolt put attachment type: %w", err)
		}
		result.Ok = true
		return nil
	})
	return result, err
}

func (db *DB) GetAttachment(id string) ([]byte, error) {
	var data []byte
	err := db.bolt.View(func(tx *bbolt.Tx) error {
		b := tx.Bucket(attachmentBucket)
		val := b.Get(docKey(id))
		if val == nil {
			return nil
		}
		data = make([]byte, len(val))
		copy(data, val)
		return nil
	})
	return data, err
}

func (db *DB) GetAttachmentType(id string) (string, error) {
	var ctype string
	err := db.bolt.View(func(tx *bbolt.Tx) error {
		mb := tx.Bucket(metaBucket)
		val := mb.Get([]byte("att_type:" + id))
		if val != nil {
			ctype = string(val)
		}
		return nil
	})
	return ctype, err
}

func (db *DB) Rev(id string) (string, error) {
	var rev string
	err := db.bolt.View(func(tx *bbolt.Tx) error {
		mb := tx.Bucket(metaBucket)
		val := mb.Get([]byte("rev:" + id))
		if val != nil {
			rev = string(val)
		}
		return nil
	})
	return rev, err
}

// --- helpers ---

func revFromTime() string {
	return fmt.Sprintf("%d", time.Now().UnixNano())
}

func bytesHasPrefix(b, prefix []byte) bool {
	return len(b) >= len(prefix) && string(b[:len(prefix)]) == string(prefix)
}
