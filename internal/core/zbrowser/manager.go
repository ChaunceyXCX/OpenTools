package zbrowser

import (
	"context"
	"encoding/base64"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/chromedp/chromedp"
)

type Tab struct {
	ID      string `json:"id"`
	URL     string `json:"url"`
	Title   string `json:"title"`
	Created int64  `json:"created"`
}

type ActionResult struct {
	Success  bool   `json:"success"`
	Data     string `json:"data,omitempty"`
	Error    string `json:"error,omitempty"`
	MIMEType string `json:"mimeType,omitempty"`
}

type Manager struct {
	mu       sync.Mutex
	ctx      context.Context
	cancel   context.CancelFunc
	allocCtx context.Context
	allocCancel context.CancelFunc
	tabs     map[string]*Tab
	running  bool
}

func NewManager() *Manager {
	return &Manager{
		tabs: make(map[string]*Tab),
	}
}

func (m *Manager) Start() error {
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.running {
		return nil
	}

	allocCtx, allocCancel := chromedp.NewExecAllocator(context.Background(),
		append(chromedp.DefaultExecAllocatorOptions[:],
			chromedp.Flag("headless", true),
			chromedp.Flag("disable-gpu", true),
			chromedp.Flag("no-sandbox", true),
			chromedp.Flag("disable-dev-shm-usage", true),
		)...,
	)

	ctx, cancel := chromedp.NewContext(allocCtx, chromedp.WithLogf(func(s string, args ...any) {
		log.Printf("[ZBrowser] %s", fmt.Sprintf(s, args...))
	}))

	m.allocCtx = allocCtx
	m.allocCancel = allocCancel
	m.ctx = ctx
	m.cancel = cancel
	m.running = true

	log.Println("[ZBrowser] started")
	return nil
}

func (m *Manager) Stop() {
	m.mu.Lock()
	defer m.mu.Unlock()
	if !m.running {
		return
	}
	if m.cancel != nil {
		m.cancel()
	}
	if m.allocCancel != nil {
		m.allocCancel()
	}
	m.running = false
	log.Println("[ZBrowser] stopped")
}

func (m *Manager) Navigate(url string) (*ActionResult, error) {
	ctx, cancel := context.WithTimeout(m.ctx, 30*time.Second)
	defer cancel()

	var buf []byte
	err := chromedp.Run(ctx,
		chromedp.Navigate(url),
		chromedp.FullScreenshot(&buf, 90),
	)
	if err != nil {
		return nil, fmt.Errorf("navigate: %w", err)
	}

	return &ActionResult{
		Success:  true,
		Data:     base64.StdEncoding.EncodeToString(buf),
		MIMEType: "image/png",
	}, nil
}

func (m *Manager) Screenshot(selector string) (*ActionResult, error) {
	ctx, cancel := context.WithTimeout(m.ctx, 15*time.Second)
	defer cancel()

	var buf []byte
	tasks := []chromedp.Action{chromedp.FullScreenshot(&buf, 90)}
	if selector != "" {
		tasks = []chromedp.Action{
			chromedp.WaitVisible(selector),
			chromedp.Screenshot(selector, &buf, chromedp.NodeVisible),
		}
	}

	err := chromedp.Run(ctx, tasks...)
	if err != nil {
		return nil, fmt.Errorf("screenshot: %w", err)
	}

	return &ActionResult{
		Success:  true,
		Data:     base64.StdEncoding.EncodeToString(buf),
		MIMEType: "image/png",
	}, nil
}

func (m *Manager) GetHTML(selector string) (*ActionResult, error) {
	ctx, cancel := context.WithTimeout(m.ctx, 10*time.Second)
	defer cancel()

	var html string
	tasks := []chromedp.Action{chromedp.OuterHTML(selector, &html, chromedp.ByQuery)}
	if selector == "" {
		tasks = []chromedp.Action{chromedp.OuterHTML("html", &html, chromedp.ByQuery)}
	}

	err := chromedp.Run(ctx, tasks...)
	if err != nil {
		return nil, fmt.Errorf("get html: %w", err)
	}

	return &ActionResult{Success: true, Data: html}, nil
}

func (m *Manager) Click(selector string) (*ActionResult, error) {
	ctx, cancel := context.WithTimeout(m.ctx, 10*time.Second)
	defer cancel()

	err := chromedp.Run(ctx,
		chromedp.WaitVisible(selector),
		chromedp.Click(selector, chromedp.NodeVisible),
	)
	if err != nil {
		return nil, fmt.Errorf("click: %w", err)
	}

	return &ActionResult{Success: true}, nil
}

func (m *Manager) Type(selector, text string) (*ActionResult, error) {
	ctx, cancel := context.WithTimeout(m.ctx, 10*time.Second)
	defer cancel()

	err := chromedp.Run(ctx,
		chromedp.WaitVisible(selector),
		chromedp.SendKeys(selector, text, chromedp.ByQuery),
	)
	if err != nil {
		return nil, fmt.Errorf("type: %w", err)
	}

	return &ActionResult{Success: true}, nil
}

func (m *Manager) Evaluate(js string) (*ActionResult, error) {
	ctx, cancel := context.WithTimeout(m.ctx, 10*time.Second)
	defer cancel()

	var result string
	err := chromedp.Run(ctx, chromedp.Evaluate(js, &result))
	if err != nil {
		return nil, fmt.Errorf("evaluate: %w", err)
	}

	return &ActionResult{Success: true, Data: result}, nil
}

func (m *Manager) GetCurrentURL() (*ActionResult, error) {
	ctx, cancel := context.WithTimeout(m.ctx, 5*time.Second)
	defer cancel()

	var url string
	err := chromedp.Run(ctx, chromedp.Location(&url))
	if err != nil {
		return nil, fmt.Errorf("get url: %w", err)
	}

	return &ActionResult{Success: true, Data: url}, nil
}

func (m *Manager) GetTitle() (*ActionResult, error) {
	ctx, cancel := context.WithTimeout(m.ctx, 5*time.Second)
	defer cancel()

	var title string
	err := chromedp.Run(ctx, chromedp.Title(&title))
	if err != nil {
		return nil, fmt.Errorf("get title: %w", err)
	}

	return &ActionResult{Success: true, Data: title}, nil
}
