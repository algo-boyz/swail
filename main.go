package main

import (
	"context"
	"embed"

	"time"

	"github.com/wailsapp/wails/v3/pkg/application"
	"github.com/wailsapp/wails/v3/pkg/updater"
	"github.com/wailsapp/wails/v3/pkg/updater/providers/github"
)

const currentVersion = "dev"

// Wails uses Go's `embed` package to embed the frontend files into the binary.
// Any files in the frontend/dist folder will be embedded into the binary and
// made available to the frontend.
// See https://pkg.go.dev/embed for more information.

//go:embed all:frontend/dist
var assets embed.FS

//go:embed tools/updater-key.pub
var updaterPublicKey []byte

func init() {
	// Register a custom event whose associated data type is string.
	// This is not required, but the binding generator will pick up registered events
	// and provide a strongly typed JS/TS API for them.
	application.RegisterEvent[string]("time")
}

// main function serves as the application's entry point. It initializes the application, creates a window,
// and starts a goroutine that emits a time-based event every second. It subsequently runs the application and
// logs any error that might occur.
func main() {

	// Create a new Wails application by providing the necessary options.
	// Variables 'Name' and 'Description' are for application metadata.
	// 'Assets' configures the asset server with the 'FS' variable pointing to the frontend files.
	// 'Bind' is a list of Go struct instances. The frontend has access to the methods of these instances.
	// 'Mac' options tailor the application when running an macOS.
	app := application.New(application.Options{
		Name:        "swail-app",
		Description: "Dashboard to SniGo",
		Services: []application.Service{
			application.NewService(&MpService{}),
		},
		Assets: application.AssetOptions{
			Handler: application.AssetFileServerFS(assets),
		},
		Mac: application.MacOptions{
			ApplicationShouldTerminateAfterLastWindowClosed: true,
		},
	})

	// Create a new window with the necessary options.
	// 'Title' is the title of the window.
	// 'Mac' options tailor the window when running on macOS.
	// 'BackgroundColour' is the background colour of the window.
	// 'URL' is the URL that will be loaded into the webview.
	app.Window.NewWithOptions(application.WebviewWindowOptions{
		Title: "Swail",
		// Window sized to the golden ratio (1000 / 618 ≈ 1.618).
		Width:  1000,
		Height: 618,
		Mac: application.MacWindow{
			InvisibleTitleBarHeight: 50,
			Backdrop:                application.MacBackdropTranslucent,
			TitleBar:                application.MacTitleBarHiddenInset,
		},
		UseApplicationMenu: true, // required for menu to appear on Win/Linux
		BackgroundColour:   application.NewRGB(6, 7, 15),
		URL:                "/",
	})

	gh, err := github.New(github.Config{
		Repository:    "algo-boyz/swail",
		ChecksumAsset: "SHA256SUMS", // sibling file with sha256 digests
	})
	if err != nil {
		app.Logger.Error("github init", "error", err)
	}

	if err := app.Updater.Init(updater.Config{
		CheckInterval:  6 * time.Hour,
		CurrentVersion: currentVersion,
		PublicKey:      updaterPublicKey,
		Providers:      []updater.Provider{gh},
		Window: &updater.BuiltinWindow{
			CSS: `:root { --accent: #ff6f00; --radius: 16px; }`,
		},
	}); err != nil {
		app.Logger.Error("update init", "error", err)
	}

	menu := app.Menu.New()
	app.Menu.SetApplicationMenu(menu)
	appMenu := menu.AddSubmenu("App")
	appMenu.Add("Check for Updates…").OnClick(func(*application.Context) {
		go func() {
			if err := app.Updater.CheckAndInstall(context.Background()); err != nil {
				app.Logger.Error("update", "error", err)
			}
		}()
	})

	// Create a goroutine that emits an event containing the current time every second.
	// The frontend can listen to this event and update the UI accordingly.
	go func() {
		for {
			app.Event.Emit("time", time.Now().Format(time.RFC1123))
			time.Sleep(time.Second)
		}
	}()

	// Run the application. This blocks until the application has been exited.
	// If an error occurred while running the application, log it and exit.
	if err = app.Run(); err != nil {
		app.Logger.Error("run", "error", err)
	}
}
