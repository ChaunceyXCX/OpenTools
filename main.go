package main

import (
	"embed"
	"log"
	"os"
	"path/filepath"

	"github.com/ChaunceyXCX/OpenTools/internal/api"
	"github.com/ChaunceyXCX/OpenTools/internal/core"
	"github.com/wailsapp/wails/v3/pkg/application"
	"github.com/wailsapp/wails/v3/pkg/events"
)

//go:embed all:frontend/dist
var assets embed.FS

func main() {
	dataDir := getDataDir()

	lock, err := core.NewLockFile(dataDir)
	if err != nil {
		log.Fatalf("[ZTools] %v", err)
	}
	defer lock.Release()

	dbPath := filepath.Join(dataDir, "ztools.db")
	if err := api.InitDatabase(dbPath); err != nil {
		log.Fatalf("failed to init database: %v", err)
	}

	app := application.New(application.Options{
		Name:        "ZTools",
		Description: "A high-performance launcher and plugin platform",
		Services: []application.Service{
			application.NewService(&GreetService{}),
			application.NewService(api.NewDatabaseService()),
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
		log.Println("[ZTools] shutting down")
		api.CloseDatabase()
	})

	if err := app.Run(); err != nil {
		log.Fatal(err)
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
