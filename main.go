package main

import (
	"embed"
	"log"
	"os"
	"os/signal"
	"path/filepath"
	"runtime"

	"github.com/ChaunceyXCX/OpenTools/internal/api"
	"github.com/ChaunceyXCX/OpenTools/internal/core/httpserver"
	"github.com/ChaunceyXCX/OpenTools/internal/core/mcp"
	"github.com/ChaunceyXCX/OpenTools/internal/core/updater"
	"github.com/ChaunceyXCX/OpenTools/internal/native"
	"github.com/wailsapp/wails/v3/pkg/application"
	"github.com/wailsapp/wails/v3/pkg/events"
	"golang.design/x/hotkey"
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
	historySvc := api.NewHistoryService()
	zbrowserSvc := api.NewZBrowserService()

	app := application.New(application.Options{
		Name:        "OpenTools",
		Description: "Open source launcher and plugin platform",
		Services: []application.Service{
			application.NewService(&GreetService{}),
			application.NewService(api.NewDatabaseService()),
			application.NewService(api.NewCommandsService()),
			application.NewService(clipSvc),
			application.NewService(pluginSvc),
			application.NewService(historySvc),
			application.NewService(zbrowserSvc),
		},
		Assets: application.AssetOptions{
			Handler: application.AssetFileServerFS(assets),
		},
		Mac: application.MacOptions{
			ApplicationShouldTerminateAfterLastWindowClosed: false,
		},
		SingleInstance: &application.SingleInstanceOptions{
			UniqueID: "com.opentools.app",
			OnSecondInstanceLaunch: func(data application.SecondInstanceData) {
				log.Println("[OpenTools] second instance")
			},
			ExitCode: 0,
		},
	})

	winW, winH := 800, 600
	cx, cy := 0, 0
	if primary := app.Screen.GetPrimary(); primary != nil {
		cx = primary.Bounds.X + (primary.Bounds.Width-winW)/2
		cy = primary.Bounds.Y + (primary.Bounds.Height-winH)/3
	}

	mainWin := app.Window.NewWithOptions(application.WebviewWindowOptions{
		Title:             "OpenTools",
		Width:             winW,
		Height:            winH,
		X:                 cx,
		Y:                 cy,
		Frameless:         true,
		AlwaysOnTop:       true,
		BackgroundColour:  application.NewRGB(27, 38, 54),
		URL:               "/",
	})

	mainWin.Hide()

	if runtime.GOOS == "linux" {
		native.SetFloatingWindow(mainWin.NativeWindow())
	}

	mainWin.OnWindowEvent(events.Common.WindowClosing, func(event *application.WindowEvent) {
		mainWin.Hide()
	})
	mainWin.OnWindowEvent(events.Common.WindowLostFocus, func(event *application.WindowEvent) {
		mainWin.Hide()
	})

	systray := app.SystemTray.New()
	systray.AttachWindow(mainWin)
	systray.WindowDebounce(200)
	systray.OnClick(func() {
		toggleWindow(mainWin)
	})

	httpSrv := httpserver.New(17891)
	if err := httpSrv.Start(); err != nil {
		log.Printf("[HTTP] failed: %v", err)
	}

	go mcp.StartMCPServer()
	clipSvc.StartMonitor()

	log.Printf("[OpenTools] version %s", updater.CurrentVersion())

	go registerGlobalHotkey(mainWin)

	app.OnShutdown(func() {
		clipSvc.StopMonitor()
		httpSrv.Stop()
		api.CloseDatabase()
	})

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt)
	go func() {
		<-sigCh
		app.Quit()
	}()

	if err := app.Run(); err != nil {
		log.Fatal(err)
	}
}

func registerGlobalHotkey(win application.Window) {
	hk := hotkey.New([]hotkey.Modifier{modAlt()}, hotkey.KeyZ)
	if hk == nil {
		log.Println("[Hotkey] failed to register Alt+Z")
		return
	}
	if err := hk.Register(); err != nil {
		log.Printf("[Hotkey] register Alt+Z error: %v", err)
		return
	}
	log.Println("[Hotkey] Alt+Z registered (global)")
	for range hk.Keydown() {
		toggleWindow(win)
	}
}

func toggleWindow(win application.Window) {
	if win.IsVisible() {
		win.Hide()
	} else {
		win.Center()
		win.Show()
		win.Focus()
	}
}

func modAlt() hotkey.Modifier {
	return hotkey.Mod1 // Mod1 = Alt on X11/Linux
}

func getDataDir() string {
	home, err := os.UserHomeDir()
	if err != nil {
		return filepath.Join(os.TempDir(), "opentools-data")
	}
	dir := filepath.Join(home, ".opentools", "data")
	if err := os.MkdirAll(dir, 0755); err != nil {
		return filepath.Join(os.TempDir(), "opentools-data")
	}
	return dir
}
