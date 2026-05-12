package api

import clipboardCore "github.com/ChaunceyXCX/OpenTools/internal/core/clipboard"

type ClipboardService struct {
	monitor *clipboardCore.Monitor
}

func NewClipboardService() *ClipboardService {
	svc := &ClipboardService{}
	svc.monitor = clipboardCore.NewMonitor(svc)
	return svc
}

func (s *ClipboardService) PutClipboardItem(item *clipboardCore.ClipboardItem) error {
	return NewDatabaseService().PutClipboardItem(item)
}

func (s *ClipboardService) GetHistory(page, pageSize int) []*clipboardCore.ClipboardItem {
	return NewDatabaseService().GetClipboardHistory(page, pageSize)
}

func (s *ClipboardService) Delete(id string) bool {
	return NewDatabaseService().DeleteClipboardItem(id)
}

func (s *ClipboardService) Clear() int {
	return NewDatabaseService().ClearClipboard()
}

func (s *ClipboardService) StartMonitor() {
	s.monitor.Start()
}

func (s *ClipboardService) StopMonitor() {
	s.monitor.Stop()
}
