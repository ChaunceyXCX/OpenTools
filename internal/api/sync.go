package api

import (
	"github.com/ChaunceyXCX/OpenTools/internal/core/database"
	"github.com/ChaunceyXCX/OpenTools/internal/core/sync"
)

type SyncBinding struct {
	svc *sync.SyncService
}

func NewSyncBinding() *SyncBinding {
	return &SyncBinding{}
}

func (s *SyncBinding) Init(db *database.DB) {
	s.svc = sync.NewSyncService(db)
}

func (s *SyncBinding) Configure(url, username, password, dir string) {
	s.svc.Configure(sync.Config{
		URL: url, Username: username,
		Password: password, Dir: dir,
	})
}

func (s *SyncBinding) Push() (map[string]any, error) {
	return s.svc.Push()
}

func (s *SyncBinding) Pull() (map[string]any, error) {
	return s.svc.Pull()
}
