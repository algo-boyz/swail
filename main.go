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

//go:embed all:frontend/dist
var assets embed.FS

//go:embed tools/updater-key.pub
var updaterPublicKey []byte

func init() {
	application.RegisterEvent[string]("time")
}

func main() {
	app := application.New(application.Options{
		Name:        "AgentJetson",
		Description: "AgentJetson dashboard — voice query, evidence & edge awareness",
		Services: []application.Service{
			application.NewService(&AJService{}),
		},
		Assets: application.AssetOptions{
			Handler: application.AssetFileServerFS(assets),
		},
		Mac: application.MacOptions{
			ApplicationShouldTerminateAfterLastWindowClosed: true,
		},
	})

	app.Window.NewWithOptions(application.WebviewWindowOptions{
		Title:  "AgentJetson",
		Width:  1100,
		Height: 720,
		Mac: application.MacWindow{
			InvisibleTitleBarHeight: 50,
			Backdrop:                application.MacBackdropTranslucent,
			TitleBar:                application.MacTitleBarHiddenInset,
		},
		UseApplicationMenu: true,
		BackgroundColour:   application.NewRGB(6, 7, 15),
		URL:                "/",
	})

	gh, err := github.New(github.Config{
		Repository:    "algo-boyz/swail",
		ChecksumAsset: "SHA256SUMS",
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

	go func() {
		for {
			app.Event.Emit("time", time.Now().Format(time.RFC1123))
			time.Sleep(time.Second)
		}
	}()

	if err = app.Run(); err != nil {
		app.Logger.Error("run", "error", err)
	}
}
