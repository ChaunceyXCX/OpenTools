package main

import (
	"embed"
	"log"
	"os"
	"path/filepath"

	"github.com/ChaunceyXCX/OpenTools/internal/api"
	"github.com/wailsapp/wails/v3/pkg/application"
	"github.com/wailsapp/wails/v3/pkg/events"
)

//go:embed all:frontend/dist
var assets embed.FS

func main() {
	dataDir := getDataDir()

	dbPath := filepath.Join(dataDir, "ztools.db")
	if err := api.InitDatabase(dbPath); err != nil {
		log.Fatalf("failed to init database: %v", err)
	}

	clipSvc := api.NewClipboardService()
	pluginSvc := api.NewPluginService()

	app := application.New(application.Options{
		Name:        "ZTools",
		Description: "A high-performance launcher and plugin platform",
		Services: []application.Service{
			application.NewService(&GreetService{}),
			application.NewService(api.NewDatabaseService()),
			application.NewService(api.NewCommandsService()),
			application.NewService(clipSvc),
			application.NewService(pluginSvc),
		},
		Assets: application.AssetOptions{
			Handler: application.AssetFileServerFS(assets),
		},
		Mac: application.MacOptions{
			ApplicationShouldTerminateAfterLastWindowClosed: false,
		},
		SingleInstance: &application.SingleInstanceOptions{
			UniqueID: "link.eiot.ztools",
			OnSecondInstanceLaunch: func(data application.SecondInstanceData) {
				log.Println("[ZTools] second instance detected, focusing existing window")
			},
			ExitCode: 0,
		},
		KeyBindings: map[string]func(window application.Window){
			"optionoralt+z": func(window application.Window) {
				toggleWindow(window)
			},
			"escape": func(window application.Window) {
				window.Hide()
			},
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

	if saved := api.LoadWindowState(); saved != nil {
		mainWin.SetBounds(application.Rect{
			X: saved.X, Y: saved.Y,
			Width: saved.Width, Height: saved.Height,
		})
	}

	mainWin.OnWindowEvent(events.Common.WindowDidMove, func(event *application.WindowEvent) {
		go saveWindowState(mainWin)
	})
	mainWin.OnWindowEvent(events.Common.WindowDidResize, func(event *application.WindowEvent) {
		go saveWindowState(mainWin)
	})

	mainWin.OnWindowEvent(events.Common.WindowClosing, func(event *application.WindowEvent) {
		saveWindowState(mainWin)
		mainWin.Hide()
	})

	mainWin.OnWindowEvent(events.Common.WindowLostFocus, func(event *application.WindowEvent) {
		mainWin.Hide()
	})

	systray := app.SystemTray.New()
	systray.AttachWindow(mainWin)
	systray.WindowDebounce(200)

	systray.OnClick(func() {
		mainWin.Show()
	})

	clipSvc.StartMonitor()

	app.OnShutdown(func() {
		log.Println("[ZTools] shutting down")
		clipSvc.StopMonitor()
		api.CloseDatabase()
	})

	if err := app.Run(); err != nil {
		log.Fatal(err)
	}
}

func saveWindowState(win application.Window) {
	b := win.Bounds()
	api.SaveWindowState(&api.WindowState{
		X: b.X, Y: b.Y,
		Width: b.Width, Height: b.Height,
	})
}

func toggleWindow(win application.Window) {
	if win.IsVisible() {
		win.Hide()
	} else {
		win.Show()
		win.Focus()
	}
}

func getDataDir() string {
	home, err := os.UserHomeDir()
	if err != nil {
		return filepath.Join(os.TempDir(), "ztools-data")
	}
	dir := filepath.Join(home, ".ztools", "data")
	if err := os.MkdirAll(dir, 0755); err != nil {
		return filepath.Join(os.TempDir(), "ztools-data")
	}
	return dir
}
