package api

import (
	"github.com/ChaunceyXCX/OpenTools/internal/core/zbrowser"
)

type ZBrowserService struct {
	mgr *zbrowser.Manager
}

func NewZBrowserService() *ZBrowserService {
	return &ZBrowserService{
		mgr: zbrowser.NewManager(),
	}
}

func (s *ZBrowserService) Start() (bool, error) {
	if err := s.mgr.Start(); err != nil {
		return false, err
	}
	return true, nil
}

func (s *ZBrowserService) Stop() bool {
	s.mgr.Stop()
	return true
}

func (s *ZBrowserService) Navigate(url string) (*zbrowser.ActionResult, error) {
	return s.mgr.Navigate(url)
}

func (s *ZBrowserService) Screenshot(selector string) (*zbrowser.ActionResult, error) {
	return s.mgr.Screenshot(selector)
}

func (s *ZBrowserService) GetHTML(selector string) (*zbrowser.ActionResult, error) {
	return s.mgr.GetHTML(selector)
}

func (s *ZBrowserService) Click(selector string) (*zbrowser.ActionResult, error) {
	return s.mgr.Click(selector)
}

func (s *ZBrowserService) Type(selector, text string) (*zbrowser.ActionResult, error) {
	return s.mgr.Type(selector, text)
}

func (s *ZBrowserService) Evaluate(js string) (*zbrowser.ActionResult, error) {
	return s.mgr.Evaluate(js)
}

func (s *ZBrowserService) GetCurrentURL() (*zbrowser.ActionResult, error) {
	return s.mgr.GetCurrentURL()
}

func (s *ZBrowserService) GetTitle() (*zbrowser.ActionResult, error) {
	return s.mgr.GetTitle()
}
