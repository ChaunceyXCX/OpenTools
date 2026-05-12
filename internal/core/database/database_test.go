package database

import (
	"path/filepath"
	"testing"
)

func setupDB(t *testing.T) *DB {
	t.Helper()
	dir := t.TempDir()
	db, err := Open(filepath.Join(dir, "test.db"))
	if err != nil {
		t.Fatalf("Open: %v", err)
	}
	t.Cleanup(func() { db.Close() })
	return db
}

func TestPutAndGet(t *testing.T) {
	db := setupDB(t)
	result, err := db.Put(&DbDoc{ID: "ZTOOLS/key1", Data: "value1"})
	if err != nil {
		t.Fatalf("Put: %v", err)
	}
	if !result.Ok {
		t.Fatalf("Put not ok")
	}
	doc, err := db.Get("ZTOOLS/key1")
	if err != nil {
		t.Fatalf("Get: %v", err)
	}
	if doc.Data != "value1" {
		t.Fatalf("Get got %q, want %q", doc.Data, "value1")
	}
}

func TestGetMissing(t *testing.T) {
	db := setupDB(t)
	doc, err := db.Get("ZTOOLS/nonexistent")
	if err != nil {
		t.Fatalf("Get: %v", err)
	}
	if doc != nil {
		t.Fatal("expected nil for missing key")
	}
}

func TestRemove(t *testing.T) {
	db := setupDB(t)
	db.Put(&DbDoc{ID: "ZTOOLS/todelete", Data: "x"})
	result, err := db.Remove("ZTOOLS/todelete")
	if err != nil {
		t.Fatalf("Remove: %v", err)
	}
	if !result.Ok {
		t.Fatal("Remove not ok")
	}
	doc, _ := db.Get("ZTOOLS/todelete")
	if doc != nil {
		t.Fatal("doc should be nil after remove")
	}
}

func TestAllDocsPrefix(t *testing.T) {
	db := setupDB(t)
	db.Put(&DbDoc{ID: "ZTOOLS/a", Data: "1"})
	db.Put(&DbDoc{ID: "ZTOOLS/b", Data: "2"})
	db.Put(&DbDoc{ID: "PLUGIN/x", Data: "3"})

	docs, err := db.AllDocs("ZTOOLS/")
	if err != nil {
		t.Fatalf("AllDocs: %v", err)
	}
	if len(docs) != 2 {
		t.Fatalf("AllDocs got %d, want 2", len(docs))
	}
}

func TestAttachments(t *testing.T) {
	db := setupDB(t)
	_, err := db.PostAttachment("ZTOOLS/att1", []byte("binary"), "text/plain")
	if err != nil {
		t.Fatalf("PostAttachment: %v", err)
	}
	data, err := db.GetAttachment("ZTOOLS/att1")
	if err != nil {
		t.Fatalf("GetAttachment: %v", err)
	}
	if string(data) != "binary" {
		t.Fatalf("got %q, want %q", string(data), "binary")
	}
	ctype, _ := db.GetAttachmentType("ZTOOLS/att1")
	if ctype != "text/plain" {
		t.Fatalf("type got %q, want %q", ctype, "text/plain")
	}
}

func TestBulkDocs(t *testing.T) {
	db := setupDB(t)
	results, err := db.BulkDocs([]*DbDoc{
		{ID: "ZTOOLS/x", Data: "1"},
		{ID: "ZTOOLS/y", Data: "2"},
	})
	if err != nil {
		t.Fatalf("BulkDocs: %v", err)
	}
	if len(results) != 2 {
		t.Fatalf("got %d results, want 2", len(results))
	}
	for _, r := range results {
		if !r.Ok {
			t.Fatalf("result not ok: %+v", r)
		}
	}
}

func TestPersistence(t *testing.T) {
	dir := t.TempDir()
	dbPath := filepath.Join(dir, "persist.db")
	db1, err := Open(dbPath)
	if err != nil {
		t.Fatalf("Open1: %v", err)
	}
	db1.Put(&DbDoc{ID: "ZTOOLS/persist", Data: "saved"})
	db1.Close()

	db2, err := Open(dbPath)
	if err != nil {
		t.Fatalf("Open2: %v", err)
	}
	defer db2.Close()
	doc, _ := db2.Get("ZTOOLS/persist")
	if doc == nil || doc.Data != "saved" {
		t.Fatal("data not persisted")
	}
}

func TestNamespaceIsolation(t *testing.T) {
	db := setupDB(t)
	db.Put(&DbDoc{ID: "ZTOOLS/setting", Data: "x"})
	db.Put(&DbDoc{ID: "PLUGIN/my/setting", Data: "y"})

	ztools, _ := db.AllDocs("ZTOOLS/")
	plugin, _ := db.AllDocs("PLUGIN/my/")

	if len(ztools) != 1 || len(plugin) != 1 {
		t.Fatalf("isolation broken: ztools=%d plugin=%d", len(ztools), len(plugin))
	}
}
