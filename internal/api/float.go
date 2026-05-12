package api

import (
	"log"

	"github.com/wailsapp/wails/v3/pkg/application"
)

type FloatBallService struct {
	win    *application.WebviewWindow
	showing bool
}

func NewFloatBallService() *FloatBallService {
	return &FloatBallService{}
}

func (s *FloatBallService) Show(app *application.App) {
	if s.showing {
		return
	}

	s.win = app.Window.NewWithOptions(application.WebviewWindowOptions{
		Title:  "OpenTools Float",
		Width:  48,
		Height: 48,
		Mac: application.MacWindow{
			TitleBar: application.MacTitleBarHiddenInset,
		},
		BackgroundColour: application.NewRGB(27, 38, 54),
		URL:              "/",
	})

	s.win.SetAlwaysOnTop(true)
	s.win.SetFrameless(true)
	s.win.Show()
	s.showing = true
	log.Println("[FloatBall] shown")
}

func (s *FloatBallService) Hide() {
	if s.win != nil {
		s.win.Hide()
		s.showing = false
	}
}

func (s *FloatBallService) Toggle(app *application.App) {
	if s.showing {
		s.Hide()
	} else {
		s.Show(app)
	}
}
