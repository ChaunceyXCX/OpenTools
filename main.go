package main

import (
	"embed"
	"log"

	"github.com/wailsapp/wails/v3/pkg/application"
	"github.com/wailsapp/wails/v3/pkg/events"
)

//go:embed all:frontend/dist
var assets embed.FS

func main() {
	app := application.New(application.Options{
		Name:        "ZTools",
		Description: "A high-performance launcher and plugin platform",
		Services: []application.Service{
			application.NewService(&GreetService{}),
		},
		Assets: application.AssetOptions{
			Handler: application.AssetFileServerFS(assets),
		},
		Mac: application.MacOptions{
			ApplicationShouldTerminateAfterLastWindowClosed: false,
		},
	})

	mainWin := app.Window.NewWithOptions(application.WebviewWindowOptions{
		Title: "ZTools",
		Mac: application.MacWindow{
			InvisibleTitleBarHeight: 50,
			Backdrop:                application.MacBackdropTranslucent,
			TitleBar:                application.MacTitleBarHiddenInset,
		},
		URL: "/",
	})

	mainWin.OnWindowEvent(events.Common.WindowClosing, func(event *application.WindowEvent) {
		mainWin.Hide()
	})

	systray := app.SystemTray.New()
	systray.AttachWindow(mainWin)
	systray.WindowDebounce(200)

	app.OnShutdown(func() {
		log.Println("[ZTools] Application shutting down")
	})

	if err := app.Run(); err != nil {
		log.Fatal(err)
	}
}
